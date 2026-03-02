package tray

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const appName = "AwesomeFS"

// EnableAutoStart adds AwesomeFS to Windows startup
func EnableAutoStart() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// Get absolute path
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return err
	}

	// Open registry key
	key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	// Set the value
	return key.SetStringValue(appName, exePath)
}

// DisableAutoStart removes AwesomeFS from Windows startup
func DisableAutoStart() error {
	// Open registry key
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	// Delete the value
	return key.DeleteValue(appName)
}

// IsAutoStartEnabled checks if AwesomeFS is set to start with Windows
func IsAutoStartEnabled() bool {
	// Open registry key
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Run`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return false
	}
	defer key.Close()

	// Check if value exists
	_, _, err = key.GetStringValue(appName)
	return err == nil
}
