package command

import (
	"log"
)

type LogCommand struct {
	command Command
	err     error
}

func (c LogCommand) Execute() error {
	log.Println(getType(c.command), "got error", c.err)
	return nil
}

type RepeatCommand struct {
	command Command
	attempt int
}

func (c RepeatCommand) Execute() error {
	return c.command.Execute()
}
