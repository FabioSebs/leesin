package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func removeExtraWhitespace(s string) string {
	words := strings.Fields(s)
	return strings.Join(words, " ")
}

func writeJSON(data []Booking, fname string) {
	balidata, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	if err = ioutil.WriteFile(fmt.Sprintf("%s.json", fname), balidata, 0644); err != nil {
		log.Println("unable to write to json file")
	}
	cars = cars[:0]
}

func writePubs(pubs []Publication, fname string) {
	pubdata, err := json.MarshalIndent(pubs, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s.json", fname), pubdata, 0644); err != nil {
		log.Println("unable to write to json file")
	}
}
