package testing

import "sync"

type Counter struct {
	mux    sync.Mutex
	counts map[string]int
}

func NewCounter() *Counter {
	return &Counter{
		counts: make(map[string]int),
	}
}

func (c *Counter) Increment(name string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, ok := c.counts[name]; !ok {
		c.counts[name] = 0
	}

	c.counts[name]++
}

func (c *Counter) Count(name string) int {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, ok := c.counts[name]; !ok {
		return 0
	}

	return c.counts[name]
}
