package command

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestCommand(t *testing.T) {
	suite.Run(t, new(CommandSuite))
}

type CommandSuite struct {
	suite.Suite

	command CommandMock
	err     error
}

func (s *CommandSuite) SetupSuite() {
	s.err = fmt.Errorf("some error")
}

func (s *CommandSuite) TestLog() {
	loggerMock := LoggerMock{}
	command := LogCommand{
		command: &s.command,
		err:     s.err,
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

	s.command.On("Execute").Return(s.err)
	err := command.Execute()
	s.Require().Error(err)
	s.Require().Equal(s.err, err)
	s.command.AssertExpectations(s.T())
}
