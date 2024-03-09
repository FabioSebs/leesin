package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/FabioSebs/leesin/config"
	"github.com/FabioSebs/leesin/logger"
	"github.com/gocolly/colly"
)

var ExpediaEntries []Booking
var newEntries []Booking
var x int = 0

type ExpediaScraper interface {
	WebScraper
}

type ExpediaProgram struct {
	Collector *colly.Collector
	Config    config.Config
	Logger    logger.Logger
}

func NewExpediaScraper() ExpediaScraper {
	env := config.NewConfig()

	return &ExpediaProgram{
		Collector: colly.NewCollector(colly.AllowedDomains(
			env.AllowedDomains...,
		)),
		Config: env,
		Logger: logger.NewLogger(),
	}
}

func (g *ExpediaProgram) CollectorSetup() *colly.Collector {
	///////////////////////////////////////////////////// ORIGINAL COLLY /////////////////////////////////////////////////////
	g.Collector.OnHTML("div.uitk-layout-grid div.uitk-layout-grid-item div.uitk-spacing.uitk-spacing-margin-large-inlinestart-three.uitk-spacing-margin-large-blockstart-three div.uitk-layout-flex div[data-stid='map-image-link'] div div div", func(element *colly.HTMLElement) {
		addy := element.Text
		if addy != "" {
			ExpediaEntries[x].Address = addy
			ExpediaEntries[x].PostCode = addy[len(addy)-5:]
		}

		newEntries = append(newEntries, ExpediaEntries[x])
		writeJSON(newEntries, "newexpedia")
		x++
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

func (g *ExpediaProgram) GetReviewsSynchronously(collector *colly.Collector) {
	// Read the Data File
	data, err := ioutil.ReadFile("expedia.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Unmarshal the JSON data into the slice
	err = json.Unmarshal(data, &ExpediaEntries)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	fmt.Println(len(ExpediaEntries))

	// Concurrently Visit all the entries for address
	for idx, val := range ExpediaEntries {
		fmt.Println(idx)
		if err := collector.Visit(val.Source); err != nil {
			g.Logger.WriteError(err.Error())
		}
	}
}
