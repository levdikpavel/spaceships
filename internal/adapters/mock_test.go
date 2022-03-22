package adapters

import (
	"github.com/stretchr/testify/mock"

	"modules/internal/vector"
)

type MovableMock struct {
	mock.Mock
}

func (m *MovableMock) GetPosition() vector.Vector {
	args := m.Called()
	return args.Get(0).(vector.Vector)
}

func (m *MovableMock) GetVelocity() vector.Vector {
	args := m.Called()
	return args.Get(0).(vector.Vector)
}

func (m *MovableMock) SetPosition(v vector.Vector) {
	m.Called(v)
}
