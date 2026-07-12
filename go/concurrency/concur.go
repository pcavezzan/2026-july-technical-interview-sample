package concurrency

type JobProgress interface {
	JobId() string
	Progression() int
}

// Counter to keep track of progression in a go routine
type Counter struct {
	identifier      string             // Identifier of the Job being processed
	value           int                // Current progression
	progressChannel chan<- JobProgress // Channel where each progression step is published
}

// Increment progression and publish it to consumer
func (c *Counter) Inc() {
	c.value++
	c.progressChannel <- *c
}

// Current progression of the job
func (c Counter) Progression() int {
	return c.value
}

func (c Counter) JobId() string {
	return c.identifier
}

func NewProgressCounter(id string, progress chan<- JobProgress) *Counter {
	return &Counter{
		identifier:      id,
		progressChannel: progress,
	}
}
