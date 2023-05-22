package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

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
		Collector: colly.NewCollector(colly.AllowedDomains(
			env.AllowedDomains...,
		)),
		Config: env,
		Logger: logger.NewLogger(),
	}
}

func (g *GoCollyProgram) GetReviews() {
	// Parsing Data to Model from URL
	reviews := make([]Review, 0)

	g.Collector.OnHTML("div.styles_mainContent__nFxAv section.styles_reviewsContainer__3_GQw", func(element *colly.HTMLElement) {
		element.ForEach("div.styles_cardWrapper__LcCPA", func(i int, h *colly.HTMLElement) {
			review := Review{
				User:   h.DOM.Find("div.styles_reviewCardInner__EwDq2 aside.styles_consumerInfoWrapper__KP3Ra div.styles_consumerDetailsWrapper__p2wdr a span").First().Text(),
				Rating: "1",
				Review: h.DOM.Find("div.styles_reviewCardInner__EwDq2 section.styles_reviewContentwrapper__zH_9M div.styles_reviewContent__0Q2Tg p.typography_body-l__KUYFJ").Text(),
				Date:   h.DOM.Find("div.styles_reviewCardInner__EwDq2 section.styles_reviewContentwrapper__zH_9M div.styles_reviewContent__0Q2Tg p.typography_body-m__xgxZ_").Text(),
			}
			reviews = append(reviews, review)
		})
	})

	// Request Feedback
	g.Collector.OnRequest(func(r *colly.Request) {
		g.Logger.WriteTrace(fmt.Sprintf("visiting url: %s", r.URL.String()))
	})

	// Error Feedback
	g.Collector.OnError(func(r *colly.Response, err error) {
		g.Logger.WriteError(fmt.Sprintf("error: %s", err.Error()))
	})

	//Visiting URLS

	jobNo, err := strconv.Atoi(g.Config.MaxPage)
	if err != nil {
		g.Logger.WriteError(err.Error())
	}

	for i := 1; i <= jobNo; i++ {
		page := strconv.Itoa(i)

		url := fmt.Sprintf(g.Config.FullDomain+"?page=%s&stars=1", page)
		if err := g.Collector.Visit(url); err != nil {
			g.Logger.WriteError(err.Error())
		}
	}

	utils.PrettyPrintStruct(reviews)
	writeJSON(reviews)
}

func writeJSON(data []Review) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = ioutil.WriteFile("leaguereviews.json", file, 0644)
}
