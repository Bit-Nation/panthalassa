package internal

type UpStream interface {
	Send(data string)
}

type Command struct {
	id          uint
	commandType string
}

func (c *Command) ID() uint {
	return c.id
}

func (c *Command) Type() string {
	return c.commandType
}

type Api struct {
	device UpStream
}

func (a *Api) SendToDevice(data string) {
	a.device.Send(data)
}
