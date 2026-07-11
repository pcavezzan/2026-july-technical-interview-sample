// Helper package for managing lists, does not need fixing
package list

// List is a wrapper around go lists that adds methods for ease of use
type List[T any] struct {
	inner []T
}

// Append pushes values to the end of the list
func (l *List[T]) Append(values ...T) {
	l.inner = append(l.inner, values...)
}

// Slice returns the contents of l as a go slice
func (l *List[T]) Slice() []T {
	return l.inner
}

// Len returns the number of elements in l
func (l *List[T]) Len() int {
	return len(l.inner)
}

// Get returns the element at position index, or the zero value of T if out of bounds
func (l *List[T]) Get(index int) T {
	if index < 0 || index >= l.Len() {
		var zero T
		return zero
	}
	return l.inner[index]
}
