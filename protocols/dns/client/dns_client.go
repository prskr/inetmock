package client

import (
	"context"
	"net"

	mdns "github.com/miekg/dns"
)

var (
	stringMapper = func(rr mdns.RR) string {
		switch e := rr.(type) {
		case *mdns.MX:
			return e.Mx
		case *mdns.PTR:
			return e.Ptr
		default:
			return ""
		}
	}

	ipMapper = func(rr mdns.RR) net.IP {
		switch e := rr.(type) {
		case *mdns.A:
			return e.A
		case *mdns.AAAA:
			return e.AAAA
		default:
			return nil
		}
	}
)

type RoundTripper interface {
	RoundTrip(ctx context.Context, question *mdns.Msg) (resp *mdns.Msg, err error)
}

type Resolver struct {
	Transport RoundTripper
}

func (r Resolver) Do(ctx context.Context, msg *mdns.Msg) (resp *mdns.Msg, err error) {
	return r.Transport.RoundTrip(ctx, msg)
}

func (r Resolver) Question(ctx context.Context, question mdns.Question) (resp *mdns.Msg, err error) {
	req := new(mdns.Msg)
	req.Id = mdns.Id()
	req.RecursionDesired = true
	req.Question = []mdns.Question{question}
	return r.Do(ctx, req)
}

func (r Resolver) LookupA(ctx context.Context, host string) (res []net.IP, err error) {
	msg := new(mdns.Msg).SetQuestion(mdns.Fqdn(host), mdns.TypeA)

	if resp, err := r.Do(ctx, msg); err != nil {
		return nil, err
	} else {
		return unwrapIPSlice(resp.Answer, ipMapper), nil
	}
}

func (r Resolver) LookupAAAA(ctx context.Context, host string) (res []net.IP, err error) {
	msg := new(mdns.Msg).SetQuestion(mdns.Fqdn(host), mdns.TypeAAAA)

	if resp, err := r.Do(ctx, msg); err != nil {
		return nil, err
	} else {
		return unwrapIPSlice(resp.Answer, ipMapper), nil
	}
}

func (r Resolver) LookupPTR(ctx context.Context, inAddrArpa string) (res []string, err error) {
	if i := net.ParseIP(inAddrArpa); i != nil {
		if inAddrArpa, err = mdns.ReverseAddr(inAddrArpa); err != nil {
			return nil, err
		}
	}
	msg := new(mdns.Msg).SetQuestion(mdns.Fqdn(inAddrArpa), mdns.TypePTR)

	if resp, err := r.Do(ctx, msg); err != nil {
		return nil, err
	} else {
		return unwrapStringSlice(resp.Answer, stringMapper), nil
	}
}

func (r Resolver) LookupMX(ctx context.Context, domain string) (res []string, err error) {
	msg := new(mdns.Msg).SetQuestion(mdns.Fqdn(domain), mdns.TypeMX)
	if resp, err := r.Do(ctx, msg); err != nil {
		return nil, err
	} else {
		return unwrapStringSlice(resp.Answer, stringMapper), nil
	}
}

func unwrapStringSlice(records []mdns.RR, mapper func(rr mdns.RR) string) []string {
	out := make([]string, 0, len(records))
	for idx := range records {
		if val := mapper(records[idx]); val != "" {
			out = append(out, val)
		}
	}
	return out
}

func unwrapIPSlice(records []mdns.RR, mapper func(rr mdns.RR) net.IP) []net.IP {
	out := make([]net.IP, 0, len(records))
	for idx := range records {
		if val := mapper(records[idx]); val != nil {
			out = append(out, val)
		}
	}
	return out
}
