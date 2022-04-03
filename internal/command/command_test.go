package command

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"modules/internal/vector"
)

var (
	errSomeError = fmt.Errorf("some error")	
)

func TestCommand(t *testing.T) {
	suite.Run(t, new(CommandSuite))
}

type CommandSuite struct {
	suite.Suite

	command CommandMock
}

func (s *CommandSuite) TestLog() {
	loggerMock := LoggerMock{}
	command := LogCommand{
		command: &s.command,
		err:     errSomeError,
		logFunc: loggerMock.Log,
	}

	err := command.Execute()
	s.Require().NoError(err)
	s.Require().Equal("CommandMock got error: 'some error'", loggerMock.message)
}

func (s *CommandSuite) TestRepeat() {
	command := RepeatCommand{
		command: &s.command,
	}

	s.command.On("Execute").Return(errSomeError)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().Equal(errSomeError, err)
	s.command.AssertExpectations(s.T())
}


func TestMove(t *testing.T) {
	suite.Run(t, new(MoveTestSuite))
}

type MoveTestSuite struct {
	suite.Suite

	mock      MovableMock
	nilVector vector.Vector
}

func (s *MoveTestSuite) SetupTest() {
	s.mock = MovableMock{}
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
	move := NewMoveCommand(&s.mock)
	err := move.Execute()
	s.Require().NoError(err)
}

func (s *MoveTestSuite) TestPosError() {
	s.mock.On("GetPosition").Return(s.nilVector, errSomeError)
	move := NewMoveCommand(&s.mock)
	err := move.Execute()
	s.Require().Error(err)
}

func (s *MoveTestSuite) TestVelocityError() {
	pos := vector.New([]int{12, 5})
	s.mock.On("GetPosition").Return(pos, nil).
		On("GetVelocity").Return(s.nilVector, errSomeError)
	move := NewMoveCommand(&s.mock)
	err := move.Execute()
	s.Require().Error(err)
}

func (s *MoveTestSuite) TestSetPosError() {
	pos := vector.New([]int{12, 5})
	v := vector.New([]int{-7, 3})
	posNew := vector.New([]int{5, 8})
	s.mock.On("GetPosition").Return(pos, nil).
		On("GetVelocity").Return(v, nil).
		On("SetPosition", posNew).Return(errSomeError)
	move := NewMoveCommand(&s.mock)
	err := move.Execute()
	s.Require().Error(err)
}

func TestRotate(t *testing.T) {
	suite.Run(t, new(RotateTestSuite))
}

type RotateTestSuite struct {
	suite.Suite

	mock RotatableMock
}

func (s *RotateTestSuite) SetupTest() {
	s.mock = RotatableMock{}
}

func (s *RotateTestSuite) TearDownTest() {
	s.mock.AssertExpectations(s.T())
}

func (s *RotateTestSuite) TestSuccess() {
	s.mock.On("GetDirection").Return(300, nil).
		On("GetAngularVelocity").Return(70, nil).
		On("GetDirectionsNumber").Return(360, nil).
		On("SetDirection", 10).Return(nil)
	rotate := NewRotateCommand(&s.mock)
	err := rotate.Execute()
	s.Require().NoError(err)
}

func (s *RotateTestSuite) TestDirectionError() {
	s.mock.On("GetDirection").Return(0, errSomeError)
	rotate := NewRotateCommand(&s.mock)
	err := rotate.Execute()
	s.Require().Error(err)
}

func (s *RotateTestSuite) TestAngularVelocityError() {
	s.mock.On("GetDirection").Return(300, nil).
		On("GetAngularVelocity").Return(0, errSomeError)
	rotate := NewRotateCommand(&s.mock)
	err := rotate.Execute()
	s.Require().Error(err)
}

func (s *RotateTestSuite) TestDirectionNumberError() {
	s.mock.On("GetDirection").Return(300, nil).
		On("GetAngularVelocity").Return(70, nil).
		On("GetDirectionsNumber").Return(0, errSomeError)
	rotate := NewRotateCommand(&s.mock)
	err := rotate.Execute()
	s.Require().Error(err)
}

