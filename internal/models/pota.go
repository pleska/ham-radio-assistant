package models

import (
	"strconv"
	"strings"
)

// ParkReference represents a POTA (Parks on the Air) park reference.
type ParkReference struct {
	Reference    string  `json:"reference" csv:"reference"`
	Name         string  `json:"name" csv:"name"`
	Active       bool    `json:"active" csv:"active"`
	EntityID     int     `json:"entityId" csv:"entityId"`
	LocationDesc string  `json:"locationDesc" csv:"locationDesc"` // e.g., "US-ME"
	Latitude     float64 `json:"latitude" csv:"latitude"`
	Longitude    float64 `json:"longitude" csv:"longitude"`
	Grid         string  `json:"grid" csv:"grid"` // Maidenhead grid square
}

// UnmarshalCSV is a custom unmarshaler for the ParkReference struct
// It handles specific field parsing
func (p *ParkReference) UnmarshalCSV(field []string, header []string) error {
	for i, h := range header {
		if i >= len(field) {
			continue
		}

		switch h {
		case "reference":
			p.Reference = field[i]
		case "name":
			p.Name = field[i]
		case "active":
			active, err := ParseBool(field[i])
			if err == nil {
				p.Active = active
			} else {
				p.Active = true // Default to active if parsing fails
			}
		case "entityId":
			if id, err := strconv.Atoi(field[i]); err == nil {
				p.EntityID = id
			}
		case "locationDesc":
			p.LocationDesc = field[i]
		case "latitude":
			if lat, err := strconv.ParseFloat(field[i], 64); err == nil {
				p.Latitude = lat
			}
		case "longitude":
			if lon, err := strconv.ParseFloat(field[i], 64); err == nil {
				p.Longitude = lon
			}
		case "grid":
			p.Grid = field[i]
		}
	}
	return nil
}

// ParseBool parses a string representation of a boolean from CSV
// It handles different formats like "1", "true", "yes", etc.
func ParseBool(value string) (bool, error) {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return false, nil
	}

	// Check for common true representations
	if value == "1" || value == "true" || value == "yes" || value == "y" {
		return true, nil
	}

	// Parse as integer (0 = false, non-zero = true)
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return i != 0, nil
	}

	// Parse as normal bool
	return strconv.ParseBool(value)
}
