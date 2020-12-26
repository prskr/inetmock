package health

import "fmt"

type checker struct {
	componentChecks map[string]Check
}

type Checker interface {
	RegisterCheck(component string, check Check) error
	IsHealthy() Result
}

func (c *checker) RegisterCheck(component string, check Check) error {
	if _, exists := c.componentChecks[component]; exists {
		return fmt.Errorf("component: %s: %w", component, ErrCheckForComponentAlreadyRegistered)
	}

	c.componentChecks[component] = check
	return nil
}

func (c *checker) IsHealthy() (r Result) {
	r.Status = HEALTHY
	r.Components = make(map[string]CheckResult)
	for component, componentCheck := range c.componentChecks {
		r.Components[component] = componentCheck()
		r.Status = max(r.Components[component].Status, r.Status)
	}
	return
}

func max(s1, s2 Status) Status {
	var max Status
	if s1 > s2 {
		max = s1
	} else {
		max = s2
	}

	if max > UNHEALTHY {
		return UNHEALTHY
	}
	return max
}
