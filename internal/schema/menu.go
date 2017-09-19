package schema

type Menu struct {
	ID      int `json:"id"`
	VenueID int `json:"venue_id"`
}

type MenuItem struct {
	ID          int     `json:"id"`
	MenuID      int     `json:"menu_id"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

type MenuDateTime struct {
	ID        int  `json:"id"`
	MenuID    int  `json:"menu_id"`
	Monday    bool `json:"monday"`
	Tuesday   bool `json:"tuesday"`
	Wednesday bool `json:"wednesday"`
	Thursday  bool `json:"thursday"`
	Friday    bool `json:"friday"`
	Saturday  bool `json:"saturday"`
	Sunday    bool `json:"sunday"`
	StartAt   bool `json:"start_at"`
	EndAt     bool `json:"end_at"`
}
