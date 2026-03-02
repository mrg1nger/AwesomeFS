package tray

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"awesomefs/config"
	"awesomefs/server"
)

// ShowSettingsDialog opens the config file in the default text editor
func ShowSettingsDialog(srv *server.Server, cfg *config.Config) {
	// Get config file path
	configPath, err := getConfigPath()
	if err != nil {
		fmt.Printf("Failed to get config path: %v\n", err)
		return
	}

	// Open config file in default editor
	cmd := exec.Command("notepad.exe", configPath)
	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to open config file: %v\n", err)
		// Fallback: try to open the directory
		cmd = exec.Command("explorer.exe", "/select,"+configPath)
		cmd.Start()
	}
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	// Use the same logic as in config/config.go
	appData, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}

	// For now, just return a simple path
	// In production, this should match the config package's getConfigPath
	return filepath.Join(appData, "config.json"), nil
}
