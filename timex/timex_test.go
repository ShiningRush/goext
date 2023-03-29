//go:generate mockgen -source timex.go  -destination timex_mock.go -package timex
package timex

import (
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNow(t *testing.T) {
	MockClockImpl(&realClock{})

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
	MockClockImpl(&realClock{})

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

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name         string
		giveText     string
		wantDuration time.Duration
		wantErr      require.ErrorAssertionFunc
	}{
		{
			name:         "normal-hour",
			giveText:     "4h",
			wantDuration: time.Hour * 4,
			wantErr:      require.NoError,
		},
		{
			name:         "normal-minute",
			giveText:     "4m",
			wantDuration: time.Minute * 4,
			wantErr:      require.NoError,
		},
		{
			name:         "normal-second",
			giveText:     "4s",
			wantDuration: time.Second * 4,
			wantErr:      require.NoError,
		},
		{
			name:         "normal-ms",
			giveText:     "4ms",
			wantDuration: time.Millisecond * 4,
			wantErr:      require.NoError,
		},
		{
			name:     "not-support-week",
			giveText: "4w",
			wantErr:  require.Error,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d, err := ParseDuration(tc.giveText)
			tc.wantErr(t, err)
			assert.Equal(t, tc.wantDuration, d)
		})
	}
}
