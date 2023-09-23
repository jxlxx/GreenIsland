package world

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/jxlxx/GreenIsland/bank"
	"github.com/jxlxx/GreenIsland/config"
	"github.com/jxlxx/GreenIsland/payloads"
	"github.com/jxlxx/GreenIsland/subjects"
	"github.com/jxlxx/GreenIsland/types"
)

type Company struct {
	FullName        string            `yaml:"full_name"`
	Name            string            `yaml:"name"`
	HQCountryCode   string            `yaml:"hq_country_code"`
	Code            string            `yaml:"code"`
	BankCode        string            `yaml:"bank_code"`
	DefaultCurrency bank.CurrencyCode `yaml:"currency_code"`

	OutstandingShares int          `yaml:"outstanding_shares"`
	BalanceSheet      BalanceSheet `yaml:"balance_sheet"`
	Income            Income       `yaml:"income"`

	Bid bank.CurrencyValue `yaml:"bid"`
	Ask bank.CurrencyValue `yaml:"ask"`

	QuarterlyBehaviour QuarterlyBehaviour `yaml:"quarterly_behaviour"`
	QuarterlyMetrics   QuarterlyMetrics   `yaml:"quarterly_metrics"`

	Employment Employment `yaml:"employment"`
	Industries Industries `yaml:"industries"`

	nc *nats.EncodedConn
	id uuid.UUID
}

