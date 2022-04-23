package command

import "modules/internal/core"

type Listener struct {
	queue        core.Queue
	errorHandler core.ErrorHandler
}

func (l *Listener) listen() {
	for {
		command, alive := l.queue.Get()
		if !alive {
			break
		}

		err := command.Execute()
		if err != nil {
			l.errorHandler(command, err)
		}
	}
}
