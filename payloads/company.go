package payloads

import "github.com/jxlxx/GreenIsland/bank"

type QuarterlyCompanyUpdate struct {
	Name         string `json:"name"`
	CurrencyCode bank.CurrencyCode
	Quarter      int
	Employees    int
	BalanceSheet BalanceSheet
	IncomeSheet  Income
	Dividends    Dividends
}

type BalanceSheet struct {
	Assets      Assets
	Liabilities Liabilities
}

type Assets struct {
	CurrencyUnit         bank.UnitType
	Liquid               int
	MarketableSecurities int
	AccountsReceivables  int
	Inventory            int
	PrepaidExpenses      int
	CapitalAssets        int
	IntangibleAssets     int
	Investments          int
}

type Liabilities struct {
	CurrencyUnit    bank.UnitType
	AccountsPayable int
	WagesPayable    int
	InterestPayable int
	DeferredRevenue int
	DeferredTaxes   int
	ShortTermDebts  int
	LongTermDebts   int
}

type Income struct {
	CurrencyUnit           bank.UnitType
	OperatingRevenue       int
	NonOperatingRevenue    int
	ProductionExpenses     int
	AdministrativeExpenses int
	Depreciation           int
}

type Dividends struct {
	CurrencyUnit bank.UnitType
	Payout       int
}
