package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"modules/internal/adapters"
	"modules/internal/core"
	"modules/internal/ioc"
	"modules/internal/mock"
)

func TestGenerate(t *testing.T) {
	a := assert.New(t)

	const inputFile = "../../core/interface.go"
	interfaces := []string{"Rotatable"}
	const expectedResult = `package adapters

import (
	"modules/internal/core"
	"modules/internal/ioc"
)

type RotatableAdapter struct{
   obj core.UObject
}

func NewRotatableAdapter(obj core.UObject) *RotatableAdapter {
	return &RotatableAdapter{
		obj: obj,
	}
}

func (a *RotatableAdapter) GetDirection() (int, error) {
	return ioc.Resolve("Operations.Rotatable:Direction.get", a.obj).(int), nil
}

func (a *RotatableAdapter) GetAngularVelocity() (int, error) {
	return ioc.Resolve("Operations.Rotatable:AngularVelocity.get", a.obj).(int), nil
}

func (a *RotatableAdapter) SetDirection(newValue int) (error) {
	return ioc.Resolve("Operations.Rotatable:Direction.set", a.obj, newValue).(core.Command).Execute()
}

func (a *RotatableAdapter) GetDirectionsNumber() (int, error) {
	return ioc.Resolve("Operations.Rotatable:DirectionsNumber.get", a.obj).(int), nil
}


func init() {
	_ = ioc.Resolve("IoC.Register", "Adapter",
		func(params ...interface{}) interface{} {
			adapterType := params[0].(string)
			switch adapterType {
			case "core.Rotatable":
				return NewRotatableAdapter(params[1].(core.UObject))
			default:
				panic("unknown adapter type" + adapterType)
			}
	}).(core.Command).Execute()
}
`
	sb := &strings.Builder{}
	err := generateAdapters(inputFile, sb, interfaces)

	a.NoError(err)
	a.Equal(expectedResult, sb.String())
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(RotatableAdapterTestSuite))
}

type RotatableAdapterTestSuite struct {
	suite.Suite

	adapter core.Rotatable
	obj     mock.UObjectMock
}

func (a *RotatableAdapterTestSuite) SetupTest() {
	a.obj = mock.UObjectMock{}
	a.adapter = adapters.NewRotatableAdapter(&a.obj)
}

func (a *RotatableAdapterTestSuite) TestGetDirection() {
	_ = ioc.Resolve("IoC.Register", "Operations.Rotatable:Direction.get", func(params ...interface{}) interface{} {
		return 5
	}).(core.Command).Execute()
	result, err := a.adapter.GetDirection()
	a.Require().NoError(err)
	a.Require().Equal(5, result)
}

func (a *RotatableAdapterTestSuite) TestGetAngularVelocity() {
	_ = ioc.Resolve("IoC.Register", "Operations.Rotatable:AngularVelocity.get", func(params ...interface{}) interface{} {
		return 5
	}).(core.Command).Execute()
	result, err := a.adapter.GetAngularVelocity()
	a.Require().NoError(err)
	a.Require().Equal(5, result)
}

func (a *RotatableAdapterTestSuite) TestSetDirection() {
	var direction int
	command := mock.CommandMock{}
	command.On("Execute").Return(func() error {
		direction = 5
		return nil
	})

	_ = ioc.Resolve("IoC.Register", "Operations.Rotatable:Direction.set", func(params ...interface{}) interface{} {
		return &command
	}).(core.Command).Execute()

	err := a.adapter.SetDirection(5)
	a.Require().NoError(err)
	a.Require().Equal(5, direction)
}

func (a *RotatableAdapterTestSuite) TestGetDirectionsNumber() {
	_ = ioc.Resolve("IoC.Register", "Operations.Rotatable:DirectionsNumber.get", func(params ...interface{}) interface{} {
		return 5
	}).(core.Command).Execute()
	result, err := a.adapter.GetDirectionsNumber()
	a.Require().NoError(err)
	a.Require().Equal(5, result)
}
