package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/FabioSebs/leesin/config"
	"github.com/FabioSebs/leesin/logger"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var PubList []Publication

type PublicationScraper interface {
	WebScraper
}

type PublicationProgram struct {
	Collector *colly.Collector
	Config    config.Config
	Logger    logger.Logger
	Data      []Publication
	Idx       int
}

func NewPubScraper() ExpediaScraper {
	env := config.NewConfig()

	return &PublicationProgram{
		Collector: colly.NewCollector(colly.AllowedDomains(
			env.AllowedDomains...,
		)),
		Config: env,
		Logger: logger.NewLogger(),
		Data:   []Publication{},
		Idx:    0,
	}
}

func (g *PublicationProgram) CollectorSetup() *colly.Collector {

	///////////////////////////////////////////////////// ORIGINAL COLLY /////////////////////////////////////////////////////
	g.Collector.OnHTML("div#page-container div#et-boc div#et-main-area div#main-content div.container div#content-area div.et_builder_inner_content", func(element *colly.HTMLElement) {
		var pub Publication = Publication{
			Year:    g.Data[g.Idx].Year,
			Format:  g.Data[g.Idx].Format,
			Title:   g.Data[g.Idx].Title,
			Author1: g.Data[g.Idx].Author1,
			Author2: g.Data[g.Idx].Author2,
			Author3: g.Data[g.Idx].Author3,
			Author4: g.Data[g.Idx].Author4,
			Author5: g.Data[g.Idx].Author5,
			Author6: g.Data[g.Idx].Author6,
			Link:    g.Data[g.Idx].Link,
		}

		// Select the parent element
		parentElement := element.DOM.Find("div")

		if parentElement.Length() == 0 {
			fmt.Println("No parent element found.")
			return
		}

		// Iterate over child elements with class "et_pb_column"
		parentElement.Find("div.et_pb_row div.et_pb_column").Each(func(i int, childElement *goquery.Selection) {
			fmt.Println("Child Element:", i)

			// Access data within each child element
			tagLabel := childElement.Find("div.et_pb_module div.et_pb_text_inner div.tag label").Text()
			fmt.Println("Tag Label:", tagLabel)

			if tagLabel == "Sector" {
				pub.Sector = childElement.Find("div.et_pb_module div.et_pb_text_inner div.tag a").Text()
			}

			if tagLabel == "Region" {
				pub.Region = childElement.Find("div.et_pb_module div.et_pb_text_inner div.tag a").Text()
			}
		})

		PubList = append(PubList, pub)
		writePubs(PubList, "publications")
		g.Idx = g.Idx + 1
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

func (g *PublicationProgram) GetReviewsSynchronously(collector *colly.Collector) {

	// Read the Data File
	data, err := ioutil.ReadFile("pubs.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Unmarshal the JSON data into the slice
	err = json.Unmarshal(data, &g.Data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Concurrently Visit all the entries for address
	for _, val := range g.Data {
		if err := g.Collector.Visit(val.Link); err != nil {
			g.Logger.WriteError(err.Error())
		}
	}
}
