package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type BookingDomain struct {
	URL string
}

type Config struct {
	FullDomain     string
	AllowedDomains []string
	MaxPage        string
}

func NewConfig() Config {

	return Config{
		FullDomain:     os.Getenv("TARGET"),
		AllowedDomains: []string{os.Getenv("ALLOWED1"), os.Getenv("ALLOWED2"), os.Getenv("ALLOWED3"), os.Getenv("ALLOWED4"), os.Getenv("ALLOWED5")},
		MaxPage:        os.Getenv("MAXPAGE"),
	}
}
