package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl    string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

func Read() (Config) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}
	}
	config_data, err := os.ReadFile(path)
	if err != nil {
		return Config{}
	}
	var config Config
	if err := json.Unmarshal(config_data, &config); err != nil {
		return Config{}
	}
	return config
}

func (cfg *Config) SetUser(user string) error {
	cfg.UserName = user
	if err := write(*cfg); err != nil {
		return err
	}
	return nil
}
