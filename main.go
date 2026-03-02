package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"awesomefs/config"
	"awesomefs/server"
	"awesomefs/tray"
)

var (
	trayMode = flag.Bool("tray", false, "Run in system tray mode (Windows only)")
	host     = flag.String("host", "", "Interface to bind to (overrides config)")
	port     = flag.Int("port", 0, "Port to host the server on (overrides config, 0 = use config)")
)

func main() {
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: Failed to load config: %v, using defaults\n", err)
		cfg = config.New()
	}

	// Override config with command-line flags if provided
	if *host != "" {
		cfg.Host = *host
	}
	if *port != 0 {
		cfg.Port = *port
	}

	// Create server
	srv := server.New(&server.Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		FilesDir: cfg.FilesDir,
	})

	// Run in appropriate mode
	if *trayMode {
		// Hide console window on Windows
		hideConsole()
		runTrayMode(srv, cfg)
	} else {
		runCLIMode(srv)
	}
}

// runCLIMode runs the server in traditional command-line mode
func runCLIMode(srv *server.Server) {
	fmt.Printf("Starting file server at %s\n", srv.GetAddress())
	fmt.Printf("Serving files from ./files directory\n")
	fmt.Println("Press Ctrl+C to stop")

	// Start the server
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-sigChan
	fmt.Println("\nShutting down server...")

	// Stop the server
	if err := srv.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}

	fmt.Println("Server stopped")
}

// runTrayMode runs the server in system tray mode
func runTrayMode(srv *server.Server, cfg *config.Config) {
	// Create and run system tray application
	app := tray.New(srv, cfg)
	app.Run()
}

// hideConsole hides the console window on Windows
func hideConsole() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	user32 := syscall.NewLazyDLL("user32.dll")

	getConsoleWindow := kernel32.NewProc("GetConsoleWindow")
	showWindow := user32.NewProc("ShowWindow")

	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd != 0 {
		// SW_HIDE = 0
		showWindow.Call(hwnd, 0)
	}
}
