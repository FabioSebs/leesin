package main

import (
	"github.com/FabioSebs/leesin/scraper"
)

func main() {
	ws := scraper.NewTripScraper()
	collector := ws.CollectorSetup()
	ws.GetReviewsSynchronously(collector)
}
