package interfaces

import (
	"fmt"
	"io"
)

// InStream scans and produces numbers from an input stream of space-separated single digits
type InStream struct {
	stream io.ReadSeeker
}

func (in *InStream) NextValue(cycle int) (int, error) {
	var val int
	// offset into stream with cycle number
	_, err := in.stream.Seek(int64(cycle*2), io.SeekStart)
	if err == nil {
		_, err = fmt.Fscan(in.stream, &val)
	}
	return val, err
}

func NewInStream(r io.ReadSeeker) *InStream {
	return &InStream{
		stream: r,
	}
}