func (s *RotateTestSuite) TestSetDirectionError() {
	s.mock.On("GetDirection").Return(300, nil).
		On("GetAngularVelocity").Return(70, nil).
		On("GetDirectionsNumber").Return(360, nil).
		On("SetDirection", 10).Return(errSomeError)
	rotate := NewRotateCommand(&s.mock)
	err := rotate.Execute()
	s.Require().Error(err)
}

func TestCheckFuel(t *testing.T) {
	suite.Run(t, new(CheckFuelTestSuite))
}

type CheckFuelTestSuite struct {
	suite.Suite

	mock FuelBurnableMock
}

func (s *CheckFuelTestSuite) SetupTest() {
	s.mock = FuelBurnableMock{}
}

func (s *CheckFuelTestSuite) TearDownTest() {
	s.mock.AssertExpectations(s.T())
}

func (s *CheckFuelTestSuite) TestSuccess() {
	s.mock.On("GetFuel").Return(300, nil).
		On("GetConsumption").Return(70, nil)
	command := NewCheckFuelCommand(&s.mock)
	err := command.Execute()
	s.Require().NoError(err)
}

func (s *CheckFuelTestSuite) TestGetFuelError() {
	s.mock.On("GetFuel").Return(0, errSomeError)
	command := NewCheckFuelCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, errSomeError)
}

func (s *CheckFuelTestSuite) TestGetConsumptionError() {
	s.mock.On("GetFuel").Return(300, nil).
		On("GetConsumption").Return(0, errSomeError)
	command := NewCheckFuelCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, errSomeError)
}

func (s *CheckFuelTestSuite) TestNotEnough() {
	s.mock.On("GetFuel").Return(300, nil).
		On("GetConsumption").Return(370, nil)
	command := NewCheckFuelCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrNotEnoughFuel)
}

func TestBurnFuel(t *testing.T) {
	suite.Run(t, new(BurnFuelTestSuite))
}

type BurnFuelTestSuite struct {
	suite.Suite

	mock FuelBurnableMock
}

func (s *BurnFuelTestSuite) SetupTest() {
	s.mock = FuelBurnableMock{}
}

func (s *BurnFuelTestSuite) TearDownTest() {
	s.mock.AssertExpectations(s.T())
}

func (s *BurnFuelTestSuite) TestSuccess() {
	s.mock.On("GetFuel").Return(300, nil).
		On("GetConsumption").Return(70, nil).
		On("SetFuel", 230).Return(nil)
	command := NewBurnFuelCommand(&s.mock)
	err := command.Execute()
	s.Require().NoError(err)
}

func (s *BurnFuelTestSuite) TestGetFuelError() {
	s.mock.On("GetFuel").Return(0, errSomeError)
	command := NewBurnFuelCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, errSomeError)
}

func (s *BurnFuelTestSuite) TestGetConsumptionError() {
	s.mock.On("GetFuel").Return(300, nil).
		On("GetConsumption").Return(0, errSomeError)
	command := NewBurnFuelCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, errSomeError)
}

func (s *BurnFuelTestSuite) TestSetFuelError() {
	s.mock.On("GetFuel").Return(300, nil).
		On("GetConsumption").Return(70, nil).
		On("SetFuel", 230).Return(errSomeError)
	command := NewBurnFuelCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, errSomeError)
}

func (s *BurnFuelTestSuite) TestNotEnough() {
	s.mock.On("GetFuel").Return(300, nil).
		On("GetConsumption").Return(370, nil)
	command := NewBurnFuelCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrNotEnoughFuel)
}

func TestMacroCommand(t *testing.T) {
	suite.Run(t, new(MacroCommandTestSuite))
}

type MacroCommandTestSuite struct {
	suite.Suite

	mock1 CommandMock
	mock2 CommandMock
}

func (s *MacroCommandTestSuite) SetupTest() {
	s.mock1 = CommandMock{}
	s.mock2 = CommandMock{}
}

func (s *MacroCommandTestSuite) TearDownTest() {
	s.mock1.AssertExpectations(s.T())
	s.mock2.AssertExpectations(s.T())
}

func (s *MacroCommandTestSuite) TestSuccess() {
	s.mock1.On("Execute").Return(nil)
	s.mock2.On("Execute").Return(nil)
	command := NewMacroCommand(&s.mock1, &s.mock2)
	err := command.Execute()
	s.Require().NoError(err)
}

