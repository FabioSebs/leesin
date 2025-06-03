package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	URL            string
	AllowedDomains []string
	MaxPage        string
	ICCTDomain     string
}

func NewConfig() Config {
	return Config{
		URL:            os.Getenv("URL"),
		AllowedDomains: []string{os.Getenv("ALLOWED1"), os.Getenv("ALLOWED2"), os.Getenv("ALLOWED3"), os.Getenv("ALLOWED4")},
		MaxPage:        os.Getenv("MAXPAGE"),
		ICCTDomain:     os.Getenv("ICCTDOMAIN"),
	}
}
