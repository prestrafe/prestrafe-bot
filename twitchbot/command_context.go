package twitchbot

type CommandContext interface {
	Channel() string
	Parameter(name string) (value string, present bool)
}

type commandContext struct {
	channel    string
	parameters map[string]string
}

func (c *commandContext) Channel() string {
	return c.channel
}

func (c *commandContext) Parameter(name string) (value string, present bool) {
	value, present = c.parameters[name]
	return
}