func (s *MacroCommandTestSuite) TestError1() {
	s.mock1.On("Execute").Return(errSomeError)
	command := NewMacroCommand(&s.mock1, &s.mock2)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, errSomeError)
}

func (s *MacroCommandTestSuite) TestError2() {
	s.mock1.On("Execute").Return(nil)
	s.mock2.On("Execute").Return(errSomeError)
	command := NewMacroCommand(&s.mock1, &s.mock2)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, errSomeError)
}

func TestMoveWithFuel(t *testing.T) {
	suite.Run(t, new(MoveWithFuelTestSuite))
}

type MoveWithFuelTestSuite struct {
	suite.Suite

	mock        MovableWithFuelMock
	position    vector.Vector
	velocity    vector.Vector
	positionNew vector.Vector
	nilVector   vector.Vector
}

func (s *MoveWithFuelTestSuite) SetupTest() {
	s.mock = MovableWithFuelMock{}
	s.position = vector.New([]int{12, 5})
	s.velocity = vector.New([]int{-7, 3})
	s.positionNew = vector.New([]int{5, 8})
}

func (s *MoveWithFuelTestSuite) TearDownTest() {
	s.mock.MovableMock.AssertExpectations(s.T())
	s.mock.FuelBurnableMock.AssertExpectations(s.T())
}

func (s *MoveWithFuelTestSuite) TestSuccess() {
	s.mock.MovableMock.On("GetPosition").Return(s.position, nil).
		On("GetVelocity").Return(s.velocity, nil).
		On("SetPosition", s.positionNew).Return(nil)
	s.mock.FuelBurnableMock.On("GetFuel").Return(300, nil).
		On("GetConsumption").Return(70, nil).
		On("SetFuel", 230).Return(nil)
	command := NewMoveWithFuelCommand(&s.mock)
	err := command.Execute()
	s.Require().NoError(err)
}

func (s *MoveWithFuelTestSuite) TestCheckError() {
	s.mock.FuelBurnableMock.On("GetFuel").Return(300, nil).
		On("GetConsumption").Return(0, errSomeError)
	command := NewMoveWithFuelCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, errSomeError)
}

func (s *MoveWithFuelTestSuite) TestMoveError() {
	s.mock.MovableMock.On("GetPosition").Return(s.nilVector, errSomeError)
	s.mock.FuelBurnableMock.On("GetFuel").Return(300, nil).
		On("GetConsumption").Return(70, nil)
	command := NewMoveWithFuelCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, errSomeError)
}

func (s *MoveWithFuelTestSuite) TestBurnError() {
	s.mock.MovableMock.On("GetPosition").Return(s.position, nil).
		On("GetVelocity").Return(s.velocity, nil).
		On("SetPosition", s.positionNew).Return(nil)
	s.mock.FuelBurnableMock.On("GetFuel").Return(300, nil).
		On("GetConsumption").Return(70, nil).
		On("SetFuel", 230).Return(errSomeError)
	command := NewMoveWithFuelCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, errSomeError)
}

func TestTurnVelocity(t *testing.T) {
	suite.Run(t, new(TurnVelocityTestSuite))
}

type TurnVelocityTestSuite struct {
	suite.Suite

	mock             RotatableMovableMock
	velocity         vector.Vector
	velocity3D       vector.Vector
	velocityNew      vector.Vector
	nilVector        vector.Vector
	direction        int
	angularVelocity  int
	directionsNumber int
}

func (s *TurnVelocityTestSuite) SetupTest() {
	s.mock = RotatableMovableMock{}
	s.velocity = vector.New([]int{100, 10})
	s.velocity3D = vector.New([]int{100, 10, 20})
	s.direction = 6
	s.angularVelocity = 4
	s.directionsNumber = 360
	s.velocityNew = vector.New([]int{98, 17})
}

func (s *TurnVelocityTestSuite) TearDownTest() {
	s.mock.RotatableMock.AssertExpectations(s.T())
	s.mock.MovableMock.AssertExpectations(s.T())
	s.mock.AcceleratingMock.AssertExpectations(s.T())
}

