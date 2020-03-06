package config

// This configuration element focuses on aspects that configure the Game State integration server that is used to
// retrieve game information from streamers.
type GsiConfig struct {
	// Configures the address that the GSI server is listening on.
	Addr string `yaml:"addr"`
	// Configures the port that the GSI server is listening on.
	Port int `yaml:"port"`
}
