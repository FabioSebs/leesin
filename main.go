package main

import (
	"github.com/FabioSebs/leesin/scraper"
)

func main() {
	ws := scraper.NewESectorScraper()
	collector := ws.CollectorSetup()
	ws.GetReviewsSynchronously(collector)
}
