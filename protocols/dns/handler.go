package dns

import (
	"errors"

	mdns "github.com/miekg/dns"
)

var ErrNoAnswerForQuestion = errors.New("cannot answer given question")

type (
	Question       mdns.Question
	ResourceRecord mdns.RR
	HandlerFunc    func(q Question) (ResourceRecord, error)
)

func (f HandlerFunc) AnswerDNSQuestion(q Question) (ResourceRecord, error) {
	return f(q)
}

type Handler interface {
	AnswerDNSQuestion(q Question) (ResourceRecord, error)
}
