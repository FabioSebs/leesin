package tests

import (
	"testing"

	"github.com/FabioSebs/leesin/scraper"
)

func TestScrape(t *testing.T) {
	syncColly := scraper.NewWebScraper()
	conColly := scraper.NewWebScraper()

	synccollector := syncColly.CollectorSetup()
	concollector := conColly.CollectorSetup()

	t.Run("Synchronous Scrape", func(t *testing.T) {
		t.Parallel()
		_, duration := syncColly.GetReviewsSynchronously(synccollector)
		if duration == 0 {
			t.Errorf("duration is %s, process failed", duration)
		}
		t.Logf("duration: %s", duration)
	})

	t.Run("Concurrent Scrape", func(t *testing.T) {
		t.Parallel()
		_, duration := conColly.GetReviewsConcurrently(concollector)
		if duration == 0 {
			t.Errorf("duration is %s, process failed", duration)
		}
		t.Logf("duration: %s", duration)
	})
}

func BenchmarkSyncScraper(b *testing.B) {
	colly := scraper.NewWebScraper()
	collector := colly.CollectorSetup()
	for i := 0; i < b.N; i++ {
		colly.GetReviewsSynchronously(collector)
	}
}

func BenchmarkConcurrentScraper(b *testing.B) {
	colly := scraper.NewWebScraper()
	collector := colly.CollectorSetup()
	for i := 0; i < b.N; i++ {
		colly.GetReviewsConcurrently(collector)
	}
}
