package config

// This config element configures properties of a twitch channel that is observed by the chat bot. It focuses on
// settings that may very between different streamers and channels.
type ChannelConfig struct {
	Name     string              `yaml:"name"`
	GsiToken string              `yaml:"gsiToken"`
	Commands []ChatCommandConfig `yaml:"commands"`
}

// A helper method that fetches a command configuration from a channel configuration with a given name. If no such
// command config is found, an empty chat configuration element will be returned.
func (cfg *ChannelConfig) GetCommandConfig(name string) *ChatCommandConfig {
	for _, cmd := range cfg.Commands {
		if cmd.Name == name {
			return &cmd
		}
	}

	return &ChatCommandConfig{
		Name:     name,
		Enabled:  nil,
		SubOnly:  nil,
		CoolDown: nil,
	}
}
