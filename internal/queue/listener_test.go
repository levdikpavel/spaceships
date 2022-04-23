package queue

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"modules/internal/mock"
)

func TestListener(t *testing.T) {
	suite.Run(t, new(ListenerTestSuite))
}

type ListenerTestSuite struct {
	suite.Suite
}

func (s *ListenerTestSuite) TestStart() {
	listener := NewListener(1)
	queue := listener.GetQueue()
	command := mock.CommandMock{}
	mockChan := make(chan int)
	command.On("Execute").Return()
}
