package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type BookingDomain struct {
	URL string
}

type Config struct {
	FullDomain     BookingDomain
	AllowedDomains []string
	MaxPage        string
	BookingDomain  string
	HotelDomain    string
	ExpediaDomain  string
	ICCTDomain     string
}

func NewConfig() Config {
	return Config{
		FullDomain: BookingDomain{
			URL: os.Getenv("BOOKINGDOMAIN"),
		},
		AllowedDomains: []string{os.Getenv("ALLOWED1"), os.Getenv("ALLOWED2")},
		MaxPage:        os.Getenv("MAXPAGE"),
		BookingDomain:  os.Getenv("BOOKINGDOMAIN"),
		HotelDomain:    os.Getenv("HOTELDOMAIN"),
		ExpediaDomain:  os.Getenv("EXPEDIADOMAIN"),
		ICCTDomain:     os.Getenv("ICCTDOMAIN"),
	}
}
