package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type LatestEVDomains struct {
	General     string
	MPV         string
	SUV         string
	Crossover   string
	Hatchback   string
	Sedan       string
	PickupTruck string
	Coupe       string
	Wagon       string
	Motorcycles string
}

type Config struct {
	FullDomain     LatestEVDomains
	AllowedDomains []string
	MaxPage        string
	ICCTDomain     string
}

func NewConfig() Config {
	return Config{
		FullDomain: LatestEVDomains{
			General:     os.Getenv("GENERAL"),
			MPV:         os.Getenv("MPV"),
			SUV:         os.Getenv("SUV"),
			Crossover:   os.Getenv("CROSSOVER"),
			Hatchback:   os.Getenv("HATCHBACK"),
			Sedan:       os.Getenv("SEDAN"),
			PickupTruck: os.Getenv("PICKUPTRUCK"),
			Coupe:       os.Getenv("COUPE"),
			Wagon:       os.Getenv("WAGON"),
			Motorcycles: os.Getenv("MOTORS"),
		},
		AllowedDomains: []string{os.Getenv("ALLOWED1"), os.Getenv("ALLOWED2")},
		MaxPage:        os.Getenv("MAXPAGE"),
		ICCTDomain:     os.Getenv("ICCTDOMAIN"),
	}
}
