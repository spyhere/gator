package config

import (
	"encoding/json"
	"os"
)

const CONFIG_FILE_NAME = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name,omitempty"`
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	return writeNewConfig(c)
}

func ReadConfigJSON() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	fileB, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}
	var res Config
	if err := json.Unmarshal(fileB, &res); err != nil {
		return Config{}, err
	}
	return res, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/" + CONFIG_FILE_NAME, nil
}

func writeNewConfig(c *Config) error {
	sliceB, err := json.Marshal(c)
	if err != nil {
		return err
	}
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, sliceB, 0666)
}
