package twitchbot

func NewGlobalCheckCommand() ChatCommandBuilder {
	return NewChatCommandBuilder("globalcheck").
		WithAlias("gc").
		WithHandler(func(ctx CommandContext) (msg string, err error) {
			return "Yes", nil
		})
}
