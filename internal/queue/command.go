package queue

type ListenerStartCommand struct {
	listener *Listener
}

func (c *ListenerStartCommand) Execute() error {
	go c.listener.run()
	return nil
}

type ListenerSoftStopCommand struct {
	listener *Listener
}

func (c *ListenerSoftStopCommand) Execute() error {
	c.listener.SoftStop()
	return nil
}

type ListenerHardStopCommand struct {
	listener *Listener
}

func (c *ListenerHardStopCommand) Execute() error {
	c.listener.HardStop()
	return nil
}
