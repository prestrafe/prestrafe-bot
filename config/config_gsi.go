package config

// This configuration element focuses on aspects that configure the Game State integration server that is used to
// retrieve game information from streamers.
type GsiConfig struct {
	// Configures the port that the GSI server is listening on.
	Port int `yaml:"port"`
	// Configures the number of seconds a retrieved game state should be valid for, before it will be considered
	// out-dated.
	TTL int `yaml:"ttl"`
}
