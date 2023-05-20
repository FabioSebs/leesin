package scraper

type Review struct {
	User   string `json:"user"`
	Rating string `json:"rating"`
	Review string `json:"review"`
	Date   string `json:"date"`
}
