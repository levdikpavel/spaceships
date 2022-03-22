package adapters

import "modules/internal/vector"

func NewMove(m Movable) *Move {
	return &Move{
		m: m,
	}
}

type Move struct {
	m Movable
}

func (m *Move) Execute() error {
	position, err := m.m.GetPosition()
	if err != nil {
		return err
	}

	velocity, err := m.m.GetVelocity()
	if err != nil {
		return err
	}

	err = m.m.SetPosition(vector.Add(position, velocity))
	if err != nil {
		return err
	}

	return nil
}

func NewRotate(r Rotatable) *Rotate {
	return &Rotate{
		r: r,
	}
}

type Rotate struct {
	r Rotatable
}

func (r *Rotate) Execute() error {
	direction, err := r.r.GetDirection()
	if err != nil {
		return err
	}

	angularVelocity, err := r.r.GetAngularVelocity()
	if err != nil {
		return err
	}

	n, err := r.r.GetDirectionsNumber()
	if err != nil {
		return err
	}

	directionNew := direction + angularVelocity
	err = r.r.SetDirection(directionNew % n)
	if err != nil {
		return err
	}

	return nil
}
