package models

// ParkReference represents a POTA (Parks on the Air) park reference.
type ParkReference struct {
	ParkID              int     `json:"parkId"`
	Reference           string  `json:"reference"`
	Name                string  `json:"name"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	Grid4               string  `json:"grid4"`
	Grid6               string  `json:"grid6"`
	ParktypeID          int     `json:"parktypeId"`
	Active              int     `json:"active"`
	ParkComments        string  `json:"parkComments"`
	Accessibility       *string `json:"accessibility"`
	Sensitivity         *string `json:"sensitivity"`
	AccessMethods       string  `json:"accessMethods"`
	ActivationMethods   string  `json:"activationMethods"`
	Agencies            *string `json:"agencies"`
	AgencyURLs          *string `json:"agencyURLs"`
	ParkURLs            *string `json:"parkURLs"`
	Website             string  `json:"website"`
	CreatedByAdmin      string  `json:"createdByAdmin"`
	ParktypeDesc        string  `json:"parktypeDesc"`
	LocationDesc        string  `json:"locationDesc"`
	LocationName        string  `json:"locationName"`
	EntityID            int     `json:"entityId"`
	EntityName          string  `json:"entityName"`
	ReferencePrefix     string  `json:"referencePrefix"`
	EntityDeleted       int     `json:"entityDeleted"`
	FirstActivator      string  `json:"firstActivator"`
	FirstActivationDate string  `json:"firstActivationDate"`
}

// IsActive returns whether the park is active or not
func (p *ParkReference) IsActive() bool {
	return p.Active == 1
}

// POTASpot represents a spot of a POTA activation
type POTASpot struct {
	SpotID       int     `json:"spotId"`
	Activator    string  `json:"activator"`
	Frequency    string  `json:"frequency"`
	Mode         string  `json:"mode"`
	Reference    string  `json:"reference"`
	SpotTime     string  `json:"spotTime"`
	Spotter      string  `json:"spotter"`
	Comments     string  `json:"comments"`
	Source       string  `json:"source"`
	Name         string  `json:"name"`
	LocationDesc string  `json:"locationDesc"`
	Grid4        string  `json:"grid4"`
	Grid6        string  `json:"grid6"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Count        int     `json:"count"`
	Expire       int     `json:"expire"`
}
