package world

import (
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/jxlxx/GreenIsland/bank"
	"github.com/jxlxx/GreenIsland/payloads"
	"github.com/jxlxx/GreenIsland/subjects"
	"github.com/jxlxx/GreenIsland/types"
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
	CommercialBanks []*bank.Bank      `yaml:"commercial_banks"`
	Population      Population        `yaml:"population"`

	nc *nats.EncodedConn
}

type Population struct {
	Total   types.Value `yaml:"total"`
	Working types.Value `yaml:"working"`
}

type CentralBank struct {
	Name    string             `yaml:"name"`
	Reserve bank.CurrencyValue `yaml:"reserve"`
}

func (c *Country) CreateBanks() {
	for _, b := range c.CommercialBanks {
		b.Setup()
	}
}

func (c *Country) Initialize() {
	for _, b := range c.CommercialBanks {
		b.Init()
	}
}

func (c *Country) DailySubscriber() func(payloads.WorldTick) {
	return func(payloads.WorldTick) {
		c.DailyUpdate()
	}
}

func (c *Country) QuarterlySubscriber() func(payloads.WorldTick) {
	return func(p payloads.WorldTick) {

		update := payloads.QuarterlyCountryUpdate{
			Name:              c.Name,
			Quarter:           p.Quarter,
			TotalPopulation:   c.Population.Total.Value,
			WorkingPopulation: c.Population.Total.Value,
			MoneySupply:       c.CalculateMoneySupply(),
		}

		if err := c.nc.Publish(subjects.QuarterlyCountryUpdate(c.Code, p.Quarter), update); err != nil {
			fmt.Println(err)
		}
	}
}

func (c *Country) DailyUpdate() {
	c.Population.Total.Value += c.Population.Total.CalcUpdate()
	c.Population.Working.Value += c.Population.Working.CalcUpdate()
	c.CentralBank.Reserve.Value += c.CentralBank.Reserve.CalcUpdate()

}
