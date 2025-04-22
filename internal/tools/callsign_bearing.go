package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pleska/ham-radio-assistant/internal/models"
)

// RegisterCallsignBearingTool registers the callsign bearing tool with the MCP server
func RegisterCallsignBearingTool(s *server.MCPServer) {
	// Add tool
	tool := mcp.NewTool("callsign-bearing",
		mcp.WithDescription("Calculate bearing between two callsigns"),
		mcp.WithString("origin-callsign",
			mcp.Required(),
			mcp.Description("Your callsign"),
			mcp.Pattern("^[A-Z0-9]{1,2}[0-9][A-Z]{1,3}$"),
		),
		mcp.WithString("destination-callsign",
			mcp.Required(),
			mcp.Description("Destination callsign"),
			mcp.Pattern("^[A-Z0-9]{1,2}[0-9][A-Z]{1,3}$"),
		),
	)

	// Add tool handler
	s.AddTool(tool, CallsignBearing)
}

// CallsignBearing is a tool handler for calculating bearing between two callsigns
func CallsignBearing(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	originCallsign, ok := request.Params.Arguments["origin-callsign"].(string)
	if !ok {
		return nil, errors.New("origin-callsign must be a string")
	}

	destCallsign, ok := request.Params.Arguments["destination-callsign"].(string)
	if !ok {
		return nil, errors.New("destination-callsign must be a string")
	}

	// Look up origin callsign
	originInfo, err := lookupCallsign(originCallsign)
	if err != nil {
		return nil, fmt.Errorf("error looking up origin callsign: %v", err)
	}
	if originInfo.Status != "VALID" {
		return mcp.NewToolResultText(fmt.Sprintf("Origin callsign %s is not valid", originCallsign)), nil
	}

	// Look up destination callsign
	destInfo, err := lookupCallsign(destCallsign)
	if err != nil {
		return nil, fmt.Errorf("error looking up destination callsign: %v", err)
	}
	if destInfo.Status != "VALID" {
		return mcp.NewToolResultText(fmt.Sprintf("Destination callsign %s is not valid", destCallsign)), nil
	}

	// Convert coordinates to float64
	originLat, originLon, err := parseCoordinates(originInfo.Location.Latitude, originInfo.Location.Longitude)
	if err != nil {
		return nil, fmt.Errorf("error parsing origin coordinates: %v", err)
	}

	destLat, destLon, err := parseCoordinates(destInfo.Location.Latitude, destInfo.Location.Longitude)
	if err != nil {
		return nil, fmt.Errorf("error parsing destination coordinates: %v", err)
	}

	// Calculate distance and bearing
	distanceKm, distanceMiles, bearing := calculateDistanceAndBearing(originLat, originLon, destLat, destLon)

	// Format response
	var result string
	result = fmt.Sprintf("## Antenna Bearing: %s to %s\n\n", originCallsign, destCallsign)
	result += fmt.Sprintf("**%s Location:** %s, %s\n", originCallsign, originInfo.Location.Latitude, originInfo.Location.Longitude)
	result += fmt.Sprintf("**Grid Square:** %s\n\n", originInfo.Location.Gridsquare)
	result += fmt.Sprintf("**%s Location:** %s, %s\n", destCallsign, destInfo.Location.Latitude, destInfo.Location.Longitude)
	result += fmt.Sprintf("**Grid Square:** %s\n\n", destInfo.Location.Gridsquare)
	result += fmt.Sprintf("**Distance:** %.2f miles (%.2f km)\n\n", distanceMiles, distanceKm)
	result += fmt.Sprintf("**Bearing:** %.1f degrees from North", bearing)

	return mcp.NewToolResultText(result), nil
}

// lookupCallsign looks up a callsign using the callook.info API
func lookupCallsign(callsign string) (*models.CallsignResponse, error) {
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

	return &result, nil
}

// parseCoordinates parses latitude and longitude strings to float64
func parseCoordinates(lat, lon string) (float64, float64, error) {
	latitude, err := parseCoordinate(lat)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid latitude: %v", err)
	}

	longitude, err := parseCoordinate(lon)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid longitude: %v", err)
	}

	return latitude, longitude, nil
}

// parseCoordinate parses a coordinate string to float64
func parseCoordinate(coord string) (float64, error) {
	var result float64
	_, err := fmt.Sscanf(coord, "%f", &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}
