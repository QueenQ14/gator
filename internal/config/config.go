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

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	config_data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err := json.Unmarshal(config_data, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func (cfg *Config) SetUser(user string) error {
	cfg.UserName = user
	if err := write(*cfg); err != nil {
		return err
	}
	return nil
}
