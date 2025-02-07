package healthcheck

import (
	"errors"
	"sync"
)

type StartStopChecker struct {
	err   error
	guard sync.RWMutex
}

func NewBaseChecker() *StartStopChecker {
	return &StartStopChecker{
		err: errors.New("not started"),
	}
}

func (c *StartStopChecker) Start() {
	c.guard.Lock()
	defer c.guard.Unlock()
	c.err = nil
}

func (c *StartStopChecker) Stop() {
	c.guard.Lock()
	defer c.guard.Unlock()
	c.err = errors.New("stopping")
}

func (c *StartStopChecker) Check() error {
	c.guard.RLock()
	defer c.guard.RUnlock()
	return c.err
}
