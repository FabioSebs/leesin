package main

import (
	"github.com/FabioSebs/leesin/scraper"
	"github.com/FabioSebs/leesin/wordcounter"
)

func main() {
	wc := wordcounter.NewWordCounter()
	ws := scraper.NewWebScraper()
	collector := ws.CollectorSetup()
	reviews, _ := ws.GetReviewsSynchronously(collector)
	wc.CountReviews(reviews)
}
