// filepath: /workspaces/ham-radio-assistant/internal/tools/pota.go
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

const (
	potaAPIBaseURL = "https://api.pota.app/park/"
)

// RegisterPotaParkLookupTool registers the POTA park lookup tool with the MCP server
func RegisterPotaParkLookupTool(s *server.MCPServer) {
	// Add tool
	tool := mcp.NewTool("pota-park-lookup",
		mcp.WithDescription("Lookup Parks on the Air (POTA) park details by reference"),
		mcp.WithString("reference",
			mcp.Required(),
			mcp.Description("POTA park reference (e.g., US-2312)"),
			mcp.Pattern("^[A-Z0-9]{1,4}-[0-9]{1,5}$"),
		),
	)

	// Add tool handler
	s.AddTool(tool, PotaParkLookup)
}

// PotaParkLookup is a tool handler for looking up POTA park details directly from the API
func PotaParkLookup(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	reference, ok := request.Params.Arguments["reference"].(string)
	if !ok {
		return nil, errors.New("reference must be a string")
	}

	// Fetch park details using the REST API
	park, err := fetchParkDetails(reference)
	if err != nil {
		return nil, fmt.Errorf("error fetching park details: %v", err)
	}

	// Format response
	var response strings.Builder
	response.WriteString(fmt.Sprintf("## POTA Park: %s\n\n", reference))
	response.WriteString(fmt.Sprintf("**Name:** %s\n", park.Name))
	response.WriteString(fmt.Sprintf("**Location:** %s, %s\n", park.LocationDesc, park.LocationName))
	response.WriteString(fmt.Sprintf("**Status:** %s\n", formatStatus(park.IsActive())))
	response.WriteString(fmt.Sprintf("**Park Type:** %s\n\n", park.ParktypeDesc))

	if park.ParkComments != "" {
		response.WriteString(fmt.Sprintf("**Comments:** %s\n\n", park.ParkComments))
	}

	response.WriteString("### Geographic Information\n")
	response.WriteString(fmt.Sprintf("**Coordinates:** %f, %f\n", park.Latitude, park.Longitude))
	response.WriteString(fmt.Sprintf("**Grid Square:** %s (%s)\n\n", park.Grid4, park.Grid6))

	if park.AccessMethods != "" {
		response.WriteString(fmt.Sprintf("**Access Methods:** %s\n", park.AccessMethods))
	}

	if park.ActivationMethods != "" {
		response.WriteString(fmt.Sprintf("**Activation Methods:** %s\n\n", park.ActivationMethods))
	}

	if park.Website != "" {
		response.WriteString(fmt.Sprintf("**Website:** [%s](%s)\n\n", park.Website, park.Website))
	}

	if park.FirstActivator != "" {
		response.WriteString(fmt.Sprintf("**First Activated By:** %s on %s\n\n", park.FirstActivator, park.FirstActivationDate))
	}

	response.WriteString(fmt.Sprintf("[View on POTA website](https://pota.app/#/park/%s)", reference))

	return mcp.NewToolResultText(response.String()), nil
}

// fetchParkDetails fetches park details from the POTA API
func fetchParkDetails(reference string) (*models.ParkReference, error) {
	// Fetch from API
	apiURL := fmt.Sprintf("%s%s", potaAPIBaseURL, reference)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to POTA API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("park with reference %s not found", reference)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response status: %s", resp.Status)
	}

	// Read and parse the JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var parkData models.ParkReference
	if err := json.Unmarshal(body, &parkData); err != nil {
		return nil, fmt.Errorf("error parsing JSON data: %v", err)
	}

	return &parkData, nil
}

// formatStatus returns a human-readable status string
func formatStatus(active bool) string {
	if active {
		return "Active"
	}
	return "Inactive"
}
