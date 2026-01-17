package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configFileName = "gatorconfig.json"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func write(cfg *Config) error {
	updatedData, err := json.Marshal(cfg)
	if err != nil {
		return nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil
	}
	if err := os.WriteFile(filepath.Join(wd, configFileName), updatedData, 0o644); err != nil {
		return nil
	}
	return nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	if err := write(c); err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}
	content, err := os.ReadFile(filepath.Join(wd, configFileName))
	if err != nil {
		fmt.Printf("failed")
		return Config{}, err
	}
	config := Config{}
	if err := json.Unmarshal(content, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
