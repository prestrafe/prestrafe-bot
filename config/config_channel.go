package config

type ChannelConfig struct {
	Name        string `yaml:"name"`
	GsiToken    string `yaml:"gsiToken"`
	ServerToken string `yaml:"serverToken"`
}
