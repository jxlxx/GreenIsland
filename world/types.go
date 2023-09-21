package world

import (
	"math/rand"

	"github.com/jxlxx/GreenIsland/bank"
)

type Value struct {
	Value   int `yaml:"value"`
	Jitter  int `yaml:"jitter"`
	Average int `yaml:"average_delta"`
}

type CurrencyValue struct {
	Currency bank.CurrencyCode `yaml:"currency"`
	Unit     bank.UnitType     `yaml:"currency_unit"`
	Value    int               `yaml:"value"`
	Jitter   int               `yaml:"jitter"`
	Average  int               `yaml:"average_delta"`
}

func (v Value) CalcUpdate() int {
	if v.Jitter == 0 {
		return 0
	}
	delta := rand.Intn(v.Jitter) - v.Jitter/2
	shift := v.Average
	return delta + shift

}

func (v CurrencyValue) CalcUpdate() int {
	if v.Jitter == 0 {
		return 0
	}
	delta := rand.Intn(v.Jitter) - v.Jitter/2
	shift := v.Average
	return delta + shift
}
