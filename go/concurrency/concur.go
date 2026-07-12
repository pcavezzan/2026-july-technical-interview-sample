package concurrency

import "sync"

type JobProgress interface {
	JobId() string
	Progression() int
}

type Progression struct {
	identifier string
	value      int
}

func (p Progression) JobId() string {
	return p.identifier
}

func (p Progression) Progression() int {
	return p.value
}

// Counter to keep track of progression in a go routine
type Counter struct {
	lck             sync.RWMutex       // RWMutex to protect value from concurrent access
	identifier      string             // Identifier of the Job being processed
	value           int                // Current progression
	progressChannel chan<- JobProgress // Channel where each progression step is published
}

// Increment progression and publish it to consumer
func (c *Counter) Inc() {
	c.lck.Lock()
	defer c.lck.Unlock()
	c.value++
	c.progressChannel <- Progression{
		identifier: c.identifier,
		value:      c.value,
	}
}

// Current progression of the job
func (c *Counter) Progression() int {
	c.lck.RLock()
	defer c.lck.RUnlock()
	return c.value
}

func (c *Counter) JobId() string {
	return c.identifier
}

func NewProgressCounter(id string, progress chan<- JobProgress) *Counter {
	return &Counter{
		identifier:      id,
		progressChannel: progress,
	}
}
