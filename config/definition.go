package config

type BotConfig struct {
	Server   ServerConfig    `yaml:"server"`
	Twitch   TwitchConfig    `yaml:"twitch"`
	Commands []CommandConfig `yaml:"commands"`
}

type ServerConfig struct {
	TTL int64 `yaml:"ttl"`
}

type TwitchConfig struct {
	BotName     string `yaml:"botName"`
	ChannelName string `yaml:"channelName"`
	AccessToken string `yaml:"accessToken"`
}

type CommandConfig struct {
	Name     string `yaml:"name"`
	Enabled  *bool  `yaml:"enabled"`
	SubOnly  *bool  `yaml:"subOnly"`
	CoolDown *int   `yaml:"coolDown"`
}

func (cfg *BotConfig) GetCommandConfig(name string) *CommandConfig {
	for _, cmd := range cfg.Commands {
		if cmd.Name == name {
			return &cmd
		}
	}

	enabled := true
	subOnly := false
	coolDown := 60

	return &CommandConfig{
		Name:     name,
		Enabled:  &enabled,
		SubOnly:  &subOnly,
		CoolDown: &coolDown,
	}
}
