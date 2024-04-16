package scraper

type EV struct {
	Name  string
	Price string
}

type EVModel struct {
	Name         string `json:"name"`
	Acceleration string `json:"accelaration"`
	TopSpeed     string `json:"top-speed"`
	Range        string `json:"range"`
	Efficiency   string `json:"efficiency"`
	FastCharge   string `json:"fast-charge"`
	// SafetyRating string `json:"safety-rating"`
	Price string `json:"price"`
}
