package command

import (
	"fmt"
	"log"
)

type LogFunc func(string)

type LogCommand struct {
	command Command
	err     error
	logFunc LogFunc
}

func StdLogFunc(message string) {
	log.Println(message)
}

func (c LogCommand) Execute() error {
	message := fmt.Sprintf("%s got error: '%s'", getType(c.command), c.err)
	c.logFunc(message)
	return nil
}

type RepeatCommand struct {
	command Command
	attempt int
}

func (c RepeatCommand) Execute() error {
	return c.command.Execute()
}
