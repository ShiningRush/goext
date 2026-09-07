package jobx

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJobDemon(t *testing.T) {
	defer func() {
		Close()
		ClearJobs()
	}()

	wg := sync.WaitGroup{}
	wg.Add(10)
	now := time.Now()

	onceJobFlag, startJobFlag, intervalJobFlag := false, false, false
	stopped := false
	RegisterJob("once-job", JobType{Once: &OnceJobDesc{}}, func(ctx context.Context) {
		assert.LessOrEqual(t, time.Since(now).Milliseconds(), int64(1000))
		wg.Done()
		onceJobFlag = true
	})

	RegisterJobDesc(&JobDescriptor{
		Name: "start-job",
		Type: JobType{Once: &OnceJobDesc{
			Delay:       800 * time.Millisecond,
			AlwaysStart: true,
		}},
		Func: func(ctx context.Context) {
			wg.Done()
			startJobFlag = true
		},
	})

	RegisterJob("interval-job", JobType{Interval: &IntervalJobDesc{Interval: 200 * time.Millisecond}}, func(ctx context.Context) {
		if stopped {
			assert.Fail(t, "interval-job should not be fired after stopped")
		}

		wg.Done()
		intervalJobFlag = true
	})
	Start()
	wg.Wait()
	Close()
	stopped = true
	assert.True(t, onceJobFlag)
	assert.True(t, startJobFlag)
	assert.True(t, intervalJobFlag)

	wg.Add(1)
	RemoveJob("interval-job")
	onceJobFlag, startJobFlag, intervalJobFlag = false, false, false
	Start()
	wg.Wait()
	Close()
	assert.False(t, onceJobFlag)
	assert.True(t, startJobFlag)
	assert.False(t, intervalJobFlag)
}

func TestJobDemonCron(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	demon := NewJobDemon()
	cronJobFlag := false
	stopped := false

	demon.RegisterJob("cron-job", JobType{Cron: &CronJobDesc{Spec: "* * * * * *"}}, func(ctx context.Context) {
		if stopped {
			assert.Fail(t, "cron-job should not be fired after stopped")
		}

		if !cronJobFlag {
			wg.Done()
			cronJobFlag = true
		}
	})

	demon.Start()
	wg.Wait()
	demon.Close()
	stopped = true
	assert.True(t, cronJobFlag)
}

func TestJobDemonStopCancelsRunningJob(t *testing.T) {
	started := make(chan struct{})
	cancelled := make(chan struct{})
	release := make(chan struct{})

	demon := NewJobDemon()
	demon.RegisterJob("once-job", JobType{Once: &OnceJobDesc{}}, func(ctx context.Context) {
		close(started)
		<-ctx.Done()
		close(cancelled)
		<-release
	})

	demon.Start()
	<-started

	stopReturned := make(chan struct{})
	go func() {
		demon.Stop()
		close(stopReturned)
	}()

	select {
	case <-stopReturned:
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "Stop should not wait for running jobs")
	}

	select {
	case <-cancelled:
	case <-time.After(time.Second):
		assert.Fail(t, "running job should receive cancellation")
	}

	close(release)
	demon.Close()
}

func TestJobDemonCloseWaitsRunningJob(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})

	demon := NewJobDemon()
	demon.RegisterJob("once-job", JobType{Once: &OnceJobDesc{}}, func(ctx context.Context) {
		close(started)
		<-release
	})

	demon.Start()
	<-started

	closeReturned := make(chan struct{})
	go func() {
		demon.Close()
		close(closeReturned)
	}()

	select {
	case <-closeReturned:
		assert.Fail(t, "Close should wait for running jobs")
	case <-time.After(100 * time.Millisecond):
	}

	close(release)
	select {
	case <-closeReturned:
	case <-time.After(time.Second):
		assert.Fail(t, "Close should return after running jobs finish")
	}
}

func TestJobDemonCloseStopsDelayedOnceJob(t *testing.T) {
	fired := make(chan struct{})

	demon := NewJobDemon()
	demon.RegisterJob("once-job", JobType{Once: &OnceJobDesc{Delay: time.Second}}, func(ctx context.Context) {
		close(fired)
	})

	demon.Start()
	demon.Close()

	select {
	case <-fired:
		assert.Fail(t, "delayed once job should not fire after Close")
	default:
	}
}

func TestJobDescriptorDo(t *testing.T) {
	fired := make(chan struct{})
	closeChan := make(chan struct{})
	job := &JobDescriptor{
		Type: JobType{Interval: &IntervalJobDesc{Interval: time.Millisecond}},
		Func: func(ctx context.Context) {
			close(fired)
			<-ctx.Done()
		},
	}

	job.Do(context.Background(), closeChan)
	select {
	case <-fired:
	case <-time.After(time.Second):
		assert.Fail(t, "job should fire")
	}
	close(closeChan)
}

func TestJobDemonZeroInterval(t *testing.T) {
	started := make(chan struct{})

	demon := NewJobDemon()
	demon.RegisterJob("interval-job", JobType{Interval: &IntervalJobDesc{}}, func(ctx context.Context) {
		close(started)
		<-ctx.Done()
	})

	demon.Start()
	select {
	case <-started:
	case <-time.After(time.Second):
		assert.Fail(t, "zero interval job should fire")
	}
	demon.Close()
}

func TestJobDemonIntervalStartsAfterPreviousRunAndDelay(t *testing.T) {
	const interval = 50 * time.Millisecond

	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})
	secondStarted := make(chan struct{})
	var runs int32

	demon := NewJobDemon()
	demon.RegisterJob("interval-job", JobType{Interval: &IntervalJobDesc{Interval: interval}}, func(ctx context.Context) {
		switch atomic.AddInt32(&runs, 1) {
		case 1:
			close(firstStarted)
			<-releaseFirst
		case 2:
			close(secondStarted)
		}
	})

	demon.Start()
	<-firstStarted
	time.Sleep(2 * interval)
	close(releaseFirst)

	select {
	case <-secondStarted:
		assert.Fail(t, "next interval job should not start immediately after the previous run")
	case <-time.After(interval / 2):
	}

	select {
	case <-secondStarted:
	case <-time.After(2 * interval):
		assert.Fail(t, "next interval job should start after the configured delay")
	}
	demon.Close()
}
