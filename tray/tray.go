package tray

import (
	"fmt"
	"os/exec"

	"awesomefs/assets"
	"awesomefs/config"
	"awesomefs/server"

	"fyne.io/systray"
)

// TrayApp represents the system tray application
type TrayApp struct {
	server *server.Server
	config *config.Config

	// Menu items
	mStatus      *systray.MenuItem
	mStartStop   *systray.MenuItem
	mOpenBrowser *systray.MenuItem
	mSettings    *systray.MenuItem
	mAutoStart   *systray.MenuItem
	mQuit        *systray.MenuItem

	// State
	autoStartEnabled bool
}

// New creates a new TrayApp instance
func New(srv *server.Server, cfg *config.Config) *TrayApp {
	return &TrayApp{
		server: srv,
		config: cfg,
	}
}

// Run starts the system tray application
func (t *TrayApp) Run() {
	systray.Run(t.onReady, t.onExit)
}

// onReady is called when the system tray is ready
func (t *TrayApp) onReady() {
	// Set icon and title
	systray.SetIcon(assets.IconData)
	systray.SetTitle("AwesomeFS")
	systray.SetTooltip("AwesomeFS File Server")

	// Create menu items
	t.mStatus = systray.AddMenuItem("Server: Stopped", "Current server status")
	t.mStatus.Disable()

	systray.AddSeparator()

	t.mStartStop = systray.AddMenuItem("Start Server", "Start or stop the file server")
	t.mOpenBrowser = systray.AddMenuItem("Open in Browser", "Open the file server in your browser")
	t.mSettings = systray.AddMenuItem("Settings...", "Configure server settings")

	systray.AddSeparator()

	// Check if auto-start is enabled
	t.autoStartEnabled = IsAutoStartEnabled()
	t.mAutoStart = systray.AddMenuItemCheckbox("Start with Windows", "Launch AwesomeFS on Windows startup", t.autoStartEnabled)

	systray.AddSeparator()

	t.mQuit = systray.AddMenuItem("Quit", "Exit AwesomeFS")

	// Set up event handlers
	go t.handleEvents()

	// Set up server event listeners
	t.server.OnStart(func() {
		t.updateStatus(true)
	})
	t.server.OnStop(func() {
		t.updateStatus(false)
	})
	t.server.OnError(func(err error) {
		// TODO: Show error notification
		fmt.Printf("Server error: %v\n", err)
	})

	// Auto-start server on launch
	if err := t.server.Start(); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

// onExit is called when the system tray is exiting
func (t *TrayApp) onExit() {
	// Stop the server if running
	if t.server.IsRunning() {
		t.server.Stop()
	}
}

// handleEvents handles menu item click events
func (t *TrayApp) handleEvents() {
	for {
		select {
		case <-t.mStartStop.ClickedCh:
			t.toggleServer()

		case <-t.mOpenBrowser.ClickedCh:
			t.openBrowser()

		case <-t.mSettings.ClickedCh:
			t.showSettings()

		case <-t.mAutoStart.ClickedCh:
			t.toggleAutoStart()

		case <-t.mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

// toggleServer starts or stops the server
func (t *TrayApp) toggleServer() {
	if t.server.IsRunning() {
		if err := t.server.Stop(); err != nil {
			fmt.Printf("Failed to stop server: %v\n", err)
		}
	} else {
		if err := t.server.Start(); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
		}
	}
}

// updateStatus updates the status menu item and icon
func (t *TrayApp) updateStatus(running bool) {
	if running {
		t.mStatus.SetTitle("Server: Running ✓")
		t.mStartStop.SetTitle("Stop Server")
		systray.SetTooltip("AwesomeFS File Server - Running")
	} else {
		t.mStatus.SetTitle("Server: Stopped")
		t.mStartStop.SetTitle("Start Server")
		systray.SetTooltip("AwesomeFS File Server - Stopped")
	}
}

// openBrowser opens the server URL in the default browser
func (t *TrayApp) openBrowser() {
	url := t.server.GetAddress()
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}
}

// showSettings opens the settings dialog
func (t *TrayApp) showSettings() {
	ShowSettingsDialog(t.server, t.config)
}

// toggleAutoStart enables or disables auto-start on Windows boot
func (t *TrayApp) toggleAutoStart() {
	if t.autoStartEnabled {
		if err := DisableAutoStart(); err != nil {
			fmt.Printf("Failed to disable auto-start: %v\n", err)
			return
		}
		t.autoStartEnabled = false
		t.mAutoStart.Uncheck()
	} else {
		if err := EnableAutoStart(); err != nil {
			fmt.Printf("Failed to enable auto-start: %v\n", err)
			return
		}
		t.autoStartEnabled = true
		t.mAutoStart.Check()
	}
}
