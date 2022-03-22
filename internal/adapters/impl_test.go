package adapters

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"modules/internal/vector"
)

func TestMove(t *testing.T) {
	suite.Run(t, new(MoveTestSuite))
}

type MoveTestSuite struct {
	suite.Suite

	mock      MovableMock
	err       error
	nilVector vector.Vector
}

func (s *MoveTestSuite) SetupTest() {
	s.mock = MovableMock{}
	s.err = fmt.Errorf("some error")
}

func (s *MoveTestSuite) TearDownTest() {
	s.mock.AssertExpectations(s.T())
}

func (s *MoveTestSuite) TestSuccess() {
	pos := vector.New([]int{12, 5})
	v := vector.New([]int{-7, 3})
	posNew := vector.New([]int{5, 8})
	s.mock.On("GetPosition").Return(pos, nil).
		On("GetVelocity").Return(v, nil).
		On("SetPosition", posNew).Return(nil)
	move := NewMove(&s.mock)
	err := move.Execute()
	s.Require().NoError(err)
}

func (s *MoveTestSuite) TestPosError() {
	s.mock.On("GetPosition").Return(s.nilVector, s.err)
	move := NewMove(&s.mock)
	err := move.Execute()
	s.Require().Error(err)
}

func (s *MoveTestSuite) TestVelocityError() {
	pos := vector.New([]int{12, 5})
	s.mock.On("GetPosition").Return(pos, nil).
		On("GetVelocity").Return(s.nilVector, s.err)
	move := NewMove(&s.mock)
	err := move.Execute()
	s.Require().Error(err)
}

func (s *MoveTestSuite) TestSetPosError() {
	pos := vector.New([]int{12, 5})
	v := vector.New([]int{-7, 3})
	posNew := vector.New([]int{5, 8})
	s.mock.On("GetPosition").Return(pos, nil).
		On("GetVelocity").Return(v, nil).
		On("SetPosition", posNew).Return(s.err)
	move := NewMove(&s.mock)
	err := move.Execute()
	s.Require().Error(err)
}

func TestRotate(t *testing.T) {
	suite.Run(t, new(RotateTestSuite))
}

type RotateTestSuite struct {
	suite.Suite

	mock RotatableMock
	err  error
}

func (s *RotateTestSuite) SetupTest() {
	s.mock = RotatableMock{}
	s.err = fmt.Errorf("some error")
}

func (s *RotateTestSuite) TearDownTest() {
	s.mock.AssertExpectations(s.T())
}

func (s *RotateTestSuite) TestSuccess() {
	s.mock.On("GetDirection").Return(300, nil).
		On("GetAngularVelocity").Return(70, nil).
		On("GetDirectionsNumber").Return(360, nil).
		On("SetDirection", 10).Return(nil)
	rotate := NewRotate(&s.mock)
	err := rotate.Execute()
	s.Require().NoError(err)
}

func (s *RotateTestSuite) TestDirectionError() {
	s.mock.On("GetDirection").Return(0, s.err)
	rotate := NewRotate(&s.mock)
	err := rotate.Execute()
	s.Require().Error(err)
}

func (s *RotateTestSuite) TestAngularVelocityError() {
	s.mock.On("GetDirection").Return(300, nil).
		On("GetAngularVelocity").Return(0, s.err)
	rotate := NewRotate(&s.mock)
	err := rotate.Execute()
	s.Require().Error(err)
}

func (s *RotateTestSuite) TestDirectionNumberError() {
	s.mock.On("GetDirection").Return(300, nil).
		On("GetAngularVelocity").Return(70, nil).
		On("GetDirectionsNumber").Return(0, s.err)
	rotate := NewRotate(&s.mock)
	err := rotate.Execute()
	s.Require().Error(err)
}

func (s *RotateTestSuite) TestSetDirectionError() {
	s.mock.On("GetDirection").Return(300, nil).
		On("GetAngularVelocity").Return(70, nil).
		On("GetDirectionsNumber").Return(360, nil).
		On("SetDirection", 10).Return(s.err)
	rotate := NewRotate(&s.mock)
	err := rotate.Execute()
	s.Require().Error(err)
}
