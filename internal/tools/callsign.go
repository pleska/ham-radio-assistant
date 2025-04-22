package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pleska/ham-radio-assistant/internal/models"
)

// RegisterCallsignLookupTool registers the callsign lookup tool with the MCP server
func RegisterCallsignLookupTool(s *server.MCPServer) {
	// Add tool
	tool := mcp.NewTool("callsign-lookup",
		mcp.WithDescription("Lookup a callsign using callook.info"),
		mcp.WithString("callsign",
			mcp.Required(),
			mcp.Description("Amateur radio callsign"),
			mcp.Pattern("^[A-Z0-9]{1,2}[0-9][A-Z]{1,3}$"),
		),
	)

	// Add tool handler
	s.AddTool(tool, CallsignLookup)
}

// CallsignLookup is a tool handler for looking up amateur radio callsigns
func CallsignLookup(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	callsign, ok := request.Params.Arguments["callsign"].(string)
	if !ok {
		return nil, errors.New("callsign must be a string")
	}

	// Make API request to callook.info
	url := fmt.Sprintf("https://callook.info/%s/json", callsign)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Decode JSON response
	var result models.CallsignResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	// Check if callsign is valid
	if result.Status != "VALID" {
		return mcp.NewToolResultText(fmt.Sprintf("Callsign %s is not valid", callsign)), nil
	}

	// Format response
	var response strings.Builder
	response.WriteString(fmt.Sprintf("## Callsign Information for %s\n\n", result.Current.Callsign))
	response.WriteString(fmt.Sprintf("**License Class:** %s\n", result.Current.OperClass))
	response.WriteString(fmt.Sprintf("**Name:** %s\n", result.Name))
	response.WriteString(fmt.Sprintf("**Type:** %s\n\n", result.Type))

	response.WriteString("### Location\n")
	response.WriteString(fmt.Sprintf("**Address:** %s, %s\n", result.Address.Line1, result.Address.Line2))
	response.WriteString(fmt.Sprintf("**Grid Square:** %s\n", result.Location.Gridsquare))
	response.WriteString(fmt.Sprintf("**Coordinates:** %s, %s\n\n", result.Location.Latitude, result.Location.Longitude))

	response.WriteString("### License Information\n")
	response.WriteString(fmt.Sprintf("**Grant Date:** %s\n", result.OtherInfo.GrantDate))
	response.WriteString(fmt.Sprintf("**Expiry Date:** %s\n", result.OtherInfo.ExpiryDate))
	response.WriteString(fmt.Sprintf("**Last Action Date:** %s\n", result.OtherInfo.LastActionDate))
	response.WriteString(fmt.Sprintf("**FRN:** %s\n", result.OtherInfo.Frn))

	if result.Previous.Callsign != "" {
		response.WriteString(fmt.Sprintf("\n**Previous Callsign:** %s (%s)\n",
			result.Previous.Callsign, result.Previous.OperClass))
	}

	response.WriteString(fmt.Sprintf("\n[View on ULS](%s)", result.OtherInfo.UlsUrl))

	return mcp.NewToolResultText(response.String()), nil
}
