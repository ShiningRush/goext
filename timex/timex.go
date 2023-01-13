package timex

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

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

var durationRE = regexp.MustCompile("^([0-9]+)(d|h|m|s|ms)$")

// ParseDuration support covert "1h,3h" to time.Duration.
func ParseDuration(durationStr string) (time.Duration, error) {
	matches := durationRE.FindStringSubmatch(durationStr)
	if len(matches) != 3 {
		return 0, fmt.Errorf("not a valid duration string: %q", durationStr)
	}
	var (
		n, _ = strconv.Atoi(matches[1])
		dur  = time.Duration(n) * time.Millisecond
	)
	switch unit := matches[2]; unit {
	case "d":
		dur *= 1000 * 60 * 60 * 24
	case "h":
		dur *= 1000 * 60 * 60
	case "m":
		dur *= 1000 * 60
	case "s":
		dur *= 1000
	case "ms":
	default:
		return 0, fmt.Errorf("invalid time unit in duration string: %q", unit)
	}
	return dur, nil
}
