package scraper

import (
	"fmt"

	"github.com/FabioSebs/leesin/config"
	"github.com/FabioSebs/leesin/logger"
	"github.com/gocolly/colly"
	"github.com/tealeg/xlsx"
)

var TripBookings []Booking

type TripAdScraper interface {
	WebScraper
}

type TripAdProgram struct {
	Collector *colly.Collector
	Config    config.Config
	Logger    logger.Logger
	Data      []Booking
	Idx       int
}

func NewTripScraper() ExpediaScraper {
	env := config.NewConfig()

	return &TripAdProgram{
		Collector: colly.NewCollector(colly.AllowedDomains(
			env.AllowedDomains...,
		)),
		Config: env,
		Logger: logger.NewLogger(),
		Data:   []Booking{},
		Idx:    0,
	}
}

func (g *TripAdProgram) CollectorSetup() *colly.Collector {

	///////////////////////////////////////////////////// ORIGINAL COLLY /////////////////////////////////////////////////////
	g.Collector.OnHTML("", func(element *colly.HTMLElement) {

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

func (g *TripAdProgram) GetReviewsSynchronously(collector *colly.Collector) {

	// Read the Data File
	file, err := xlsx.OpenFile("tripadvisor.xlsx")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	sheet := file.Sheets[0]

	// Assuming the order of columns in Excel matches the order in the struct
	for _, row := range sheet.Rows {
		// Create an instance of your struct to hold row data
		var entry Booking

		// Loop through the cells in each row
		for i, cell := range row.Cells {
			// Assign cell value to the corresponding struct field based on the index
			switch i {
			case 0:
				entry.Title = cell.String()
			case 1:
				entry.Location = cell.String()
			case 2:
				entry.Address = cell.String()
			case 3:
				entry.PostCode = cell.String()
			case 4:
				entry.Review = cell.String()
			case 5:
				entry.ReviewNumber = cell.String()
			case 6:
				entry.Price = cell.String()
			case 7:
				entry.PropertyType = cell.String()
			case 8:
				entry.Source = cell.String()
			}
		}

		// Append the struct to the data slice
		g.Data = append(g.Data, entry)
	}

	// Concurrently Visit all the entries for address
	for _, val := range g.Data {
		if err := g.Collector.Visit(val.Source); err != nil {
			g.Logger.WriteError(err.Error())
		}
	}
}
