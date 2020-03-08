package config

type ChannelConfig struct {
	Name                  string              `yaml:"name"`
	GsiToken              string              `yaml:"gsiToken"`
	JumpsOnlyWithoutBinds bool                `yaml:"jumpsOnlyWithoutBinds"`
	Commands              []ChatCommandConfig `yaml:"commands"`
}

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
