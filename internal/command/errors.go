package command

import "fmt"

var (
	ErrNotEnoughFuel = fmt.Errorf("not enough fuel")

	ErrUnsupportedDimension = fmt.Errorf("unsupported dimension")
)
