//go:generate mockgen -source=$GOFILE -destination=./../../internal/mock/cert/time_source.mock.go -package=cert_mock

package cert

import "time"

type TimeSource interface {
	UTCNow() time.Time
}

func NewTimeSource() TimeSource {
	return &defaultTimeSource{}
}

type defaultTimeSource struct{}

func (d defaultTimeSource) UTCNow() time.Time {
	return time.Now().UTC()
}
