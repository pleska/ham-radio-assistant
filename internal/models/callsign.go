package models

// CallsignResponse represents the response from callook.info API
type CallsignResponse struct {
	Status  string `json:"status"`
	Type    string `json:"type"`
	Current struct {
		Callsign  string `json:"callsign"`
		OperClass string `json:"operClass"`
	} `json:"current"`
	Previous struct {
		Callsign  string `json:"callsign"`
		OperClass string `json:"operClass"`
	} `json:"previous"`
	Trustee struct {
		Callsign string `json:"callsign"`
		Name     string `json:"name"`
	} `json:"trustee"`
	Name    string `json:"name"`
	Address struct {
		Line1 string `json:"line1"`
		Line2 string `json:"line2"`
		Attn  string `json:"attn"`
	} `json:"address"`
	Location struct {
		Latitude   string `json:"latitude"`
		Longitude  string `json:"longitude"`
		Gridsquare string `json:"gridsquare"`
	} `json:"location"`
	OtherInfo struct {
		GrantDate      string `json:"grantDate"`
		ExpiryDate     string `json:"expiryDate"`
		LastActionDate string `json:"lastActionDate"`
		Frn            string `json:"frn"`
		UlsUrl         string `json:"ulsUrl"`
	} `json:"otherInfo"`
}
