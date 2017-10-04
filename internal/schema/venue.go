package schema

type Venue struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	State    string `json:"state"`
	Zip      string `json:"zip"`
	Country  string `json:"country"`
	Image    string `json:"image"`
}

type VenueList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type VenueLists struct {
	ID          int `json:"id"`
	VenueID     int `json:"venue_id"`
	VenueListID int `json:"venue_list_id"`
}

type VenueListAdd struct {
	VenueID       int    `json:"venue_id"`
	VenueName     string `json:"venue_name"`
	VenueListID   int    `json:"venue_list_id"`
	VenueListName string `json:"venue_list_name"`
}
