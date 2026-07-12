package interfaces

import "imagino.com/interfaces/v2/list"

// InList produces values from a list
type InList struct {
	l list.List[int]
}

// Implements the Input interface
func (in *InList) NextValue(cycle int) (int, error) {
	return in.l.Get(cycle), nil
}

func NewInList(values ...int) InList {
	var inlist InList
	inlist.l.Append(values...)
	return inlist
}
