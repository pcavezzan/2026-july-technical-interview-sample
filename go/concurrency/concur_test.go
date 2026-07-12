package concurrency

import (
	"fmt"
	"sync"
	"testing"
)

const JOBID_FORMAT = "Worker %d"

// Helper to track job status
// Used only by one go routine (reader)
// for each progression step 0 ... 100
// stores all workerIds that reached this designated step
type JobTracker map[int]map[string]struct{}

func (jt JobTracker) RegisterProgression(jobId string, progress int) {
	track, exists := jt[progress]
	if !exists {
		jt[progress] = make(map[string]struct{})
		track = jt[progress]
	}
	track[jobId] = struct{}{}
	fmt.Println("RegisterProgression", jobId, progress)
}

// Check progress
// - <worker> go routines simulates Job workers that generates progress information from 1 to countTo
// - 1 go routine collects each progression step of each job worker, and checks:
//   - Each step of the progression is generated
//   - All workers get to the end
func TestCheckProgress(t *testing.T) {

	// Constant to ease tests
	worker := 5
	countTo := 100

	progress := make(chan JobProgress)
	readerWg := sync.WaitGroup{}
	readerWg.Add(1)

	readerTestSuccess := false
	// Creates a go routine inspecting progression to check all workers fetch the end
	go func() {
		result := JobTracker{}
		for p := range progress {
			result.RegisterProgression(p.JobId(), p.Progression())
		}

		if len(result) != countTo {
			t.Errorf("Some workers results are missing expected %d - got %d", countTo, len(result))
			return
		}

		for value := 1; value <= countTo; value++ {
			finishedJobs := result[value]
			if finishedCount := len(finishedJobs); finishedCount != worker {
				t.Errorf("Invalid count for progression value %d step (expected %d - got %d)", value, worker, finishedCount)
				return
			}
		}
		readerTestSuccess = true
		readerWg.Done()
	}()

	// Creates <worker> go routines to count steps up to <countTo> each
	i := 0
	workerWg := sync.WaitGroup{}
	for i < worker {
		workerWg.Add(1)
		go func(i int) {
			c := NewProgressCounter(fmt.Sprintf(JOBID_FORMAT, i), progress)
			for i := 0; i < countTo; i++ {
				c.Inc()
			}
			workerWg.Done()
		}(i)
		i++
	}

	// Waits for the end of all workers go routine
	workerWg.Wait()
	close(progress)
	readerWg.Wait()
	if !readerTestSuccess {
		t.Errorf("Reader tests not performed")
	}
}
