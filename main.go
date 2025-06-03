package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/FabioSebs/leesin/scraper"
)

func main() {
	// prompt
	fmt.Print("Enter search query: ")
	reader := bufio.NewReader(os.Stdin)
	query, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// trim new lines
	query = query[:len(query)-1]

	// scrape
	ws := scraper.NewWebScraper()
	collector := ws.CollectorSetup()
	ws.GetData(collector, query)
}
