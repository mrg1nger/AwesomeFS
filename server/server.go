package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Server represents the HTTP file server
type Server struct {
	httpServer *http.Server
	config     *Config
	running    bool
	mu         sync.RWMutex

	// Event callbacks
	onStart []func()
	onStop  []func()
	onError []func(error)
}

// Config holds the server configuration
type Config struct {
	Host     string
	Port     int
	FilesDir string
}

// New creates a new Server instance
func New(cfg *Config) *Server {
	return &Server{
		config: cfg,
	}
}

// Start begins serving files
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("server is already running")
	}

	// Create file server handler
	fs := http.FileServer(http.Dir(s.config.FilesDir))

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: fs,
	}

	// Start server in goroutine
	go func() {
		s.mu.Lock()
		s.running = true
		s.mu.Unlock()

		// Notify listeners
		for _, fn := range s.onStart {
			fn()
		}

		// ListenAndServe blocks until server stops
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			// Notify error listeners
			for _, fn := range s.onError {
				fn(err)
			}
			log.Printf("Server error: %v", err)
		}

		s.mu.Lock()
		s.running = false
		s.mu.Unlock()

		// Notify stop listeners
		for _, fn := range s.onStop {
			fn()
		}
	}()

	return nil
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running || s.httpServer == nil {
		return fmt.Errorf("server is not running")
	}

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	s.running = false
	s.httpServer = nil

	return nil
}

// IsRunning returns whether the server is currently running
func (s *Server) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetAddress returns the server address
func (s *Server) GetAddress() string {
	return fmt.Sprintf("http://%s:%d", s.config.Host, s.config.Port)
}

// OnStart registers a callback for when the server starts
func (s *Server) OnStart(fn func()) {
	s.onStart = append(s.onStart, fn)
}

// OnStop registers a callback for when the server stops
func (s *Server) OnStop(fn func()) {
	s.onStop = append(s.onStop, fn)
}

// OnError registers a callback for when the server encounters an error
func (s *Server) OnError(fn func(error)) {
	s.onError = append(s.onError, fn)
}
