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
	GoogleData = make([]GoogleMapsData, 0)
)

type WebScraper interface {
	CollectorSetup() *colly.Collector
	GetData(*colly.Collector, string) ([]GoogleMapsData, time.Duration)
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
	g.Collector.OnHTML("div.XltNde.tTVLSc", func(element *colly.HTMLElement) {
		// Target the "sub parent"
		element.ForEach("div.m6QErb.DxyBCb", func(_ int, sub *colly.HTMLElement) {

			// Handle normal children
			sub.ForEach("div.TFQHme", func(i int, h *colly.HTMLElement) {
				// Do your scraping here for TFQHme children
				fmt.Println("TFQHme child:", h.Text)
			})

			// Handle the odd first child
			sub.ForEach("div.Nv2PK.THOPZb.CpccDe", func(i int, h *colly.HTMLElement) {
				fmt.Println("Odd child:", h.Text)
			})
		})
	})

	// Request logging
	g.Collector.OnRequest(func(r *colly.Request) {
		g.Logger.WriteTrace(fmt.Sprintf("visiting url: %s", r.URL.String()))
	})

	// Error handling
	g.Collector.OnError(func(_ *colly.Response, err error) {
		g.Logger.WriteError(fmt.Sprintf("error: %s", err.Error()))
	})

	return g.Collector
}

func (g *GoCollyProgram) GetData(collector *colly.Collector, query string) ([]GoogleMapsData, time.Duration) {
	start := time.Now()

	if err := collector.Visit(fmt.Sprintf(g.Config.URL, query)); err != nil {
		g.Logger.WriteError(err.Error())
	}

	writeJSON(GoogleData, "googlemaps")

	return GoogleData, time.Since(start)
}

func writeJSON(data []GoogleMapsData, fname string) {
	mapsdata, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	if err = ioutil.WriteFile(fmt.Sprintf("%s.json", fname), mapsdata, 0644); err != nil {
		log.Println("unable to write to json file")
	}
	GoogleData = GoogleData[:0]
}
