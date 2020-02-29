package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// The bot configuration containers all required settings to operate the bot. It can be read from and written to a YAML
// format.
type BotConfig struct {
	// Contains settings that configure the internal Game State Integration Server.
	Gsi GsiConfig `yaml:"gsi"`
	// Contains settings that configure the integration into the Twitch API for the chat bot.
	Twitch TwitchConfig `yaml:"twitch"`
	// Contains a list of chat channel settings, that configure different channels that the bot operates in.
	Channels []ChannelConfig `yaml:"channels"`
}

// Reads a YAML file that may contain a configuration. If the configuration is not valid, an error is returned instead.
func ReadConfig(filename string) (*BotConfig, error) {
	configContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	botConfig := new(BotConfig)
	if err = yaml.Unmarshal(configContent, botConfig); err != nil {
		return nil, err
	}

	return botConfig, nil
}

// Writes the configuration to a YAML file.
func (config *BotConfig) WriteConfig(filename string) error {
	configContent, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, configContent, 0644)
}