func (s *TurnVelocityTestSuite) TestSuccess() {
	s.mock.MovableMock.On("GetVelocity").Return(s.velocity, nil)
	s.mock.RotatableMock.On("GetDirection").Return(s.direction, nil).
		On("GetAngularVelocity").Return(s.angularVelocity, nil).
		On("GetDirectionsNumber").Return(s.directionsNumber, nil)
	s.mock.AcceleratingMock.On("SetVelocity", s.velocityNew).Return(nil)
	command := NewTurnVelocityCommand(&s.mock)
	err := command.Execute()
	s.Require().NoError(err)
}

func (s *TurnVelocityTestSuite) TestGetVelocityError() {
	s.mock.MovableMock.On("GetVelocity").Return(s.nilVector, errSomeError)
	command := NewTurnVelocityCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, errSomeError)
}

func (s *TurnVelocityTestSuite) TestVelocity3D() {
	s.mock.MovableMock.On("GetVelocity").Return(s.velocity3D, nil)
	command := NewTurnVelocityCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrUnsupportedDimension)
}

func (s *TurnVelocityTestSuite) TestGetDirectionError() {
	s.mock.MovableMock.On("GetVelocity").Return(s.velocity, nil)
	s.mock.RotatableMock.On("GetDirection").Return(0, errSomeError)
	command := NewTurnVelocityCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, errSomeError)
}

func (s *TurnVelocityTestSuite) TestSetVelocityError() {
	s.mock.MovableMock.On("GetVelocity").Return(s.velocity, nil)
	s.mock.RotatableMock.On("GetDirection").Return(s.direction, nil).
		On("GetAngularVelocity").Return(s.angularVelocity, nil).
		On("GetDirectionsNumber").Return(s.directionsNumber, nil)
	s.mock.AcceleratingMock.On("SetVelocity", s.velocityNew).Return(errSomeError)
	command := NewTurnVelocityCommand(&s.mock)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().ErrorIs(err, errSomeError)
}

func TestRotateWithVelocity(t *testing.T) {
	suite.Run(t, new(RotateWithVelocityTestSuite))
}

type RotateWithVelocityTestSuite struct {
	suite.Suite

	mock             RotatableMovableMock
	velocity         vector.Vector
	velocityNew      vector.Vector
	direction        int
	angularVelocity  int
	directionsNumber int
	directionNew     int
}

func (s *RotateWithVelocityTestSuite) SetupTest() {
	s.mock = RotatableMovableMock{}
	s.velocity = vector.New([]int{100, 10})
	s.direction = 6
	s.angularVelocity = 4
	s.directionsNumber = 360
	s.directionNew = 10
	s.velocityNew = vector.New([]int{98, 17})
}

func (s *RotateWithVelocityTestSuite) TearDownTest() {
	s.mock.RotatableMock.AssertExpectations(s.T())
	s.mock.MovableMock.AssertExpectations(s.T())
	s.mock.AcceleratingMock.AssertExpectations(s.T())
}

func (s *RotateWithVelocityTestSuite) TestSuccess() {
	s.mock.MovableMock.On("GetVelocity").Return(s.velocity, nil)
	s.mock.RotatableMock.On("GetDirection").Return(s.direction, nil).
		On("GetAngularVelocity").Return(s.angularVelocity, nil).
		On("GetDirectionsNumber").Return(s.directionsNumber, nil).
		On("SetDirection", s.directionNew).Return(nil)
	s.mock.AcceleratingMock.On("SetVelocity", s.velocityNew).Return(nil)
	command := NewRotateWithVelocityCommand(&s.mock)
	err := command.Execute()
	s.Require().NoError(err)
}

func (s *RotateWithVelocityTestSuite) TestSuccessOnlyRotate() {
	s.mock.RotatableMock.On("GetDirection").Return(s.direction, nil).
		On("GetAngularVelocity").Return(s.angularVelocity, nil).
		On("GetDirectionsNumber").Return(s.directionsNumber, nil).
		On("SetDirection", s.directionNew).Return(nil)
	command := NewRotateWithVelocityCommand(&s.mock.RotatableMock)
	err := command.Execute()
	s.Require().NoError(err)
}
