package clock

import "time"

//go:generate mockgen -destination=../mock/clock.go -package mock -typed -source clock.go

type Clock interface {
	Parse(layout, value string) (time.Time, error)
	Since(t time.Time) time.Duration
}

type realClock struct{}

func NewClock() Clock {
	return &realClock{}
}

func (realClock) Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

func (realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}
