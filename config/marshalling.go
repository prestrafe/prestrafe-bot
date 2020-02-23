package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func ReadConfig(filename string) (*BotConfig, error) {
	configContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var botConfig BotConfig
	if err = yaml.Unmarshal(configContent, &botConfig); err != nil {
		return nil, err
	}

	return &botConfig, nil
}

func WriteConfig(filename string, config *BotConfig) error {
	configContent, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, configContent, 0644)
}
