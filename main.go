package main

import "github.com/FabioSebs/leesin/scraper"

func main() {
	ws := scraper.NewWebScraper()
	ws.GetReviews()
}
