package command

type Command interface {
	Execute() error
}

type ErrorHandler func(command Command, err error)

type Queue interface {
	Get() (Command, bool)
}
