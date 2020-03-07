package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ChannelsConfig struct {
	Channels []ChannelConfig `yaml:"channels"`
}

func ReadConfig(filename string) (*ChannelsConfig, error) {
	configContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	channelsConfig := new(ChannelsConfig)
	if err = yaml.Unmarshal(configContent, channelsConfig); err != nil {
		return nil, err
	}

	return channelsConfig, nil
}
