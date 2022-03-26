package adapters

import "modules/internal/vector"

type Movable interface {
	GetPosition() (vector.Vector, error)
	GetVelocity() (vector.Vector, error)
	SetPosition(vector.Vector) error
}

type Rotatable interface {
	GetDirection() (int, error)
	GetAngularVelocity() (int, error)
	SetDirection(direction int) error
	GetDirectionsNumber() (int, error)
}
