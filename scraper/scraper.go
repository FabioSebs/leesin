package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func (g *GoCollyProgram) CollectorSetup() *colly.Collector {
	g.Collector.OnHTML("div.d4924c9e74 div.c82435a4b8", func(element *colly.HTMLElement) {
		var book Booking
		element.ForEach("div.c066246e13", func(_ int, h *colly.HTMLElement) {
			book.Title = h.ChildText("div.f6431b446c")
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
			location := h.ChildText("div.abf093bdfe span.aee5343fdb")
			result := strings.Replace(location, "Show on map", "", -1)
			result = strings.TrimSpace(result)
			book.Location = result
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
			review, _ := strconv.ParseFloat(h.ChildText("div.aca0ade214 div.a3b8729ab1.d86cee9b25"), 32)
			book.Review = float32(review)

			book.Source = h.ChildAttr("a.a78ca197d0", "href")

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
		if err := collector.Visit(fmt.Sprintf(g.Config.BookingDomain, i)); err != nil {
			g.Logger.WriteError(err.Error())
		}
		writeJSON(bookings, "balidata")
	}

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
