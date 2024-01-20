package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	GetReviewsSynchronously(*colly.Collector) ([]EV, time.Duration)
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

// TODO: STARS
// TODO: VISIT NESTED URL (CREATE NEW COLLECTOR SET UP SEPERATE OnHTML and VISIT)
// TODO: GET FULL ADDRESS
// TODO: GET PRICE
func (g *GoCollyProgram) CollectorSetup() *colly.Collector {
	///////////////////////////////////////////////////// ORIGINAL COLLY /////////////////////////////////////////////////////
	g.Collector.OnHTML("div.d4924c9e74 div.c82435a4b8", func(element *colly.HTMLElement) {
		element.ForEach("div.c066246e13", func(_ int, h *colly.HTMLElement) {
			var book Booking
			book.Title = h.ChildText("div.c1edfbabcb div.d6767e681c h3.aab71f8e4e")
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
			location := h.ChildText("div.abf093bdfe span.aee5343fdb")
			result := strings.Replace(location, "Show on map", "", -1)
			result = strings.TrimSpace(result)
			book.Location = result
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
			review, _ := strconv.ParseFloat(h.ChildText("div.aca0ade214 div.a3b8729ab1.d86cee9b25"), 32)
			book.Review = float32(review)
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
			book.Source = h.ChildAttr("a.a78ca197d0", "href")
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
			book.Stars = h.ChildAttr("div.d8c86a593f div.b3f3c831be", "aria-label")
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
			reviewNo := strings.Split(h.ChildText("div.aca0ade214.ebac6e22e9.cd2e7d62b0.a0ff1335a1 div.abf093bdfe.f45d8e4c32.d935416c47"), " ")
			book.ReviewNumber, _ = strconv.Atoi(reviewNo[0])
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

			bookings = append(bookings, book)
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

func getMoreInfo(books []Booking) []Booking {
	var i int = 0
	env := config.NewConfig()
	l := logger.NewLogger()

	///////////////////////////////////////////////////// NESTED COLLY /////////////////////////////////////////////////////
	nestedColly := colly.NewCollector(colly.AllowedDomains(
		env.AllowedDomains...,
	))

	nestedColly.OnHTML("div#bodyconstraint div#bodyconstraint-inner div.k2-hp--gallery-header", func(element *colly.HTMLElement) {
		books[i].Address = element.ChildText("p.address span.hp_address_subtitle")
		books[i].PostCode = extractPostalCode(books[i].Address)
		i++
	})

	// Request Feedback
	nestedColly.OnRequest(func(r *colly.Request) {
		l.WriteTrace(fmt.Sprintf("visiting url: %s", r.URL.String()))
	})
	nestedColly.OnError(func(_ *colly.Response, err error) {
		l.WriteError(fmt.Sprintf("error: %s", err.Error()))
	})

	for _, val := range books {
		launchSecondVisit(nestedColly, val.Source)
	}
	return books
}

func launchSecondVisit(collector *colly.Collector, source string) {
	if err := collector.Visit(source); err != nil {
		fmt.Println(err.Error())
	}
}

func extractPostalCode(inputString string) string {
	// Define a regular expression pattern for matching postal codes
	postalCodePattern := regexp.MustCompile(`\b\d{5}\b`)

	// Find the first occurrence of the pattern in the input string
	match := postalCodePattern.FindString(inputString)

	// Return the extracted postal code if found, otherwise return an empty string
	return match
}

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

func (g *GoCollyProgram) GetReviewsSynchronously(collector *colly.Collector) ([]EV, time.Duration) {

	start := time.Now()

	//Visiting URLS
	for i := 0; i <= 975; i += 25 {
		fmt.Println(i)
		if err := collector.Visit(fmt.Sprintf(g.Config.BookingDomain, i)); err != nil {
			g.Logger.WriteError(err.Error())
		}
	}

	newbooks := getMoreInfo(bookings)

	writeJSON(newbooks, "balidata")
	return cars, time.Since(start)
}

func writeJSON(data []Booking, fname string) {
	balidata, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	if err = ioutil.WriteFile(fmt.Sprintf("%s.json", fname), balidata, 0644); err != nil {
		log.Println("unable to write to json file")
	}
	cars = cars[:0]
}
