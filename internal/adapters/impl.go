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

func (m *Move) Execute() {
	position := m.m.GetPosition()
	velocity := m.m.GetVelocity()
	m.m.SetPosition(vector.Add(position, velocity))
}

func NewRotate(r Rotatable) *Rotate {
	return &Rotate{
		r: r,
	}
}

type Rotate struct {
	r Rotatable
}

func (r *Rotate) Execute() {
	direction := r.r.GetDirection()
	angularVelocity := r.r.GetAngularVelocity()
	n := r.r.GetDirectionsNumber()
	directionNew := direction + angularVelocity
	r.r.SetDirection(directionNew % n)
}
