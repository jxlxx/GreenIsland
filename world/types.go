package world

import "github.com/jxlxx/GreenIsland/bank"

type Industry string

const (
	Constuction    Industry = "construction"
	Agriculture    Industry = "agriculture"
	Healthcare     Industry = "health_care"
	Food           Industry = "food"
	Manufacturing  Industry = "manufacturing"
	Retail         Industry = "retail"
	Transportation Industry = "transportation"
	Mining         Industry = "mining"
	Energy         Industry = "energy"
)

type Country struct {
	Name        string
	CentralBank CentralBank

	Unemployment       Rate
	ConsumerPriceIndex int
}

type CentralBank struct {
	Name     string
	Currency bank.Currency
	M1       int
	M2       int
	M3       int
}

type Rate struct {
	Rate   int
	Jitter int
}
