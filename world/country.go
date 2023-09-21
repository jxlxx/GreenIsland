package world

import (
	"fmt"

	"github.com/jxlxx/GreenIsland/bank"
	"github.com/jxlxx/GreenIsland/payloads"
	"github.com/jxlxx/GreenIsland/subjects"
	"github.com/nats-io/nats.go"
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
	nc                 *nats.EncodedConn
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

func (b CommercialBank) Update() CommercialBank {
	b.Deposits.Value += b.Deposits.CalcUpdate()
	b.Reserve.Value += b.Reserve.CalcUpdate()
	return b
}

func (c Country) CalculateMoneySupply() payloads.MoneySupply {
	m1 := c.CalculateM1()
	m2 := c.CalculateM2()
	m3 := c.CalculateM3()
	return payloads.MoneySupply{
		Currency:     c.CentralBank.Reserve.Currency,
		CurrencyUnit: c.CentralBank.Reserve.Unit,
		M1:           m1,
		M2:           m1 + m2,
		M3:           m1 + m2 + m3,
	}
}

func (c Country) CalculateM1() int {
	sum := 0
	for _, b := range c.CommercialBanks {
		sum += b.Deposits.Value
	}
	return sum
}
func (c Country) CalculateM2() int {
	sum := 0
	// get info from bank (traders)
	return sum
}
func (c Country) CalculateM3() int {
	sum := 0
	// get info from bank (companies)
	return sum
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

	for i, b := range c.CommercialBanks {
		c.CommercialBanks[i] = b.Update()
	}
}
