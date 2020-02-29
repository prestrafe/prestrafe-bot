package twitchbot

type CommandContext interface {
	GsiToken() string
	Parameter(name string) (value string, present bool)
}

type commandContext struct {
	gsiToken   string
	parameters map[string]string
}

func (c *commandContext) GsiToken() string {
	return c.gsiToken
}

func (c *commandContext) Parameter(name string) (value string, present bool) {
	value, present = c.parameters[name]
	return
}
