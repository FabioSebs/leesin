package scraper

type EV struct {
	Name  string
	Price string
}

// body div.bodyconstraint--full-width div.bodyconstraint-inner div.af5895d4b2 div.df7e6ba27d div.bcbf33c5c3 div.d4924c9e74 div.c82435a4b8
type Booking struct {
	Title        string `json:"Title"`
	Location     string `json:"Location"`
	Address      string `json:"Address"`
	PostCode     string `json:"PostCode"`
	Review       string `json:"Review"`
	ReviewNumber string `json:"ReviewNumber"`
	Price        string `json:"Price"`
	PropertyType string `json:"PropertyType"`
	Source       string `json:"Source"`
}

type Publication struct {
	Year        string `json:"Publication year"`
	Sector      string `json:"Sector"`
	VehicleType string `json:"Vehicle type"`
	Region      string `json:"Region"`
	Metric      string `json:"Metric"`
	Format      string `json:"Format"`
	Title       string `json:"Title"`
	Author1     string `json:"Author 1"`
	Author2     string `json:"Author 2"`
	Author3     string `json:"Author 3"`
	Author4     string `json:"Author 4"`
	Author5     string `json:"Author 5"`
	Author6     string `json:"Author 6"`
	Link        string `json:"Link"`
}
