package scraper

import (
	"testing"
)

func TestScrape(t *testing.T) {
	scraper := NewWebScraper()
	collector := scraper.CollectorSetup()
	t.Run("Synchronous Scrape", func(t *testing.T) {
		t.Parallel()
		duration := scraper.GetReviewsSynchronously(collector)
		if duration == 0 {
			t.Errorf("duration is %s, process failed", duration)
		}
		t.Logf("duration: %s", duration)
	})

	t.Run("Asynchronous Scrape", func(t *testing.T) {
		t.Parallel()
		duration := scraper.GetReviewsConcurrently(collector)
		if duration == 0 {
			t.Errorf("duration is %s, process failed", duration)
		}
		t.Logf("duration: %s", duration)
	})
}
