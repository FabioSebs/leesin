package main

import "github.com/FabioSebs/leesin/scraper"

func main() {
	ws := scraper.NewWebScraper()
	collector := ws.CollectorSetup()
	ws.GetReviewsConcurrently(collector)
}
