package command

type CompositeErrorHandler struct {
	handlers       map[string]ErrorHandler
	defaultHandler ErrorHandler
}

func (h *CompositeErrorHandler) Handle(command Command, err error) {

}

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
			attempt: 0,
		})
		return
	}

	if repeatCommand.attempt >= h.attempts {
		h.defaultHandler(command, err)
	}

	repeatCommand.attempt += 1
	h.queue.Put(repeatCommand)
}
