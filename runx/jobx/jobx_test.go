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
		assert.LessOrEqual(t, time.Now().Sub(now).Milliseconds(), int64(1000))
		wg.Done()
		onceJobFlag = true
	})

	RegisterJob("start-job", JobType{Once: &OnceJobDesc{
		Delay:       800 * time.Millisecond,
		AlwaysStart: true,
	}}, func(ctx context.Context) {
		wg.Done()
		startJobFlag = true
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
