package adapters

import "modules/internal/vector"

type Movable interface {
	GetPosition() vector.Vector
	GetVelocity() vector.Vector
	SetPosition(vector.Vector)
}

type Rotatable interface {
	GetDirection() int
	GetAngularVelocity() int
	SetDirection(direction int)
	GetDirectionsNumber() int
}
