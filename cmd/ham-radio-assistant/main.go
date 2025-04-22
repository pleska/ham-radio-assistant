package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pleska/ham-radio-assistant/internal/api"
	"github.com/pleska/ham-radio-assistant/internal/config"
)

func main() {
	// Get executable directory to find config relative to it
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		os.Exit(1)
	}
	execDir := filepath.Dir(execPath)
	configPath := filepath.Join(execDir, "config.json")

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Create and start the MCP server
	server := api.NewServer(cfg)
	if err := server.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
