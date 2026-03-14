package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config holds all non-secret CLI configuration.
type Config struct {
	Profile        string        `mapstructure:"profile"`
	Pretty         bool          `mapstructure:"pretty"`
	Verbose        bool          `mapstructure:"verbose"`
	Timeout        time.Duration `mapstructure:"timeout"`
	MaxRetries     int           `mapstructure:"max_retries"`
	RetryMutations   bool          `mapstructure:"retry_mutations"`
	PairingServiceURL string        `mapstructure:"pairing_service_url"`
}

// Load reads configuration from file, environment, and defaults.
// Precedence: env vars > config file > defaults.
func Load() Config {
	v := viper.New()

	// Defaults
	v.SetDefault("profile", "default")
	v.SetDefault("pretty", false)
	v.SetDefault("verbose", false)
	v.SetDefault("timeout", "15s")
	v.SetDefault("max_retries", 3)
	v.SetDefault("retry_mutations", false)
	v.SetDefault("pairing_service_url", "https://trello-connector-production.up.railway.app")

	// Environment binding
	v.SetEnvPrefix("TRELLO")
	v.AutomaticEnv()

	// Config file
	configPath := os.Getenv("TRELLO_CONFIG_PATH")
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			configPath = filepath.Join(home, ".config", "trello-cli", "config.yaml")
		}
	}
	if configPath != "" {
		v.SetConfigFile(configPath)
		_ = v.ReadInConfig() // Ignore error — missing file is fine
	}

	var cfg Config
	_ = v.Unmarshal(&cfg)

	// Parse timeout string into duration
	if cfg.Timeout == 0 {
		dur, err := time.ParseDuration(v.GetString("timeout"))
		if err == nil {
			cfg.Timeout = dur
		} else {
			cfg.Timeout = 15 * time.Second
		}
	}

	return cfg
}
