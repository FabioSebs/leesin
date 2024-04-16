package scraper

import (
	"fmt"
	"strconv"

	"github.com/FabioSebs/leesin/config"
	"github.com/FabioSebs/leesin/logger"
	"github.com/gocolly/colly"
)

var (
	data = make([]EVModel, 0)
)

type WebScraper interface {
	CollectorSetup() *colly.Collector
	StartScraper(*colly.Collector)
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
	g.Collector.OnHTML("main div.content div.list", func(element *colly.HTMLElement) {
		element.ForEach("div.list-item div.data-wrapper", func(_ int, h *colly.HTMLElement) {
			var (
				model EVModel
			)
			//Initializing
			model.Name = h.ChildText("div.title-wrap h2")

			model.Acceleration = h.ChildText("div.specs p.left span.acceleration")
			model.TopSpeed = h.ChildText("div.specs p.left span.topspeed")
			model.Range = h.ChildText("div.specs p.left span.erange_real")
			model.Efficiency = h.ChildText("div.specs p.left span.efficiency")
			model.FastCharge = h.ChildText("div.specs p.left span.fastcharge_speed_print")
			model.Price = h.ChildText("div.pricing span.price_buy span")

			// Appending
			data = append(data, model)
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

func (g *GoCollyProgram) StartScraper(collector *colly.Collector) {
	max, err := strconv.Atoi(g.Config.MaxPage)
	if err != nil {
		g.Logger.WriteError(fmt.Sprintf("error: %s", err.Error()))
	}

	for i := 0; i <= max; i++ {
		var url string
		var page string = strconv.Itoa(i)

		url = fmt.Sprintf(g.Config.FullDomain, page)

		if err := collector.Visit(url); err != nil {
			g.Logger.WriteError(err.Error())
		}
	}

	collector.Wait()

	writeJSON(data, "ev-latest")
}
