package world

import (
	"github.com/jxlxx/GreenIsland/payloads"
)

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
	// get info from bank
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
