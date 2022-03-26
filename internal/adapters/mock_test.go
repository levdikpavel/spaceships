package adapters

import (
	"github.com/stretchr/testify/mock"

	"modules/internal/vector"
)

type MovableMock struct {
	mock.Mock
}

func (m *MovableMock) GetPosition() (vector.Vector, error) {
	args := m.Called()
	return args.Get(0).(vector.Vector), args.Error(1)
}

func (m *MovableMock) GetVelocity() (vector.Vector, error) {
	args := m.Called()
	return args.Get(0).(vector.Vector), args.Error(1)
}

func (m *MovableMock) SetPosition(v vector.Vector) error {
	args := m.Called(v)
	return args.Error(0)
}

type RotatableMock struct {
	mock.Mock
}

func (m *RotatableMock) GetDirection() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *RotatableMock) GetAngularVelocity() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *RotatableMock) SetDirection(direction int) (error) {
	args := m.Called(direction)
	return args.Error(0)
}

func (m *RotatableMock) GetDirectionsNumber() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

