package queue

import "time"

type Entry struct {
	Key     string
	Value   any
	timeout time.Time
	index   int
}

func (e Entry) WithTTL(ttl time.Duration) *Entry {
	e.timeout = time.Now().UTC().Add(ttl)
	return &e
}

func (e Entry) TTL() time.Time {
	return e.timeout
}
