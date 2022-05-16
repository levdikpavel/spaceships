package mock

import (
	"github.com/stretchr/testify/mock"

	"modules/internal/core"
	"modules/internal/vector"
)

type QueueMock struct {
	mock.Mock
}

func (q *QueueMock) Get() (core.Command, bool) {
	args := q.Called()
	return args.Get(0).(core.Command), args.Bool(1)
}

func (q *QueueMock) Put(command core.Command) {
	q.Called(command)
}

type CommandMock struct {
	mock.Mock
}

func (c *CommandMock) Execute() (err error) {
	args := c.Called()

	if execute, ok := args.Get(0).(func() error); ok {
		err = execute()
	} else {
		err = args.Error(0)
	}

	return err
}

type LoggerMock struct {
	Message string
}

func (l *LoggerMock) Log(message string) {
	l.Message = message
}

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

type AcceleratingMock struct {
	mock.Mock
}

func (m *AcceleratingMock) SetVelocity(v vector.Vector) error {
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

func (m *RotatableMock) SetDirection(value int) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *RotatableMock) GetDirectionsNumber() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

type FuelBurnableMock struct {
	mock.Mock
}

func (m *FuelBurnableMock) GetFuel() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *FuelBurnableMock) GetConsumption() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *FuelBurnableMock) SetFuel(value int) error {
	args := m.Called(value)
	return args.Error(0)
}

type MovableWithFuelMock struct {
	MovableMock
	FuelBurnableMock
}

type RotatableMovableMock struct {
	RotatableMock
	MovableMock
	AcceleratingMock
}

type UObjectMock struct {
	mock.Mock
}

func (o *UObjectMock) GetValue(valueName string) (interface{}, error) {
	args := o.Called(valueName)
	return args.Get(0), args.Error(1)
}

func (o *UObjectMock) SetValue(valueName string, newValue interface{}) error {
	args := o.Called(valueName, newValue)
	return args.Error(0)
}
