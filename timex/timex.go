package timex

import "time"

var clockImpl Clock = &realClock{}

// MockClockImpl just used for mocking
func MockClockImpl(c Clock) {
	clockImpl = c
}

type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

type realClock struct {
}

func (r *realClock) Now() time.Time {
	return time.Now()
}

func (r *realClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func Now() time.Time {
	return clockImpl.Now()
}

func After(d time.Duration) <-chan time.Time {
	return clockImpl.After(d)
}
