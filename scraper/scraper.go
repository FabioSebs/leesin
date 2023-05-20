package scraper

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/FabioSebs/leesin/config"
	"github.com/FabioSebs/leesin/logger"
	"github.com/FabioSebs/leesin/utils"
	"github.com/gocolly/colly"
)

type WebScraper interface {
	GetReviews()
}

type GoCollyProgram struct {
	Collector *colly.Collector
	Config    config.Config
	Logger    logger.Logger
}

func NewWebScraper() WebScraper {
	env := config.NewConfig()

	return &GoCollyProgram{
		Collector: colly.NewCollector(),
		Config:    env,
		Logger:    logger.NewLogger(),
	}
}

func (g *GoCollyProgram) GetReviews() {

	// Parsing Data to Model from URL
	reviews := make([]Review, 0)

	g.Collector.OnHTML(".data-business-unit-reviews-section", func(element *colly.HTMLElement) {
		info := element.DOM

		review := Review{
			User:   info.Find(".data-consumer-name-typography").Text(),
			Rating: "1",
			Review: info.Find(".data-review-content").Text(),
			Date:   info.Find(".data-service-review-date-of-experience-typography").Text(),
		}
		reviews = append(reviews, review)
	})

	// Request Feedback
	g.Collector.OnRequest(func(r *colly.Request) {
		g.Logger.WriteTrace(fmt.Sprintf("visiting url: %s", r.URL))
	})

	// Creating WaitGroups
	var wg sync.WaitGroup
	jobNo, err := strconv.Atoi(g.Config.MaxPage)
	if err != nil {
		g.Logger.WriteError(err.Error())
	}
	wg.Add(jobNo)

	//Visiting URLS
	for i := 1; i <= jobNo; i++ {
		go func(page string) {
			defer wg.Done()
			url := fmt.Sprintf(g.Config.FullDomain+"?page=%s&stars=1", page)
			g.Collector.Visit(url)
		}(strconv.Itoa(i))
	}
	wg.Wait()
	utils.PrettyPrintStruct(reviews)
}
