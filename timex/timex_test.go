//go:generate mockgen -source timex.go  -destination timex_mock.go -package timex
package timex

import (
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	xNow, realNow := Now(), time.Now()
	assert.True(t, realNow.Sub(xNow).Milliseconds() < 1, "default implement should be same(gap is less than 1ms)")

	ctr := gomock.NewController(t)
	defer ctr.Finish()

	mClock := NewMockClock(ctr)
	mClock.EXPECT().Now().Return(xNow)
	MockClockImpl(mClock)
	mockNow := Now()
	assert.Equal(t, mockNow, xNow)
}

func TestAfter(t *testing.T) {
	xAfter, realAfter := After(time.Second), time.After(time.Second)
	wg := sync.WaitGroup{}
	wg.Add(2)

	var xAfterTime, realAfterTime time.Time
	go func() {
		xAfterTime = <-xAfter
		wg.Done()
	}()
	go func() {
		realAfterTime = <-realAfter
		wg.Done()
	}()
	wg.Wait()

	assert.True(t, realAfterTime.Sub(xAfterTime).Milliseconds() < 1, "default implement should be same(gap is less than 1ms)")

	ctr := gomock.NewController(t)
	defer ctr.Finish()

	mClock := NewMockClock(ctr)
	mClock.EXPECT().After(gomock.Any()).Return(nil)
	MockClockImpl(mClock)
	assert.Nil(t, After(time.Second))
}
