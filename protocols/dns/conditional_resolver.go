package dns

type QuestionPredicate interface {
	Matches(q Question) bool
}

type QuestionPredicateFunc func(q Question) bool

func (f QuestionPredicateFunc) Matches(q Question) bool {
	return f(q)
}

type ConditionalResolver struct {
	IPResolver
	Predicates []QuestionPredicate
}

func (c ConditionalResolver) Matches(q Question) bool {
	for idx := range c.Predicates {
		if !c.Predicates[idx].Matches(q) {
			return false
		}
	}
	return true
}