func (c *Company) InitializeCompany() {
	nc := config.Connect()
	defer func() {
		if err := nc.Drain(); err != nil {
			fmt.Println(err)
		}
	}()
	c.id = uuid.New()
	req := bank.NewAccountRequest{
		UserID: c.id,
	}
	resp, err := nc.Request("%s.%s.create", payloads.Bytes(req), time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	account := bank.Account{}
	if err := json.Unmarshal(resp.Data, &account); err != nil {
		log.Fatalln(err)
	}
	r := bank.Deposit{
		AccountID: account.AccountID,
		Currency:  c.BalanceSheet.Assets.LiquidAssets.Currency,
		Unit:      c.BalanceSheet.Assets.LiquidAssets.Unit,
		Sum:       c.BalanceSheet.Assets.LiquidAssets.Value,
	}
	subject := fmt.Sprintf("admin.%s.%s.deposit", c.HQCountryCode, c.BankCode)
	fmt.Println(subject)
	if _, err := nc.Request(subject, payloads.Bytes(r), time.Second); err != nil {
		log.Fatalln(err)
	}
}

func (c *Company) DailySubscriber() func(payloads.WorldTick) {
	return func(payloads.WorldTick) {
		c.DailyUpdate()
	}
}

func (c *Company) QuarterlySubscriber() func(payloads.WorldTick) {
	return func(p payloads.WorldTick) {
		update := payloads.QuarterlyCompanyUpdate{
			Name:         c.Name,
			Quarter:      p.Quarter,
			CurrencyCode: c.DefaultCurrency,
			Employees:    c.Employment.Employees.Value,
			BalanceSheet: c.CreateBalanceSheet(),
			IncomeSheet:  c.CreateIncome(),
			Dividends:    c.CreateDividends(),
		}

		if err := c.nc.Publish(subjects.QuarterlyCompanyUpdate(c.Code, p.Quarter), update); err != nil {
			fmt.Println(err)
		}
	}
}

func (c Company) CreateBalanceSheet() payloads.BalanceSheet {
	return payloads.BalanceSheet{
		Assets:      c.CreateAssets(),
		Liabilities: c.CreateLiabilities(),
	}
}

func (c Company) CreateAssets() payloads.Assets {
	a := c.BalanceSheet.Assets
	return payloads.Assets{
		CurrencyUnit:         a.LiquidAssets.Unit,
		Liquid:               a.LiquidAssets.Value,
		MarketableSecurities: a.MarketableSecurities.Value,
		AccountsReceivables:  a.AccountsReceivables.Value,
		Inventory:            a.Inventory.Value,
		PrepaidExpenses:      a.PrepaidExpenses.Value,
		CapitalAssets:        a.CapitalAssets.Value,
		IntangibleAssets:     a.IntangibleAssets.Value,
		Investments:          a.Investments.Value,
	}
}
func (c Company) CreateLiabilities() payloads.Liabilities {
	l := c.BalanceSheet.Liabilities
	return payloads.Liabilities{
		CurrencyUnit:    l.AccountsPayable.Unit,
		AccountsPayable: l.AccountsPayable.Value,
		WagesPayable:    l.WagesPayable.Value,
		InterestPayable: l.InterestPayable.Value,
		DeferredRevenue: l.DeferredRevenue.Value,
		DeferredTaxes:   l.DeferredTaxes.Value,
		ShortTermDebts:  l.ShortTermDebts.Value,
		LongTermDebts:   l.LongTermDebts.Value,
	}
}
func (c Company) CreateIncome() payloads.Income {
	i := c.Income
	return payloads.Income{
		CurrencyUnit:           i.OperatingRevenue.Unit,
		OperatingRevenue:       i.OperatingRevenue.Value,
		NonOperatingRevenue:    i.NonOperatingRevenue.Value,
		ProductionExpenses:     i.ProductionExpenses.Value,
		AdministrativeExpenses: i.AdministrativeExpenses.Value,
		Depreciation:           i.Depreciation.Value,
	}
}

func (c Company) CreateDividends() payloads.Dividends {
	return payloads.Dividends{
		CurrencyUnit: c.QuarterlyBehaviour.DividendPayout.Unit,
		Payout:       c.QuarterlyBehaviour.DividendPayout.Value,
	}
}

func (c *Company) DailyUpdate() {
	c.BalanceSheet = c.BalanceSheet.Update()
	c.Income = c.Income.Update()
	c.QuarterlyBehaviour = c.QuarterlyBehaviour.Update()
	c.QuarterlyMetrics = c.QuarterlyMetrics.Update()
	c.Employment = c.Employment.Update()
	c.Bid, c.Ask = c.UpdateBidAsk()
}

func (c Company) UpdateBidAsk() (bank.CurrencyValue, bank.CurrencyValue) {
	return c.Bid, c.Ask
}

type BalanceSheet struct {
	Assets      Assets      `yaml:"assets"`
	Liabilities Liabilities `yaml:"liabilities"`
}

func (b BalanceSheet) Update() BalanceSheet {
	b.Assets = b.Assets.Update()
	b.Liabilities = b.Liabilities.Update()
	return b
}

type Employment struct {
	Employees            types.Value        `yaml:"employees"`
	EmployeeSatisfaction types.Value        `yaml:"employee_satisfaction"`
	DailyTurnover        types.Value        `yaml:"daily_turnover"`
	HighestAnnualSalary  bank.CurrencyValue `yaml:"highest_annual_salary"`
	AverageAnnualSalary  bank.CurrencyValue `yaml:"average_annual_salary"`
	LowestAnnualSalary   bank.CurrencyValue `yaml:"lowest_annual_salary"`
}

func (e Employment) Update() Employment {
	e.Employees.Value += e.Employees.CalcUpdate()
	e.EmployeeSatisfaction.Value += e.EmployeeSatisfaction.CalcUpdate()
	e.DailyTurnover.Value += e.DailyTurnover.CalcUpdate()
	e.HighestAnnualSalary.Value += e.HighestAnnualSalary.CalcUpdate()
	e.AverageAnnualSalary.Value += e.AverageAnnualSalary.CalcUpdate()
	e.LowestAnnualSalary.Value += e.LowestAnnualSalary.CalcUpdate()
	return e
}

type Assets struct {
	LiquidAssets         bank.CurrencyValue `yaml:"liquid_assets"`
	MarketableSecurities bank.CurrencyValue `yaml:"marketable_securities"`
	AccountsReceivables  bank.CurrencyValue `yaml:"accounts_receivables"`
	Inventory            bank.CurrencyValue `yaml:"inventory"`
	PrepaidExpenses      bank.CurrencyValue `yaml:"prepaid_expenses"`
	CapitalAssets        bank.CurrencyValue `yaml:"capital_assets"`
	IntangibleAssets     bank.CurrencyValue `yaml:"intangible_assets"`
	Investments          bank.CurrencyValue `yaml:"investments"`
}

func (a Assets) Update() Assets {
	a.LiquidAssets.Value += a.LiquidAssets.CalcUpdate()
	a.MarketableSecurities.Value += a.MarketableSecurities.CalcUpdate()
	a.AccountsReceivables.Value += a.AccountsReceivables.CalcUpdate()
	a.Inventory.Value += a.Inventory.CalcUpdate()
	a.PrepaidExpenses.Value += a.PrepaidExpenses.CalcUpdate()
	a.CapitalAssets.Value += a.CapitalAssets.CalcUpdate()
	a.IntangibleAssets.Value += a.IntangibleAssets.CalcUpdate()
	a.Investments.Value += a.Investments.CalcUpdate()
	return a
}

type Liabilities struct {
	AccountsPayable bank.CurrencyValue `yaml:"accounts_payable"`
	WagesPayable    bank.CurrencyValue `yaml:"wages_payable"`
	InterestPayable bank.CurrencyValue `yaml:"interest_payable"`
	DeferredRevenue bank.CurrencyValue `yaml:"deferred_revenue"`
	DeferredTaxes   bank.CurrencyValue `yaml:"deferred_taxes"`
	ShortTermDebts  bank.CurrencyValue `yaml:"short_term_debts"`
	LongTermDebts   bank.CurrencyValue `yaml:"long_term_debts"`
}

func (l Liabilities) Update() Liabilities {
	l.AccountsPayable.Value += l.AccountsPayable.CalcUpdate()
	l.WagesPayable.Value += l.WagesPayable.CalcUpdate()
	l.InterestPayable.Value += l.InterestPayable.CalcUpdate()
	l.DeferredRevenue.Value += l.DeferredRevenue.CalcUpdate()
	l.DeferredTaxes.Value += l.DeferredTaxes.CalcUpdate()
	l.ShortTermDebts.Value += l.ShortTermDebts.CalcUpdate()
	l.LongTermDebts.Value += l.LongTermDebts.CalcUpdate()
	return l
}

type Income struct {
	OperatingRevenue       bank.CurrencyValue `yaml:"operating_revenue"`
	NonOperatingRevenue    bank.CurrencyValue `yaml:"non_operating_revenue"`
	ProductionExpenses     bank.CurrencyValue `yaml:"production_expenses"`
	AdministrativeExpenses bank.CurrencyValue `yaml:"administrative_expenses"`
	Depreciation           bank.CurrencyValue `yaml:"depreciation"`
}

func (i Income) Update() Income {
	i.OperatingRevenue.Value += i.OperatingRevenue.CalcUpdate()
	i.NonOperatingRevenue.Value += i.NonOperatingRevenue.CalcUpdate()
	i.ProductionExpenses.Value += i.ProductionExpenses.CalcUpdate()
	i.AdministrativeExpenses.Value += i.AdministrativeExpenses.CalcUpdate()
	i.Depreciation.Value += i.Depreciation.CalcUpdate()
	return i
}

type QuarterlyBehaviour struct {
	DividendPayout bank.CurrencyValue `yaml:"dividend_payout"`
	ShareBuyback   types.Value        `yaml:"share_buyback"`
}

func (q QuarterlyBehaviour) Update() QuarterlyBehaviour {
	q.DividendPayout.Value += q.DividendPayout.CalcUpdate()
	if q.DividendPayout.Value < 0 {
		q.DividendPayout.Value = 0
	}
	q.ShareBuyback.Value += q.ShareBuyback.CalcUpdate()
	return q
}

type QuarterlyMetrics struct {
	DividendGrowthRate   bank.CurrencyValue `yaml:"dividend_growth_rate"`
	RequiredRateOfReturn types.Value        `yaml:"required_rate_of_return"`
	CurrentStockPrice    bank.CurrencyValue `yaml:"current_stock_price"`
	ProjectedDividends   bank.CurrencyValue `yaml:"projected_dividends"`
}

func (q QuarterlyMetrics) Update() QuarterlyMetrics {
	q.DividendGrowthRate.Value += q.DividendGrowthRate.CalcUpdate()
	q.RequiredRateOfReturn.Value += q.RequiredRateOfReturn.CalcUpdate()
	q.CurrentStockPrice.Value += q.CurrentStockPrice.CalcUpdate()
	q.ProjectedDividends.Value += q.ProjectedDividends.CalcUpdate()
	return q
}

type Industries struct {
	PrimaryIndustries   []Industry `yaml:"primary_industries"`
	SecondaryIndustries []Industry `yaml:"secondary_industries"`
}

type CompanyCycle string

const (
	Peak      CompanyCycle = "peak"
	Recession CompanyCycle = "recession"
	Trough    CompanyCycle = "trough"
	Recovery  CompanyCycle = "recovery"
	Expansion CompanyCycle = "expansion"
)
