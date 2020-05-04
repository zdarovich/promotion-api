package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type (
	// Configuration struct
	Configuration struct {
		ReleaseMode        bool   `yaml:"releaseMode"`
		APIBasePath        string `yaml:"apiBasePath"`
		LogsEnabled        bool   `yaml:"logsEnabled"`
		LogFilePath        string `yaml:"logFilePath"`
		LogFileName        string `yaml:"logFileName"`
		LogRotateMegabytes int    `yaml:"logRotateMegaBytes"`
		LogRotateDuration  int    `yaml:"logRotateDuration"`
		LogRotateFiles     int    `yaml:"logRotateFiles"`
		LogDebugMode       bool   `yaml:"logDebugMode"`
		Port               int    `yaml:"port"`
		Database           struct {
			Discovery struct {
				Enabled bool   `yaml:"enabled"`
				Server  string `yaml:"server"`
				Timeout int    `yaml:"timeout"`
			} `yaml:"discovery"`
			Driver   string `yaml:"driver"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Name     string `yaml:"name"`
			Server   string `yaml:"server"`
			Port     int    `yaml:"port"`
		} `yaml:"database"`
		Redis struct {
			Server       string `yaml:"server"`
			DB           int    `yaml:"db"`
			Password     string `yaml:"password"`
			DialTimeout  int    `yaml:"dialTimeout"`
			ReadTimeout  int    `yaml:"readTimeout"`
			WriteTimeout int    `yaml:"writeTimeout"`
		} `yaml:"redis"`
		Identity struct {
			Server  string `yaml:"server"`
			Timeout int    `yaml:"timeout"`
			Token   string `yaml:"token"`
		} `yaml:"identity"`
	}
)

var configFileName string = "config.yml"

// Get reads configuration from the config file and
// returns the data structure
func Get() Configuration {

	configFile, err := os.Open(configFileName)

	if err != nil {
		// Application cannot start without the configuration
		// Thus it's ok to terminate the application startup here
		panic("Failed to read configuration file, please check that it exists")
	}

	var configuration Configuration
	decoder := yaml.NewDecoder(configFile)
	decoder.Decode(&configuration)

	return configuration
}
