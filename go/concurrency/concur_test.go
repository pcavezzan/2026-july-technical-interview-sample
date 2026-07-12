package concurrency

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
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
}

// Check progress
// - <worker> go routines simulates Job workers that generates progress information from 1 to countTo
// - 1 go routine collects each progression step of each job worker, and checks:
//   - Each step of the progression is generated
//   - All workers get to the end
func TestCheckProgress(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		worker  int
		countTo int
	}{
		{name: "5 workers / 100 pas", worker: 5, countTo: 100},
		{name: "1 worker / 1000 pas", worker: 1, countTo: 1000},
		{name: "50 workers / 20 pas", worker: 50, countTo: 20},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			progress := make(chan JobProgress)
			wg := sync.WaitGroup{}

			var readerErr error
			tracker := sync.WaitGroup{}
			tracker.Go(func() {
				result := JobTracker{}
				for p := range progress {
					result.RegisterProgression(p.JobId(), p.Progression())
				}
				readerErr = checkResult(result, tc.worker, tc.countTo)
			})

			for i := 1; i <= tc.worker; i++ {
				wg.Go(func() {
					c := NewProgressCounter(fmt.Sprintf(JOBID_FORMAT, i), progress)
					for j := 0; j < tc.countTo; j++ {
						c.Inc()
					}
				})
			}

			wg.Wait()
			close(progress)
			tracker.Wait()

			require.NoError(t, readerErr)
		})
	}
}

func checkResult(result JobTracker, worker int, countTo int) error {
	for value := 1; value <= countTo; value++ {
		finishedJobs := result[value]
		if finishedCount := len(finishedJobs); finishedCount != worker {
			return fmt.Errorf("Invalid count for progression value %d step (expected %d - got %d)", value, worker, finishedCount)
		}
	}
	return nil
}
