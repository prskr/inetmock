package health

import "context"

func NewCheckFunc(name string, delegate func(ctx context.Context) error) Check {
	return &checkDelegate{
		name:           name,
		statusDelegate: delegate,
	}
}

type checkDelegate struct {
	name           string
	statusDelegate func(ctx context.Context) error
}

func (c checkDelegate) Name() string {
	return c.name
}

func (c checkDelegate) Status(ctx context.Context) error {
	if c.statusDelegate == nil {
		return nil
	}
	return c.statusDelegate(ctx)
}
