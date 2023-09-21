package world

import "github.com/jxlxx/GreenIsland/bank"

type Value struct {
	Total         int `yaml:"total"`
	Jitter        int `yaml:"jitter"`
	JitterAverage int `yaml:"jitter_average"`
}

type CurrencyValue struct {
	Currency      bank.CurrencyCode `yaml:"currency"`
	Unit          bank.UnitType     `yaml:"currency_unit"`
	Total         int               `yaml:"total"`
	Jitter        int               `yaml:"jitter"`
	JitterAverage int               `yaml:"jitter_average"`
}
