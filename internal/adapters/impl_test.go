package adapters

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"modules/internal/vector"
)

func TestMove(t *testing.T) {
	suite.Run(t, new(MoveTestSuite))
}

type MoveTestSuite struct {
	suite.Suite
}

func (s *MoveTestSuite) Test1() {
	pos := vector.New([]int{12, 5})
	v := vector.New([]int{-7, 3})
	posNew := vector.New([]int{5, 8})
	var mock MovableMock
	mock.On("GetPosition").Return(pos).
		On("GetVelocity").Return(v).
		On("SetPosition", posNew).Return()
	move := NewMove(&mock)
	move.Execute()
	mock.AssertExpectations(s.T())
}
