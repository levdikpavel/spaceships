package command

import "modules/internal/core"

type LogErrorHandler struct {
	queue   core.Queue
	logFunc LogFunc
}

func (h *LogErrorHandler) Handle(command core.Command, err error) {
	logCommand := LogCommand{
		command: command,
		err:     err,
		logFunc: h.logFunc,
	}
	h.queue.Put(logCommand)
}

type RepeatErrorHandler struct {
	queue          core.Queue
	attempts       int
	defaultHandler core.ErrorHandler
}

func (h *RepeatErrorHandler) Handle(command core.Command, err error) {
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
