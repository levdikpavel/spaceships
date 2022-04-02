package command

type Listener struct {
	queue        Queue
	errorHandler ErrorHandler
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
