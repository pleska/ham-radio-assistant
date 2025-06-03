package tools

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/pleska/ham-radio-assistant/internal/models"
)

const (
	potaDataURL      = "https://pota.app/all_parks_ext.csv"
	cacheMaxAgeHours = 24
)

var (
	parksCache      = make(map[string]models.ParkReference)
	parksCacheMutex sync.RWMutex
	lastCacheUpdate time.Time
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

// PotaParkLookup is a tool handler for looking up POTA park details
func PotaParkLookup(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	reference, ok := request.Params.Arguments["reference"].(string)
	if !ok {
		return nil, errors.New("reference must be a string")
	}

	// Ensure parks data is loaded and up to date
	if err := loadParksData(); err != nil {
		return nil, fmt.Errorf("error loading POTA data: %v", err)
	}

	// Look up park in cache
	parksCacheMutex.RLock()
	park, exists := parksCache[reference]
	parksCacheMutex.RUnlock()

	if !exists {
		return mcp.NewToolResultText(fmt.Sprintf("Park with reference %s not found", reference)), nil
	}

	// Format response
	var response strings.Builder
	response.WriteString(fmt.Sprintf("## POTA Park: %s\n\n", reference))
	response.WriteString(fmt.Sprintf("**Name:** %s\n", park.Name))
	response.WriteString(fmt.Sprintf("**Location:** %s\n", park.LocationDesc))
	response.WriteString(fmt.Sprintf("**Status:** %s\n\n", formatStatus(park.Active)))

	response.WriteString("### Geographic Information\n")
	response.WriteString(fmt.Sprintf("**Coordinates:** %f, %f\n", park.Latitude, park.Longitude))
	response.WriteString(fmt.Sprintf("**Grid Square:** %s\n\n", park.Grid))

	response.WriteString(fmt.Sprintf("[View on POTA website](https://pota.app/#/park/%s)", reference))

	return mcp.NewToolResultText(response.String()), nil
}

// loadParksData ensures the POTA parks data is loaded and up to date
func loadParksData() error {
	parksCacheMutex.RLock()
	cacheEmpty := len(parksCache) == 0
	cacheExpired := time.Since(lastCacheUpdate) > cacheMaxAgeHours*time.Hour
	parksCacheMutex.RUnlock()

	if cacheEmpty || cacheExpired {
		return refreshParksData()
	}
	return nil
}

// refreshParksData downloads the latest POTA parks data directly into memory
func refreshParksData() error {
	// Download fresh data
	resp, err := http.Get(potaDataURL)
	if err != nil {
		return fmt.Errorf("error downloading POTA data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response status: %s", resp.Status)
	}

	// Load directly from the response body
	return loadParksFromReader(resp.Body)
}

// loadParksFromReader reads and parses POTA parks data from an io.Reader
func loadParksFromReader(r io.Reader) error {
	var parks []*models.ParkReference

	// Use gocsv to unmarshal the CSV data directly into structs
	if err := gocsv.Unmarshal(r, &parks); err != nil {
		return fmt.Errorf("error parsing CSV data: %v", err)
	}

	// Create a new map to store the parks
	newParksCache := make(map[string]models.ParkReference)

	// Add valid parks to the cache
	for _, park := range parks {
		if park == nil || park.Reference == "" || park.Name == "" {
			// Skip invalid records
			continue
		}
		newParksCache[park.Reference] = *park
	}

	// Update cache with acquired data
	parksCacheMutex.Lock()
	parksCache = newParksCache
	lastCacheUpdate = time.Now()
	parksCacheMutex.Unlock()

	return nil
}

// customCSVUnmarshaler configures gocsv to use our custom unmarshaling logic
func init() {
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		reader := csv.NewReader(in)
		reader.FieldsPerRecord = -1 // Allow variable number of fields
		return reader
	})
}

// formatStatus returns a human-readable status string
func formatStatus(active bool) string {
	if active {
		return "Active"
	}
	return "Inactive"
}
