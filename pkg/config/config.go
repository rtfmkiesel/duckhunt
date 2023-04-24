package config

import (
	"duckhunt/pkg/file"
	"duckhunt/pkg/logger"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// Config file name
	defaultConfigFile string = "duckhunt.yml"
)

// Struct for the config
type Cfg struct {
	Alert         bool   `yaml:"alert"`
	AlertOnTop    bool   `yaml:"alertontop"`
	AlertTitle    string `yaml:"alerttitle"`
	AlertMsg      string `yaml:"alertmessage"`
	BlockDuration int    `yaml:"blockduration"`
	LogFile       string `yaml:"logfile"`
	MaxInterval   int64  `yaml:"maxInterval"`
	HistorySize   int    `yaml:"historySize"`
	IgnoredKeys   []int  `yaml:"ignoredKeys"`
}

// LoadConfig() will load the config from a yml file into a struct
func LoadConfig() (cfg Cfg, err error) {
	// Set the full path for the config file
	var configFile = defaultConfigFile
	exePath, err := file.ExeDir()
	if err != nil {
		return cfg, err
	}
	configFile = filepath.Join(exePath + "/" + configFile)

	// Read config file
	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		return cfg, fmt.Errorf("reading config file '%s' failed with: %s", configFile, err)
	}

	// Parse to yml
	err = yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("unmarshaling config file '%s' failed with: %s", configFile, err)
	}

	// Set full path for the log file
	cfg.LogFile = filepath.Join(exePath + "/" + cfg.LogFile)

	// Setup the log file
	err = logger.LogInit(cfg.LogFile)
	if err != nil {
		return cfg, err
	}
	logger.LogWrite("Config file loaded")
	return cfg, nil
}
