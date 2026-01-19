package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config represents the application configuration
type Config struct {
	Version      int        `json:"version"`
	Presets      []Preset   `json:"presets"`
	HistoryLimit int        `json:"history_limit"`
	History      []History  `json:"history"`
	Settings     Settings   `json:"settings"`
	ActiveJob    *ActiveJob `json:"active_job"`
}

// Preset represents a quick duration preset
type Preset struct {
	Label   string `json:"label"`
	Minutes int    `json:"minutes"`
}

// History represents a past shutdown event
type History struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	DurationSeconds int       `json:"duration_seconds"`
	ScheduledFor    time.Time `json:"scheduled_for"`
	Status          string    `json:"status"` // ok, cancelled, failed, dry-run
	OS              string    `json:"os"`
	Command         string    `json:"command"`
}

// Settings represents application settings
type Settings struct {
	Confirm       bool `json:"confirm"`
	DryRunDefault bool `json:"dry_run_default"`
}

// ActiveJob represents currently running shutdown job
type ActiveJob struct {
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	DurationSec int       `json:"duration_sec"`
	Command     string    `json:"command"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Version: 1,
		Presets: []Preset{
			{Label: "15m", Minutes: 15},
			{Label: "30m", Minutes: 30},
			{Label: "45m", Minutes: 45},
			{Label: "60m", Minutes: 60},
			{Label: "90m", Minutes: 90},
			{Label: "120m", Minutes: 120},
		},
		HistoryLimit: 20,
		History:      []History{},
		Settings: Settings{
			Confirm:       true,
			DryRunDefault: false,
		},
		ActiveJob: nil,
	}
}

// ConfigPath returns the path to the config file
func ConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config dir: %w", err)
	}

	gtsDir := filepath.Join(configDir, "gts")
	if err := os.MkdirAll(gtsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config dir: %w", err)
	}

	return filepath.Join(gtsDir, "state.json"), nil
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	// If file doesn't exist, return default config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// AddHistory adds a new history entry and maintains the limit
func (c *Config) AddHistory(h History) {
	c.History = append([]History{h}, c.History...)
	if len(c.History) > c.HistoryLimit {
		c.History = c.History[:c.HistoryLimit]
	}
}

// UpdateHistoryStatus updates the status of the most recent history entry
func (c *Config) UpdateHistoryStatus(status string) {
	if len(c.History) > 0 {
		c.History[0].Status = status
	}
}

// DeleteHistory removes a history entry by ID
func (c *Config) DeleteHistory(id string) {
	for i, h := range c.History {
		if h.ID == id {
			c.History = append(c.History[:i], c.History[i+1:]...)
			return
		}
	}
}
