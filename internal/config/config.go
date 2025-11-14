package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	MattermostURL string
	Username      string
	Password      string
}

// Load reads configuration from INI file and environment variables
// Priority: environment variables > config file > defaults
func Load() (*Config, error) {
	cfg := &Config{
		MattermostURL: "",
		Username:      "",
		Password:      "",
	}

	// Define config file search paths
	configPaths := []string{
		"./m2sh.ini",
	}
	
	if home, err := os.UserHomeDir(); err == nil {
		configPaths = append(configPaths,
			filepath.Join(home, ".m2sh.ini"),
			filepath.Join(home, ".config", "m2sh.ini"),
		)
	}
	configPaths = append(configPaths, "/etc/m2sh/m2sh.ini")

	// Try to load INI file from search paths
	var iniFile *ini.File
	var configFileFound bool
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			iniFile, err = ini.Load(path)
			if err != nil {
				return nil, fmt.Errorf("error reading config file %s: %w", path, err)
			}
			configFileFound = true
			break
		}
	}

	// Read values from INI file if found
	if configFileFound && iniFile != nil {
		section := iniFile.Section("")
		if section.HasKey("url") {
			cfg.MattermostURL = section.Key("url").String()
		}
		if section.HasKey("username") {
			cfg.Username = section.Key("username").String()
		}
		if section.HasKey("password") {
			cfg.Password = section.Key("password").String()
		}
	}

	// Use Viper for environment variables (higher priority than config file)
	v := viper.New()
	v.SetEnvPrefix("MM")
	v.AutomaticEnv()

	// Override with environment variables if set
	if v.IsSet("URL") {
		cfg.MattermostURL = v.GetString("URL")
	}
	if v.IsSet("USERNAME") {
		cfg.Username = v.GetString("USERNAME")
	}
	if v.IsSet("PASSWORD") {
		cfg.Password = v.GetString("PASSWORD")
	}

	return cfg, nil
}

// Validate checks if required configuration values are present
func (c *Config) Validate() error {
	if c.MattermostURL == "" {
		return fmt.Errorf("MM_URL (or url in config file) must be set")
	}
	if c.Username == "" {
		return fmt.Errorf("MM_USERNAME (or username in config file) must be set")
	}
	return nil
}
