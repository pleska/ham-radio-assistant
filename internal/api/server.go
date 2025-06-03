package api

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"github.com/pleska/ham-radio-assistant/internal/config"
	"github.com/pleska/ham-radio-assistant/internal/tools"
)

// Server represents the MCP API server
type Server struct {
	config    *config.Config
	mcpServer *server.MCPServer
}

// NewServer creates a new MCP server instance
func NewServer(cfg *config.Config) *Server {
	mcpServer := server.NewMCPServer(
		"Ham Radio Assistant",
		"1.0.0",
	)

	return &Server{
		config:    cfg,
		mcpServer: mcpServer,
	}
}

// RegisterTools registers all the tools with the MCP server
func (s *Server) RegisterTools() {
	// Register the callsign lookup tool
	tools.RegisterCallsignLookupTool(s.mcpServer)
	tools.RegisterAntennaBearingTool(s.mcpServer)
	tools.RegisterCallsignBearingTool(s.mcpServer)
	tools.RegisterPotaParkLookupTool(s.mcpServer)
	tools.RegisterPotaSpotsTool(s.mcpServer)

	// Additional tools can be registered here in the future
}

// Start starts the MCP server
func (s *Server) Start() error {
	s.RegisterTools()

	// Start the stdio server
	if err := server.ServeStdio(s.mcpServer); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
