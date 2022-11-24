package dns_test

import (
	"testing"

	"inetmock.icb4dc0.de/inetmock/protocols/dns"
)

func TestConditionalResolver_Matches(t *testing.T) {
	t.Parallel()
	type fields struct {
		Predicates []dns.QuestionPredicate
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Empty predicates - expect not to match",
			want: true,
		},
		{
			name: "Single predicate - expect to match",
			fields: fields{
				Predicates: []dns.QuestionPredicate{dns.QuestionPredicateFunc(func(dns.Question) bool {
					return true
				})},
			},
			want: true,
		},
		{
			name: "Single predicate - expect to not match",
			fields: fields{
				Predicates: []dns.QuestionPredicate{dns.QuestionPredicateFunc(func(dns.Question) bool {
					return false
				})},
			},
			want: false,
		},
		{
			name: "Multiple predicates - only first one matches",
			fields: fields{
				Predicates: []dns.QuestionPredicate{
					dns.QuestionPredicateFunc(func(dns.Question) bool {
						return true
					}),
					dns.QuestionPredicateFunc(func(dns.Question) bool {
						return false
					}),
				},
			},
			want: false,
		},
		{
			name: "Multiple predicates - only second one matches",
			fields: fields{
				Predicates: []dns.QuestionPredicate{
					dns.QuestionPredicateFunc(func(dns.Question) bool {
						return false
					}),
					dns.QuestionPredicateFunc(func(dns.Question) bool {
						return true
					}),
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := dns.ConditionalResolver{
				Predicates: tt.fields.Predicates,
			}
			if got := c.Matches(dns.Question{}); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
