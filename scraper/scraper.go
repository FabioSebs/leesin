package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/FabioSebs/leesin/config"
	"github.com/FabioSebs/leesin/logger"
	"github.com/gocolly/colly"
)

var (
	cars         = make([]EV, 0)
	publications = make([]Publication, 0)
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
	g.Collector.OnHTML("div.facetwp-template ", func(element *colly.HTMLElement) {
		element.ForEach("div.article_content", func(_ int, h *colly.HTMLElement) {
			pub := Publication{
				Title:  h.ChildText("h3"),
				Source: h.ChildAttr("h3 a", "href"), // Corrected line
				Year:   h.ChildText("p.post_meta"),
			}
			h.ForEach("div.authors", func(_ int, e *colly.HTMLElement) {
				pub.Author = append(pub.Author, e.ChildText("a"))
			})
			publications = append(publications, pub)
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
	for i := 1; i < 7; i++ {
		if err := collector.Visit(fmt.Sprintf("https://theicct.org/insight-analysis/publications/?_icct_authors=253&_paged=%d&_sort=date_desc", i)); err != nil {
			g.Logger.WriteError(err.Error())
		}
		writeJSON(publications, "publications")
	}

	return cars, time.Since(start)
}

func writeJSON(data []Publication, fname string) {
	cardata, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	if err = ioutil.WriteFile(fmt.Sprintf("%s.json", fname), cardata, 0644); err != nil {
		log.Println("unable to write to json file")
	}
	cars = cars[:0]
}
