package main

import "time"

type timeSource interface {
	UTCNow() time.Time
}

func createTimeSource() timeSource {
	return &defaultTimeSource{}
}

type defaultTimeSource struct {
}

func (d defaultTimeSource) UTCNow() time.Time {
	return time.Now().UTC()
}
