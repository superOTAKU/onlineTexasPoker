package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	Server ServerConfig
	Log    LogConfig
}

type ServerConfig struct {
	Host         string
	Port         int
	ProtocolType string
}

type LogConfig struct {
	FilePath string
	Console  bool
}

func LoadConfig(configPath string) (*AppConfig, error) {
	viper := viper.New()
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
