package scraper

type EV struct {
	Name  string
	Price string
}

// body div.bodyconstraint--full-width div.bodyconstraint-inner div.af5895d4b2 div.df7e6ba27d div.bcbf33c5c3 div.d4924c9e74 div.c82435a4b8
type Booking struct {
	Title    string
	Location string
	Review   float32
	Source   string
}
