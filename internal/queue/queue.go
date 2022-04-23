package queue

import "modules/internal/core"

type Queue struct {
	commandsChan chan core.Command
	aliveChan    chan int
}

func (q *Queue) Get() (core.Command, bool) {
	command, ok := <- q.commandsChan
	return command, ok
}

func (q *Queue) Put(command core.Command) {
	select {
	case <-q.aliveChan:
		return
	default:
	}

	q.commandsChan <- command
}
