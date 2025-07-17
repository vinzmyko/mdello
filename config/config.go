package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vinzmyko/mdello/trello"
)

const (
	DateFormatISO = "2006-01-02 15:04"
	DateFormatUS  = "01-02-2006 15:04"
	DateFormatEU  = "02-01-2006 15:04"
)

type DateFormatOption struct {
	Display string
	Value   string
}

type Config struct {
	Token          string `json:"token"`
	CurrentBoardID string `json:"currentBoardId"`
	DateFormat string `json:"DateFormat"`
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

func (cfg *Config) GetCurrentBoard(trelloClient *trello.TrelloClient) (*trello.Board, error) {
	if cfg.CurrentBoardID == "" {
		return nil, fmt.Errorf("no current board set")
	}
	return trelloClient.GetBoard(cfg.CurrentBoardID)
}

func (cfg *Config) UpdateToken(newToken string) {
	cfg.Token = newToken
}

func (cfg *Config) UpdateBoardID(newBoardID string) {
	cfg.CurrentBoardID = newBoardID
}

func (cfg *Config) UpdateDateFormat(newDateFormat string) {
	cfg.DateFormat = newDateFormat
}

func (cfg *Config) Save() error {
	return SaveConfig(*cfg)
}

func GetDateFormatOptions() []DateFormatOption {
	return []DateFormatOption{
		{"International (YYYY-MM-DD)", DateFormatISO},
		{"US (MM-DD-YYYY)", DateFormatUS},
		{"European (DD-MM-YYYY)", DateFormatEU},
	}
}

func GetDisplayOptions() []string {
	options := GetDateFormatOptions()
	displays := make([]string, len(options))
	for i, option := range options {
		displays[i] = option.Display
	}
	return displays
}

func GetFormatFromDisplay(display string) (string, bool) {
	for _, option := range GetDateFormatOptions() {
		if option.Display == display {
			return option.Value, true
		}
	}
	return "", false
}
