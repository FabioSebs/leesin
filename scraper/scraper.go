package scraper

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/FabioSebs/leesin/config"
	"github.com/FabioSebs/leesin/logger"
	"github.com/gocolly/colly"
	"github.com/tealeg/xlsx"
)

var (
	cars         = make([]EV, 0)
	publications = make([]PublicationProject, 0)
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
			pub := PublicationProject{
				Title:  h.ChildText("h3"),
				Format: h.ChildText("div.article_content div.tax_term"),
				Year:   h.ChildText("p.post_meta"),
				Link:   h.ChildAttr("h3 a", "href"),
			}

			h.ForEach("div.authors", func(i int, e *colly.HTMLElement) {
				switch i {
				case 0:
					pub.Author1 = e.ChildText("a")
				case 1:
					pub.Author2 = e.ChildText("a")
				case 2:
					pub.Author3 = e.ChildText("a")
				case 3:
					pub.Author4 = e.ChildText("a")
				case 4:
					pub.Author5 = e.ChildText("a")
				case 5:
					pub.Author6 = e.ChildText("a")
				}
			})

			if strings.Contains(strings.ToLower(pub.Title), "battery") ||
				strings.Contains(strings.ToLower(pub.Title), "cost") ||
				strings.Contains(strings.ToLower(pub.Title), "econ") ||
				strings.Contains(strings.ToLower(pub.Title), "perform") ||
				strings.Contains(strings.ToLower(pub.Title), "econ") {
				publications = append(publications, pub)
			}
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
	for i := 1; i < 11; i++ {
		if err := collector.Visit(fmt.Sprintf("https://theicct.org/insight-analysis/publications/?_year=2024-02-01%%2C2024-10-14&_paged=%d", i)); err != nil {
			g.Logger.WriteError(err.Error())
		}
		writeExcel(publications)
	}

	return cars, time.Since(start)
}

// func writeJSON(data []Publication, fname string) {
// 	cardata, err := json.MarshalIndent(data, "", " ")
// 	if err != nil {
// 		log.Println("Unable to create json file")
// 		return
// 	}

// 	if err = ioutil.WriteFile(fmt.Sprintf("%s.json", fname), cardata, 0644); err != nil {
// 		log.Println("unable to write to json file")
// 	}
// 	cars = cars[:0]
// }

func writeExcel(data []PublicationProject) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Projects")
	if err != nil {
		log.Fatalf("Failed to create sheet : %v", err)
	}

	header := sheet.AddRow()
	headerData := []string{
		"Year", "Sector", "Vehicle Type", "Region", "Metric", "Format", "Title",
		"Author 1", "Author 2", "Author 3", "Author 4", "Author 5", "Author 6", "Link",
	}

	for _, h := range headerData {
		cell := header.AddCell()
		cell.Value = h
	}

	for _, project := range data {
		row := sheet.AddRow()
		row.AddCell().Value = project.Year
		row.AddCell().Value = project.Sector
		row.AddCell().Value = project.VehicleType
		row.AddCell().Value = project.Region
		row.AddCell().Value = project.Metric
		row.AddCell().Value = project.Format
		row.AddCell().Value = project.Title
		row.AddCell().Value = project.Author1
		row.AddCell().Value = project.Author2
		row.AddCell().Value = project.Author3
		row.AddCell().Value = project.Author4
		row.AddCell().Value = project.Author5
		row.AddCell().Value = project.Author6
		row.AddCell().Value = project.Link
	}

	err = file.Save("projects.xlsx")
	if err != nil {
		log.Fatalf("Failed to save file: %v", err)
	}
}
