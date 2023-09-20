package company

import (
	"github.com/jxlxx/GreenIsland/bank"
	"github.com/jxlxx/GreenIsland/world"
)

type Company struct {
	FullName      string `yaml:"full_name"`
	Name          string `yaml:"name"`
	HQCountryCode string `yaml:"hq_country_code"`

	StockSymbol       string            `yaml:"stock_symbol"`
	DefaultCurrency   bank.CurrencyCode `yaml:"currency_code"`
	OutstandingShares int               `yaml:"outstanding_shares"`

	// Balance Sheet
	BalanceSheet            BalanceSheet `yaml:"balance_sheet"`
	DailyBalanceSheetJitter BalanceSheet `yaml:"daily_balance_sheet_jitter"`

	// Revenue + expenses
	Income            Income `yaml:"income"`
	DailyIncomeJitter Income `yaml:"daily_income_jitter"`

	QuarterlyBehaviour QuarterlyBehaviour `yaml:"quarterly_behaviour"`
	QuarterlyMetrics   QuarterlyMetrics   `yaml:"quarterly_metrics"`

	// world context
	Employment Employment `yaml:"employment"`
	Industries Industries `yaml:"industries"`
}

type BalanceSheet struct {
	UnitType    bank.UnitType `yaml:"unit_type"`
	Assets      Assets        `yaml:"assets"`
	Liabilities Liabilities   `yaml:"liabilities"`
}

type Employment struct {
	UnitType             bank.UnitType `yaml:"unit_type"`
	Employees            int           `yaml:"employees"`
	EmployeeSatisfaction int           `yaml:"employee_satisfaction"`
	DailyTurnover        int           `yaml:"daily_turnover"`
	DailyTurnoverJitter  int           `yaml:"daily_turnover_jitter"`
	HighestAnnualSalary  int           `yaml:"highest_annual_salary"`
	AverageAnnualSalary  int           `yaml:"average_annual_salary"`
	LowestAnnualSalary   int           `yaml:"lowest_annual_salary"`
}

type Industries struct {
	PrimaryIndustries   []world.Industry `yaml:"primary_industries"`
	SecondaryIndustries []world.Industry `yaml:"secondary_industries"`
}

type Assets struct {
	LiquidAssets         int `yaml:"liquid_assets"`
	MarketableSecurities int `yaml:"marketable_securities"`
	AccountsReceivables  int `yaml:"accounts_receivables"`
	Inventory            int `yaml:"inventory"`
	PrepaidExpenses      int `yaml:"prepaid_expenses"`
	CapitalAssets        int `yaml:"capital_assets"`
	IntangibleAssets     int `yaml:"intangible_assets"`
	Investments          int `yaml:"investments"`
}

type Liabilities struct {
	AccountsPayable int `yaml:"accounts_payable"`
	WagesPayable    int `yaml:"wages_payable"`
	InterestPayable int `yaml:"interest_payable"`
	DeferredRevenue int `yaml:"deferred_revenue"`
	DeferredTaxes   int `yaml:"deferred_taxes"`
	ShortTermDebts  int `yaml:"short_term_debts"`
	LongTermDebts   int `yaml:"long_term_debts"`
}

type Income struct {
	UnitType               bank.UnitType `yaml:"unit_type"`
	OperatingRevenue       int           `yaml:"operating_revenue"`
	NonOperatingRevenue    int           `yaml:"non_operating_revenue"`
	ProductionExpenses     int           `yaml:"production_expenses"`
	AdministrativeExpenses int           `yaml:"administrative_expenses"`
	Depreciation           int           `yaml:"depreciation"`
}

type Delta struct {
	Assets      Assets      `yaml:"assets_delta"`
	Liabilities Liabilities `yaml:"liabilities_delta"`
	Income      Income      `yaml:"income_delta"`
}

type QuarterlyBehaviour struct {
	UnitType       bank.UnitType `yaml:"unit_type"`
	DividendPayout int           `yaml:"dividend_payout"`
	ShareBuyback   int           `yaml:"share_buyback"`
}

type QuarterlyMetrics struct {
	UnitType           bank.UnitType `yaml:"unit_type"`
	DividendGrowthRate int           `yaml:"dividend_growth_rate"`
	EquityGrowthRate   int           `yaml:"equity_growth_rate"`
	CurrentStockPrice  int           `yaml:"current_stock_price"`
	ProjectedDividends int           `yaml:"projected_dividends"`
}

type CompanyCycle string

const (
	Peak      CompanyCycle = "peak"
	Recession CompanyCycle = "recession"
	Trough    CompanyCycle = "trough"
	Recovery  CompanyCycle = "recovery"
	Expansion CompanyCycle = "expansion"
)
