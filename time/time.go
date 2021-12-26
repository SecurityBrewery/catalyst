package time

import "time"

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

var DefaultClock Clock = &realClock{}

func Now() time.Time {
	return DefaultClock.Now()
}
