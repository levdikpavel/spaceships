package command

import "github.com/stretchr/testify/mock"

type QueueMock struct {
	mock.Mock
}

func (q *QueueMock) Get() (Command, bool) {
	args := q.Called()
	return args.Get(0).(Command), args.Bool(1)
}

func (q *QueueMock) Put(command Command) {
	q.Called(command)
}

type CommandMock struct {
	mock.Mock
}

func (c *CommandMock) Execute() error {
	args := c.Called()
	return args.Error(0)
}

type LoggerMock struct {
	message string
}

func (l *LoggerMock) Log(message string) {
	l.message = message
}
