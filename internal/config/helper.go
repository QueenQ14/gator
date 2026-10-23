package config

import (
	"encoding/json"
	"os"
)

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := home_dir + "/" + configFileName
	return path, nil
}

func write(cfg Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, jsonData, 0666); err != nil {
		return err
	}
	return nil
}
