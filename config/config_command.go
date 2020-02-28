package config

// This config element configures properties of chat commands that might differ from channel to channel.
// Command configurations focus on aspects that impact the command itself and not the actual command handling.
type ChatCommandConfig struct {
	// The name of the command that should be configured. It has to be the name, as configuration by alias is not
	// supported. If the name is "*", the configuration should be applied as a default to all commands.
	Name string `yaml:"name"`
	// The enabled property can be used to enable or disable commands. Disabled commands are ignored by the bot.
	Enabled *bool `yaml:"enabled"`
	// The sub-only property can be used to control which users may trigger commands. Sub-only commands can only be
	// triggered by subscribers, moderators and broadcasters.
	SubOnly *bool `yaml:"subOnly"`
	// The cool down is the number of seconds that need to pass between two successful invocations of the same command.
	// Additional calls to the command will be ignored.
	CoolDown *int `yaml:"coolDown"`
}
