package command

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestErrorHandler(t *testing.T) {
	suite.Run(t, new(ErrorHandlerSuite))
}

type ErrorHandlerSuite struct {
	suite.Suite

	queue   QueueMock
	err     error
	command *CommandMock
}

func (s *ErrorHandlerSuite) SetupTest() {
	s.queue = QueueMock{}
	s.err = fmt.Errorf("some error")
	s.command = &CommandMock{}
}

func (s *ErrorHandlerSuite) TestLog() {
	h := LogErrorHandler{
		queue: &s.queue,
	}

	logCommand := LogCommand{
		command: s.command,
		err:     s.err,
	}
	s.queue.On("Put", logCommand).Return()

	h.Handle(s.command, s.err)
	s.queue.AssertExpectations(s.T())
}

func (s *ErrorHandlerSuite) TestRepeat() {
	h := RepeatErrorHandler{
		queue: &s.queue,
	}

	repeatCommand := RepeatCommand{
		command: s.command,
		attempt: 0,
	}
	s.queue.On("Put", repeatCommand).Return()

	h.Handle(s.command, s.err)
	s.queue.AssertExpectations(s.T())
}

func (s *ErrorHandlerSuite) TestRepeatDouble() {
	logHandler := LogErrorHandler{
		queue: &s.queue,
	}

	h := RepeatErrorHandler{
		queue:          &s.queue,
		attempts:       1,
		defaultHandler: logHandler.Handle,
	}

	repeatCommand := RepeatCommand{
		command: s.command,
		attempt: 0,
	}
	s.queue.On("Put", repeatCommand).Return()

	h.Handle(s.command, s.err)

	logCommand := LogCommand{
		command: s.command,
		err:     s.err,
	}
	s.queue.On("Put", logCommand).Return()
	repeatCommand.attempt = 1
	h.Handle(repeatCommand, s.err)

	s.queue.AssertExpectations(s.T())
}
