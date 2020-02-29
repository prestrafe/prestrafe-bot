package twitchbot

type CommandContext interface {
	Parameter(name string) (value string, present bool)
}

type commandContext struct {
	parameters map[string]string
}

func (c *commandContext) Parameter(name string) (value string, present bool) {
	value, present = c.parameters[name]
	return
}
