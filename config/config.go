package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	FilesDir string `json:"filesDir"`
}

// Default configuration values
var defaults = Config{
	Host:     "0.0.0.0",
	Port:     8081,
	FilesDir: "./files",
}

// New creates a new Config with default values
func New() *Config {
	cfg := defaults
	return &cfg
}

// Load reads configuration from the config file
// If the file doesn't exist, it returns defaults
func Load() (*Config, error) {
	cfg := New()

	configPath, err := getConfigPath()
	if err != nil {
		return cfg, fmt.Errorf("failed to get config path: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// File doesn't exist, return defaults
		return cfg, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// Save writes the configuration to the config file
func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigPath returns the path to the config file
// Uses %APPDATA%\AwesomeFS\config.json on Windows
func getConfigPath() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return "", fmt.Errorf("APPDATA environment variable not set")
	}

	configPath := filepath.Join(appData, "AwesomeFS", "config.json")
	return configPath, nil
}

// Update modifies the configuration and saves it
func (c *Config) Update(host string, port int, filesDir string) error {
	if host != "" {
		c.Host = host
	}
	if port > 0 {
		c.Port = port
	}
	if filesDir != "" {
		c.FilesDir = filesDir
	}

	return c.Save()
}
