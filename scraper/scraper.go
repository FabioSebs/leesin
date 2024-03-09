package scraper

import (
	"fmt"

	"github.com/FabioSebs/leesin/config"
	"github.com/FabioSebs/leesin/logger"
	"github.com/gocolly/colly"
)

var (
	cars     = make([]EV, 0)
	bookings = make([]Booking, 0)
)

type WebScraper interface {
	CollectorSetup() *colly.Collector
	// GetReviewsConcurrently(*colly.Collector) ([]EV, time.Duration)
	GetReviewsSynchronously(*colly.Collector)
}

type GoCollyProgram struct {
	Collector *colly.Collector
	Config    config.Config
	Logger    logger.Logger
}

func NewWebScraper() WebScraper {
	env := config.NewConfig()

	return &GoCollyProgram{
		Collector: colly.NewCollector(colly.AllowedDomains(
			env.AllowedDomains...,
		)),
		Config: env,
		Logger: logger.NewLogger(),
	}
}

func (g *GoCollyProgram) CollectorSetup() *colly.Collector {
	///////////////////////////////////////////////////// ORIGINAL COLLY /////////////////////////////////////////////////////
	g.Collector.OnHTML("div[data-stid='section-results'] div[data-stid='property-listing-results']", func(element *colly.HTMLElement) {
		element.ForEach("div.uitk-spacing div.uitk-card div.uitk-layout-grid div.uitk-card-content-section", func(_ int, h *colly.HTMLElement) {

		})

	})

	// Request Feedback
	g.Collector.OnRequest(func(r *colly.Request) {
		g.Logger.WriteTrace(fmt.Sprintf("visiting url: %s", r.URL.String()))
	})

	// Error Feedback
	g.Collector.OnError(func(_ *colly.Response, err error) {
		g.Logger.WriteError(fmt.Sprintf("error: %s", err.Error()))
	})
	return g.Collector
}

// func getMoreInfo(books []Booking) []Booking {
// 	var i int = 0
// 	env := config.NewConfig()
// 	l := logger.NewLogger()

// 	///////////////////////////////////////////////////// NESTED COLLY /////////////////////////////////////////////////////
// 	nestedColly := colly.NewCollector(colly.AllowedDomains(
// 		env.AllowedDomains...,
// 	))

// 	nestedColly.OnHTML("div#bodyconstraint div#bodyconstraint-inner div.k2-hp--gallery-header", func(element *colly.HTMLElement) {
// 		books[i].Address = element.ChildText("p.address span.hp_address_subtitle")
// 		books[i].PostCode = extractPostalCode(books[i].Address)
// 		i++
// 	})

// 	// Request Feedback
// 	nestedColly.OnRequest(func(r *colly.Request) {
// 		l.WriteTrace(fmt.Sprintf("visiting url: %s", r.URL.String()))
// 	})
// 	nestedColly.OnError(func(_ *colly.Response, err error) {
// 		l.WriteError(fmt.Sprintf("error: %s", err.Error()))
// 	})

// 	for _, val := range books {
// 		launchSecondVisit(nestedColly, val.Source)
// 	}
// 	return books
// }

// func launchSecondVisit(collector *colly.Collector, source string) {
// 	if err := collector.Visit(source); err != nil {
// 		fmt.Println(err.Error())
// 	}
// }

// func extractPostalCode(inputString string) string {
// 	// Define a regular expression pattern for matching postal codes
// 	postalCodePattern := regexp.MustCompile(`\b\d{5}\b`)

// 	// Find the first occurrence of the pattern in the input string
// 	match := postalCodePattern.FindString(inputString)

// 	// Return the extracted postal code if found, otherwise return an empty string
// 	return match
// }

// func (g *GoCollyProgram) GetReviewsConcurrently(collector *colly.Collector) ([]EV, time.Duration) {
// 	//empty slice
// 	defer emptyReviews(&cars)
// 	start := time.Now()

// 	//Visiting URLS
// 	jobNo, err := strconv.Atoi(g.Config.MaxPage)
// 	if err != nil {
// 		g.Logger.WriteError(err.Error())
// 	}

// 	var wg sync.WaitGroup
// 	wg.Add(jobNo)
// 	for i := 1; i <= jobNo; i++ {
// 		page := strconv.Itoa(i)
// 		go func(page string) {
// 			defer wg.Done()
// 			url := fmt.Sprintf(g.Config.FullDomain+"?page=%s&stars=1", page)
// 			if err := collector.Visit(url); err != nil {
// 				g.Logger.WriteError(err.Error())
// 			}
// 		}(page)
// 	}
// 	wg.Wait()
// 	// writeJSON(reviews)
// 	return cars, time.Since(start)
// }

func (g *GoCollyProgram) GetReviewsSynchronously(collector *colly.Collector) {
	if err := collector.Visit(g.Config.ExpediaDomain); err != nil {
		g.Logger.WriteError(err.Error())
	}

	collector.Wait()

	writeJSON(bookings, "expedia")
}
