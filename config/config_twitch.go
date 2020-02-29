package config

// This configuration element bundles settings that are required by Twitch, so that the bot can communicate with the
// Twitch API correctly.
type TwitchConfig struct {
	// The user name that the bot user should have.
	UserName string `yaml:"userName"`
	// The access token of the bot user that is used to authenticate against the Twitch API.
	AccessToken string `yaml:"accessToken"`
}
