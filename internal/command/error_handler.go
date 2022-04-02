package command

type LogErrorHandler struct {
	queue Queue
}

func (h *LogErrorHandler) Handle(command Command, err error) {
	logCommand := LogCommand{
		command: command,
		err:     err,
	}
	h.queue.Put(logCommand)
}

type RepeatErrorHandler struct {
	queue          Queue
	attempts       int
	defaultHandler ErrorHandler
}

func (h *RepeatErrorHandler) Handle(command Command, err error) {
	repeatCommand, ok := command.(RepeatCommand)
	if !ok {
		h.queue.Put(RepeatCommand{
			command: command,
			attempt: 1,
		})
		return
	}

	if repeatCommand.attempt >= h.attempts {
		h.defaultHandler(repeatCommand.command, err)
		return
	}

	repeatCommand.attempt += 1
	h.queue.Put(repeatCommand)
}
