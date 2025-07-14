package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vinzmyko/mdello/trello"
)

type Config struct {
	Token        string        `json:"token"`
	CurrentBoard *trello.Board `json:"currentBoard"`
	DateFormat   string        `json:"DateFormat"`
}

func SaveConfig(config Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".mdello")
	if err := os.MkdirAll(configDir, 0700); err != nil { // 0700 = owner read/write/execute only, only current user can access dir
		return fmt.Errorf("could not create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	return os.WriteFile(configFile, data, 0600) // 0600 = owner read/write only, only current user can read/write
}

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not find home directory: %w", err)
	}

	configFile := filepath.Join(homeDir, ".mdello", "config.json")
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}

	return &config, nil
}

func (cfg *Config) UpdateToken(newToken string) {
	cfg.Token = newToken
}

func (cfg *Config) UpdateCurrentBoard(newBoard *trello.Board) {
	cfg.CurrentBoard = newBoard
}

func (cfg *Config) UpdateDateFormat(newDateFormat string) {
	cfg.DateFormat = newDateFormat
}

func (cfg *Config) Save() error {
	return SaveConfig(*cfg)
}
