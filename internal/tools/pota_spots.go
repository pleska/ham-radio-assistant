// Package tools provides tool implementations for the ham radio assistant
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pleska/ham-radio-assistant/internal/models"
)

const (
	potaSpotsAPIURL = "https://api.pota.app/spot/activator"
)

// RegisterPotaSpotsTool registers the POTA activator spots lookup tool with the MCP server
func RegisterPotaSpotsTool(s *server.MCPServer) {
	// Add tool
	tool := mcp.NewTool("pota-spots",
		mcp.WithDescription("Display active POTA activations by callsign or mode"),
		mcp.WithString("callsign",
			mcp.Description("Activator callsign to filter by"),
		),
		mcp.WithString("mode",
			mcp.Description("Mode to filter by (e.g., SSB, CW, FT8)"),
		),
	)

	// Add tool handler
	s.AddTool(tool, PotaSpotsLookup)
}

// PotaSpotsLookup is a tool handler for looking up current POTA activations
func PotaSpotsLookup(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get optional parameters
	callsign, _ := request.Params.Arguments["callsign"].(string)
	mode, _ := request.Params.Arguments["mode"].(string)

	// Fetch spots from the API
	spots, err := fetchPotaSpots(callsign, mode)
	if err != nil {
		return nil, fmt.Errorf("error fetching POTA spots: %v", err)
	}

	// Check if any spots were found
	if len(spots) == 0 {
		message := "No active POTA spots found"
		if callsign != "" {
			message += fmt.Sprintf(" for activator %s", callsign)
		}
		if mode != "" {
			if callsign != "" {
				message += " and"
			} else {
				message += " for"
			}
			message += fmt.Sprintf(" mode %s", mode)
		}
		return mcp.NewToolResultText(message), nil
	}

	// Format response
	var response strings.Builder
	response.WriteString("# Current POTA Activations\n\n")

	if callsign != "" {
		response.WriteString(fmt.Sprintf("Filtered by activator: **%s**\n\n", callsign))
	}
	if mode != "" {
		response.WriteString(fmt.Sprintf("Filtered by mode: **%s**\n\n", mode))
	}

	response.WriteString("| Activator | Reference | Park Name | Frequency | Mode | Location | Spotted At | Spotted By | Comments |\n")
	response.WriteString("|-----------|-----------|-----------|-----------|------|----------|------------|------------|----------|\n")

	for _, spot := range spots {
		// Parse and format the spot time
		spotTime, err := time.Parse("2006-01-02T15:04:05", spot.SpotTime)
		timeStr := spot.SpotTime
		if err == nil {
			timeStr = spotTime.Format("15:04 UTC")
		}

		// Format row
		response.WriteString(fmt.Sprintf("| %s | [%s](https://pota.app/#/park/%s) | %s | %s | %s | %s | %s | %s | %s |\n",
			spot.Activator,
			spot.Reference,
			spot.Reference,
			spot.Name,
			spot.Frequency,
			spot.Mode,
			spot.LocationDesc,
			timeStr,
			spot.Spotter,
			spot.Comments,
		))
	}

	response.WriteString("\n\nData provided by [Parks on the Air API](https://pota.app)")

	return mcp.NewToolResultText(response.String()), nil
}

// fetchPotaSpots fetches current POTA activations from the API
// and filters them based on callsign and mode if provided
func fetchPotaSpots(callsign, mode string) ([]models.POTASpot, error) {
	// Make the HTTP request to get all spots
	resp, err := http.Get(potaSpotsAPIURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to POTA API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response status: %s", resp.Status)
	}

	// Read and parse the JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var spots []models.POTASpot
	if err := json.Unmarshal(body, &spots); err != nil {
		return nil, fmt.Errorf("error parsing JSON data: %v", err)
	}

	// If no filters are provided, return all spots
	if callsign == "" && mode == "" {
		return spots, nil
	}

	// Filter spots based on callsign and/or mode
	var filteredSpots []models.POTASpot
	for _, spot := range spots {
		// Check if spot matches the callsign filter
		callsignMatch := callsign == "" || strings.EqualFold(spot.Activator, callsign)

		// Check if spot matches the mode filter
		modeMatch := mode == "" || strings.EqualFold(spot.Mode, mode)

		// Add to filtered results if both conditions are met
		if callsignMatch && modeMatch {
			filteredSpots = append(filteredSpots, spot)
		}
	}

	return filteredSpots, nil
}
