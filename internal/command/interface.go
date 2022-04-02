package command

import "modules/internal/vector"

type Command interface {
	Execute() error
}

type ErrorHandler func(command Command, err error)

type Queue interface {
	Get() (Command, bool)
	Put(Command)
}

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

type FuelBurnable interface {
	GetFuel() (int, error)
	GetConsumption() (int, error)
	SetFuel(int) error
}

type MovableWithFuel interface {
	Movable
	FuelBurnable
}
