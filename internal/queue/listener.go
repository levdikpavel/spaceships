package queue

import (
	"context"

	"modules/internal/core"
)

func noOpErrorHandler(command core.Command, err error) {}

type Listener struct {
	ctx          context.Context
	cancel       func()
	commandsChan chan core.Command
	aliveChan    chan int
	queue        core.Queue
	errorHandler core.ErrorHandler
}

func NewListener(bufferLength int) *Listener {
	ctx, cancel := context.WithCancel(context.Background())
	commandsChan := make(chan core.Command, bufferLength)
	aliveChan := make(chan int)
	queue := &Queue{
		commandsChan: commandsChan,
		aliveChan:    aliveChan,
	}
	listener := &Listener{
		ctx:          ctx,
		cancel:       cancel,
		commandsChan: commandsChan,
		aliveChan:    aliveChan,
		queue:        queue,
		errorHandler: noOpErrorHandler,
	}

	return listener
}

func (l *Listener) GetQueue() core.Queue {
	return l.queue
}

func (l *Listener) SetErrorHandler(errorHandler core.ErrorHandler) {
	l.errorHandler = errorHandler
}

func (l *Listener) run() {
	for {
		select {
		case <-l.ctx.Done():
			// hard stop
			return
		case command, ok := <-l.commandsChan:
			if !ok {
				// soft stop
				return
			}

			err := command.Execute()
			if err != nil {
				l.errorHandler(command, err)
			}

			if l.ctx.Err() != nil {
				// hard stop
				return
			}
		}
	}
}

func (l *Listener) SoftStop() {
	close(l.aliveChan)
	close(l.commandsChan)
}

func (l *Listener) HardStop() {
	close(l.aliveChan)
	l.cancel()
}

func (l *Listener) StartCommand() core.Command {
	return &ListenerStartCommand{
		listener: l,
	}
}

func (l *Listener) SoftStopCommand() core.Command {
	return &ListenerSoftStopCommand{
		listener: l,
	}
}

func (l *Listener) HardStopCommand() core.Command {
	return &ListenerHardStopCommand{
		listener: l,
	}
}
