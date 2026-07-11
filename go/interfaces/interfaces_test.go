package interfaces

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterfaces(t *testing.T) {
	// initialize interfaces
	inlist := NewInList(1, 3, 5, 7, 11)
	instream := NewInStream(strings.NewReader("2 4 6 8 0"))
	var outlist OutList

	interfaces := []any{&inlist, instream, &outlist}

	// run for 5 cycles
	for cycle := 0; cycle < 5; cycle++ {
		for i, src := range interfaces {
			// iterate over all inputs
			input, isInput := src.(Input)
			if !isInput {
				// not an input, nothing to do
				continue
			}
			// get next value from this input
			val, err := input.NextValue(cycle)
			if !assert.NoErrorf(t, err, "interface %d errored when producing value at cycle %d", i, cycle) {
				return
			}
			// send value to all outputs
			for j, dest := range interfaces {
				// iterate over all outputs
				output, isOutput := dest.(Output)
				if !isOutput {
					// not an output, nothing to do
					continue
				}
				// send value to this output
				err = output.SendValue(val)
				if !assert.NoErrorf(t, err, "interface %d errored when sending value at cycle %d", j, cycle) {
					return
				}
			}
		}
	}

	// outlist should contain values from both inlist and instream
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 11, 0}, outlist.Result())
}
