// filepath: g:\repos\ham-radio-mcp\internal\tools\bearing.go
package tools

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterAntennaBearingTool(s *server.MCPServer) {
	tool := mcp.NewTool("antenna-bearing",
		mcp.WithDescription("Calculate antenna bearing"),
		mcp.WithString("origin-latitude",
			mcp.Required(),
			mcp.Description("Origin station latitude in decimal degrees"),
		),
		mcp.WithString("origin-longitude",
			mcp.Required(),
			mcp.Description("Origin station longtitude in decimal degrees"),
		),
		mcp.WithString("destination-latitude",
			mcp.Required(),
			mcp.Description("Destination station latitude in decimal degrees"),
		),
		mcp.WithString("destination-longitude",
			mcp.Required(),
			mcp.Description("Destination station longtitude in decimal degrees"),
		),
	)

	// Add tool handler
	s.AddTool(tool, AntennaBearing)
}

func AntennaBearing(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract coordinates from request
	originLatStr, ok := request.Params.Arguments["origin-latitude"].(string)
	if !ok {
		return nil, errors.New("origin-latitude must be a string")
	}

	originLonStr, ok := request.Params.Arguments["origin-longitude"].(string)
	if !ok {
		return nil, errors.New("origin-longitude must be a string")
	}

	destLatStr, ok := request.Params.Arguments["destination-latitude"].(string)
	if !ok {
		return nil, errors.New("destination-latitude must be a string")
	}

	destLonStr, ok := request.Params.Arguments["destination-longitude"].(string)
	if !ok {
		return nil, errors.New("destination-longitude must be a string")
	}

	// Convert coordinates to float64
	originLat, err := strconv.ParseFloat(originLatStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid origin latitude: %v", err)
	}
	originLon, err := strconv.ParseFloat(originLonStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid origin longitude: %v", err)
	}
	destLat, err := strconv.ParseFloat(destLatStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid destination latitude: %v", err)
	}
	destLon, err := strconv.ParseFloat(destLonStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid destination longitude: %v", err)
	}

	// Calculate distance and bearing
	distanceKm, distanceMiles, bearing := calculateDistanceAndBearing(originLat, originLon, destLat, destLon)

	// Prepare result
	result := fmt.Sprintf("## Antenna Bearing Results\n\n"+
		"**Distance:** %.2f miles (%.2f km)\n\n"+
		"**Bearing:** %.2f degrees from North",
		distanceMiles, distanceKm, bearing)

	return mcp.NewToolResultText(result), nil
}

// toRadians converts degrees to radians
func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// toDegrees converts radians to degrees
func toDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

// calculateDistanceAndBearing calculates the distance in km and miles, and the initial bearing in degrees
// between two points on the Earth's surface using their latitude and longitude coordinates
func calculateDistanceAndBearing(lat1, lon1, lat2, lon2 float64) (distanceKm, distanceMiles, bearing float64) {
	// Earth's radius in kilometers
	earthRadiusKm := 6371.0

	// Convert degrees to radians
	lat1Rad := toRadians(lat1)
	lon1Rad := toRadians(lon1)
	lat2Rad := toRadians(lat2)
	lon2Rad := toRadians(lon2)

	// Calculate differences
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	// Haversine formula for distance
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distanceKm = earthRadiusKm * c

	// Convert distance to miles (1 km = 0.621371 miles)
	distanceMiles = distanceKm * 0.621371

	// Calculate initial bearing
	y := math.Sin(lon2Rad-lon1Rad) * math.Cos(lat2Rad)
	x := math.Cos(lat1Rad)*math.Sin(lat2Rad) -
		math.Sin(lat1Rad)*math.Cos(lat2Rad)*math.Cos(lon2Rad-lon1Rad)

	bearingRad := math.Atan2(y, x)
	bearing = toDegrees(bearingRad)

	// Convert bearing to 0-360 range
	if bearing < 0 {
		bearing += 360
	}

	return distanceKm, distanceMiles, bearing
}
