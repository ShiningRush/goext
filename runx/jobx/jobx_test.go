package jobx

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJobDemon(t *testing.T) {
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
	Stop()
	stopped = true
	assert.True(t, onceJobFlag)
	assert.True(t, startJobFlag)
	assert.True(t, intervalJobFlag)

	wg.Add(1)
	RemoveJob("interval-job")
	onceJobFlag, startJobFlag, intervalJobFlag = false, false, false
	Start()
	wg.Wait()
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
	demon.Stop()
	stopped = true
	assert.True(t, cronJobFlag)
}
