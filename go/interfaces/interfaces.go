package interfaces

type Input interface {
	NextValue(cycle int) (int, error)
}

type Output interface {
	SendValue(value int) error
}
