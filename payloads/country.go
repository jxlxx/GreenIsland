package payloads

import "github.com/jxlxx/GreenIsland/bank"

type QuarterlyCountryUpdate struct {
	Name              string
	Quarter           int
	TotalPopulation   int
	WorkingPopulation int
	MoneySupply       MoneySupply
}

type MoneySupply struct {
	CentralBankName string
	Currency        bank.CurrencyCode
	CurrencyUnit    bank.UnitType
	M1              int
	M2              int
	M3              int
}
