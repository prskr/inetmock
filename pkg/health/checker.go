package health

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type checker map[string]Check

func (c *checker) AddCheck(check Check) error {
	self := *c
	name := check.Name()
	if _, ok := self[name]; ok {
		return ErrAmbiguousCheckName
	}
	self[name] = check
	return nil
}

func (c checker) Status(ctx context.Context) Result {
	rw := NewResultWriter()
	grp, grpCtx := errgroup.WithContext(ctx)

	for k, v := range c {
		// pin variables
		checkName := k
		check := v
		grp.Go(func() error {
			rw.WriteResult(checkName, check.Status(grpCtx))
			return nil
		})
	}

	_ = grp.Wait()
	return rw.GetResult()
}
