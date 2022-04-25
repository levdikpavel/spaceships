package queue

import (
	"runtime"
	"testing"
	"time"

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

	command := mock.CommandMock{}
	executeChan := make(chan int)
	executeStarted := false
	command.On("Execute").Return(func() error {
		executeStarted = true
		close(executeChan)
		return nil
	})

	listener.GetQueue().Put(&command)
	s.Require().False(executeStarted)

	err := listener.StartCommand().Execute()
	s.Require().NoError(err)

	<- executeChan
	s.Require().True(executeStarted)
}

func (s *ListenerTestSuite) TestHardStop() {
	listener := NewListener(1)

	command1 := mock.CommandMock{}
	executeChan1Wait := make(chan int)
	executeChan1Go := make(chan int)
	command1.On("Execute").Return(func() error {
		close(executeChan1Wait)
		<- executeChan1Go
		return nil
	})

	executeChan2 := make(chan int)
	executeStarted2 := false
	command2 := mock.CommandMock{}
	command2.On("Execute").Return(func() error {
		executeStarted2 = true
		close(executeChan2)
		return nil
	})

	err := listener.StartCommand().Execute()
	s.Require().NoError(err)

	listener.GetQueue().Put(&command1)
	<-executeChan1Wait
	listener.GetQueue().Put(&command2)

	err = listener.HardStopCommand().Execute()
	s.Require().NoError(err)
	s.Require().Error(listener.ctx.Err())

	close(executeChan1Go)

	runtime.Gosched()
	select {
	case <-executeChan2:
		// should not get here
		s.Require().False(true)
	case <-time.After(time.Millisecond * 10):
	}

	s.Require().False(executeStarted2)
}

func (s *ListenerTestSuite) TestSoftStop() {
	listener := NewListener(1)

	command1 := mock.CommandMock{}
	executeChan1Wait := make(chan int)
	executeChan1Go := make(chan int)
	command1.On("Execute").Return(func() error {
		close(executeChan1Wait)
		<- executeChan1Go
		return nil
	})

	executeChan2 := make(chan int)
	executeStarted2 := false
	command2 := mock.CommandMock{}
	command2.On("Execute").Return(func() error {
		executeStarted2 = true
		close(executeChan2)
		return nil
	})

	err := listener.StartCommand().Execute()
	s.Require().NoError(err)

	listener.GetQueue().Put(&command1)
	<-executeChan1Wait
	listener.GetQueue().Put(&command2)

	err = listener.SoftStopCommand().Execute()
	s.Require().NoError(err)

	close(executeChan1Go)

	<-executeChan2

	s.Require().True(executeStarted2)
}
