package test

import (
	"fmt"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/td"
)

type genericMatcher struct {
	tb       testing.TB
	expected any
}

func GenericMatcher(tb testing.TB, expected any) gomock.Matcher {
	tb.Helper()
	return &genericMatcher{
		tb:       tb,
		expected: expected,
	}
}

func (g genericMatcher) Matches(x any) bool {
	g.tb.Helper()
	return td.Cmp(g.tb, x, g.expected)
}

func (g genericMatcher) String() string {
	return g.tb.Name()
}

func IP(rawIP string) td.TestDeep {
	parsed := net.ParseIP(rawIP)
	return td.Code(func(other net.IP) error {
		if !parsed.Equal(other) {
			return fmt.Errorf("expected IP %s", rawIP)
		}
		return nil
	})
}
