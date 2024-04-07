package scraper

import (
	"fmt"

	"github.com/FabioSebs/leesin/config"
	"github.com/FabioSebs/leesin/logger"
	"github.com/gocolly/colly"
)

var esectors []ESector

type ESectorScraper interface {
	WebScraper
}

type ESectorProgram struct {
	Collector *colly.Collector
	Config    config.Config
	Logger    logger.Logger
}

func NewESectorScraper() ESectorScraper {
	env := config.NewConfig()

	return &ESectorProgram{
		Collector: colly.NewCollector(colly.AllowedDomains(
			env.AllowedDomains...,
		)),
		Config: env,
		Logger: logger.NewLogger(),
	}
}

func (g *ESectorProgram) CollectorSetup() *colly.Collector {
	///////////////////////////////////////////////////// ORIGINAL COLLY /////////////////////////////////////////////////////
	g.Collector.OnHTML("main div.clearfix div#left-col-7 div#op-table-related table.dp-table tbody", func(e *colly.HTMLElement) {
		// Iterate through each row of the table
		e.ForEach("tr", func(_ int, row *colly.HTMLElement) {
			var sectorRow ESector
			// Extract data from each column of the row
			row.ForEach("td", func(i int, h *colly.HTMLElement) {
				switch i {
				case 0:
					sectorRow.Region = h.ChildText("a")
				case 1:
					sectorRow.Number = h.ChildText("span.c-red")
				case 3:
					sectorRow.Year = "2018"
				}
			})
			esectors = append(esectors, sectorRow)
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

func (g *ESectorProgram) GetReviewsSynchronously(collector *colly.Collector) {
	fmt.Println(g.Config.ICCTDomain)
	if err := collector.Visit(g.Config.ICCTDomain); err != nil {
		g.Logger.WriteError(fmt.Sprintf("error: %s", err.Error()))
	}

	writeESector(esectors, "esec")
}
