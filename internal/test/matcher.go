package test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maxatome/go-testdeep/td"
)

type genericMatcher struct {
	tb       testing.TB
	expected interface{}
}

func GenericMatcher(tb testing.TB, expected interface{}) gomock.Matcher {
	tb.Helper()
	return &genericMatcher{
		tb:       tb,
		expected: expected,
	}
}

func (g genericMatcher) Matches(x interface{}) bool {
	g.tb.Helper()
	return td.Cmp(g.tb, x, g.expected)
}

func (g genericMatcher) String() string {
	return g.tb.Name()
}
