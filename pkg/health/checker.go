package health

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

type ResultWriter interface {
	WriteResult(checkName string, result error)
	GetResult() Result
}

func NewResultWriter() ResultWriter {
	return &resultWriter{
		lock:   new(sync.Mutex),
		result: Result{},
	}
}

type resultWriter struct {
	lock   sync.Locker
	result Result
}

func (r *resultWriter) WriteResult(checkName string, result error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.result[checkName] = result
}

func (r resultWriter) GetResult() Result {
	return r.result
}

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

func (c checker) Status(ctx context.Context) (res Result, err error) {
	rw := NewResultWriter()
	grp, grpCtx := errgroup.WithContext(ctx)

	for k, v := range c {
		grp.Go(func() error {
			checkErr := v.Status(grpCtx)
			rw.WriteResult(k, checkErr)
			return checkErr
		})
	}

	err = grp.Wait()
	res = rw.GetResult()
	return
}
