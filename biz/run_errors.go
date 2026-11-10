package biz

import (
	"fmt"
	"sync"
)

// concurrentRunErrors never lets a successful worker clear an earlier error.
type concurrentRunErrors struct {
	mu  sync.Mutex
	err error
}

func (e *concurrentRunErrors) add(err error) {
	if err == nil {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.err = combineRunErrors(e.err, err)
}

func (e *concurrentRunErrors) get() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.err
}

func combineRunErrors(first, next error) error {
	if first == nil {
		return next
	}
	if next == nil {
		return first
	}
	return fmt.Errorf("%v; %v", first, next)
}
