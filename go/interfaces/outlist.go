package interfaces

import "imagino.com/interfaces/v2/list"

// OutList stores output values in a list
type OutList struct {
	l list.List[int]
}

// Implements the Output interface
func (out *OutList) SendValue(value int) error {
	out.l.Append(value)
	return nil
}

// Get the values that were stored in out
func (out *OutList) Result() []int {
	return out.l.Slice()
}
