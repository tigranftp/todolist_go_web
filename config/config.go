package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port string
	DSN  string
}

func GetConfig() (*Config, error) {
	configContent, err := os.ReadFile("res/config.json")
	if err != nil {
		return nil, err
	}
	cfg := new(Config)
	err = json.Unmarshal(configContent, cfg)
	fmt.Println(*cfg)
	return cfg, err
}
