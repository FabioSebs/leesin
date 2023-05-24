package wordcounter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/FabioSebs/leesin/config"
	"github.com/FabioSebs/leesin/logger"
	"github.com/FabioSebs/leesin/scraper"
)

type WordCounter interface {
	CountReviews([]scraper.Review)
	CountSubjects([]scraper.Review)
}

type Counter struct {
	Words  map[string]int
	Logger logger.Logger
	Config config.Config
}

type Count struct {
	Key   string
	Value int
}

type CountList []Count

func (c CountList) Len() int { return len(c) }

func (c CountList) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

func (c CountList) Less(i, j int) bool { return c[i].Value < c[j].Value }

func NewWordCounter() WordCounter {
	return &Counter{
		Words:  make(map[string]int),
		Logger: logger.NewLogger(),
		Config: config.NewConfig(),
	}
}

func (c *Counter) CountReviews(reviews []scraper.Review) {
	defer func(words *map[string]int) {
		for k := range *words {
			delete(*words, k)
		}
	}(&c.Words)

	// Tokenizing
	for _, val := range reviews {
		words := strings.Split(val.Review, " ")
		for _, word := range words {
			token := strings.ToLower(word)
			_, exists := c.Words[token]
			if exists {

				c.Words[token] += 1
			} else {
				c.Words[token] = 1
			}
		}
	}

	// Sorting
	countList := make(CountList, len(c.Words))
	for k, v := range c.Words {
		countList = append(countList, Count{Key: k, Value: v})
	}

	sort.Sort(sort.Reverse(countList))

	fmt.Println(countList)
}

func (c *Counter) CountSubjects([]scraper.Review) {

}
