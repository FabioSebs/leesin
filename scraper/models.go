package scraper

type EV struct {
	Name  string
	Price string
}

type Publication struct {
	Author []string
	Year   string
	Title  string
	Source string
}

// for icct-projects

type PublicationProject struct {
	Year        string `xlsx:"Publication Year"`
	Sector      string `xlsx:"Sector"`
	VehicleType string `xlsx:"Vehicle Type"`
	Region      string `xlsx:"Region"`
	Metric      string `xlsx:"Metric"`
	Format      string `xlsx:"Format"`
	Title       string `xlsx:"Title"`
	Author1     string `xlsx:"Author 1"`
	Author2     string `xlsx:"Author 2"`
	Author3     string `xlsx:"Author 3"`
	Author4     string `xlsx:"Author 4"`
	Author5     string `xlsx:"Author 5"`
	Author6     string `xlsx:"Author 6"`
	Link        string `xlsx:"Link"`
}
