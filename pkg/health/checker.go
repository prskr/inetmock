package health

import (
	"context"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type checker struct {
	registeredChecks map[string]Check
	checkClient      *http.Client
}

func (c *checker) AddCheck(check Check) error {
	name := check.Name()
	if _, ok := c.registeredChecks[name]; ok {
		return ErrAmbiguousCheckName
	}
	c.registeredChecks[name] = check
	return nil
}

func (c checker) Status(ctx context.Context) Result {
	rw := NewResultWriter()
	grp, grpCtx := errgroup.WithContext(ctx)

	for k, v := range c.registeredChecks {
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
