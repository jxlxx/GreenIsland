package world

import (
	"github.com/jxlxx/GreenIsland/bank"
)

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
	Name            string            `yaml:"name"`
	Code            string            `yaml:"code"`
	Currency        bank.CurrencyCode `yaml:"currency_code"`
	CentralBank     CentralBank       `yaml:"central_bank"`
	CommercialBanks []CommercialBank  `yaml:"commercial_banks"`
	Population      Population        `yaml:"population"`

	unemployment       int
	consumerPriceIndex int
}

type Population struct {
	Total   Value `yaml:"total"`
	Working Value `yaml:"working"`
}

type CentralBank struct {
	Name    string        `yaml:"name"`
	Reserve CurrencyValue `yaml:"reserve"`

	m1 int
	m2 int
	m3 int
}

type CommercialBank struct {
	Name     string        `yaml:"name"`
	Deposits CurrencyValue `yaml:"deposits"`
	Reserve  CurrencyValue `yaml:"reserve"`
}
