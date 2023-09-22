package bank

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	"github.com/jxlxx/GreenIsland/config"
	"github.com/jxlxx/GreenIsland/requests"
	"github.com/jxlxx/GreenIsland/responses"
)

type Bank struct {
	ID             int            `yaml:"id"`
	Name           string         `yaml:"name"`
	Code           string         `yaml:"code"`
	CountryCode    string         `yaml:"country_code"`
	HomeCurrencies []CurrencyCode `yaml:"home_currencies"`

	js          nats.JetStreamContext
	service     micro.Service
	accounts    nats.KeyValue
	currencies  []Currency
	currencyMap map[CurrencyCode]Currency
}

func (b Bank) Bucket() string {
	return fmt.Sprintf("bank-%s-%d", b.CountryCode, b.ID)
}

func (b Bank) ServiceName() string {
	return fmt.Sprintf("%sBankingService", b.Code)
}

func (b Bank) ServiceGroup() string {
	return fmt.Sprintf("%s.%s", b.CountryCode, b.Code)
}

func (b Bank) AdminGroup() string {
	return fmt.Sprintf("admin.%s.%s", b.CountryCode, b.Code)
}

func (b Bank) Description() string {
	return fmt.Sprintf("This is the banking microservice for %s.", b.Name)
}

func (b Bank) Version() string {
	return config.GetEnvOrDefault("VERSION", "0.0.1")
}

type CurrencyCode string

type Currency struct {
	Name              string       `yaml:"name"`
	Code              CurrencyCode `yaml:"code"`
	MicroUnit         CurrencyUnit `yaml:"micro_unit"`
	MinorUnit         CurrencyUnit `yaml:"minor_unit"`
	MajorUnit         CurrencyUnit `yaml:"major_unit"`
	HighVolumeUnit    CurrencyUnit `yaml:"high_volume_unit"`
	MegaVolumeUnit    CurrencyUnit `yaml:"mega_volume_unit"`
	HighestVolumeUnit CurrencyUnit `yaml:"highest_volume_unit"`
	UnitMap           map[UnitType]CurrencyUnit
}

type CurrencyValue struct {
	Currency CurrencyCode `yaml:"currency"`
	Unit     UnitType     `yaml:"currency_unit"`
	Value    int          `yaml:"value"`
	Jitter   int          `yaml:"jitter"`
	Average  int          `yaml:"average_delta"`
}

func (v CurrencyValue) CalcUpdate() int {
	if v.Jitter == 0 {
		return 0
	}
	delta := rand.Intn(v.Jitter) - v.Jitter/2
	shift := v.Average
	return delta + shift
}

type UnitType string

const (
	Micro     UnitType = "micro"
	Minor     UnitType = "minor"
	Major     UnitType = "major"
	Millions  UnitType = "millions"
	Billions  UnitType = "billions"
	Trillions UnitType = "trillions"
)

type CurrencyUnit struct {
	UnitType     UnitType `yaml:"unit_type"`
	NameSingular string   `yaml:"name_singular"`
	NamePlural   string   `yaml:"name_plural"`
	Symbol       string   `yaml:"symbol"`
	MinorRatio   int      `yaml:"minor_ratio"`
}

type Account struct {
	UserID uuid.UUID
	Funds  map[CurrencyCode]Funds
}

type Funds struct {
	TotalMinor     int
	AvailableMinor int
	OnHoldMinor    int
	TotalMajor     int
	AvailableMajor int
	OnHoldMajor    int
	Currency       CurrencyCode
}

func (b *Bank) Setup() {
	currencies, _ := InitCurrencies()
	cm := make(map[CurrencyCode]Currency)
	for _, c := range currencies {
		units := make(map[UnitType]CurrencyUnit)

		units[Micro] = c.MicroUnit
		units[Minor] = c.MinorUnit
		units[Major] = c.MajorUnit
		units[Millions] = c.HighVolumeUnit
		units[Billions] = c.MegaVolumeUnit
		units[Trillions] = c.HighestVolumeUnit

		c.UnitMap = units
		cm[c.Code] = c
	}
	b.currencies = currencies
	b.currencyMap = cm
}

func (b *Bank) Update() {

}

func (b *Bank) ServiceConfig() micro.Config {
	conf := micro.Config{
		Name:        b.ServiceName(),
		Version:     b.Version(),
		Description: b.Description(),
	}
	return conf
}

func (b *Bank) Connect() {
	js := config.JetStream()
	kv, err := js.KeyValue(b.Bucket())
	if err != nil {
		log.Fatalln(err)
	}
	b.js = js
	b.accounts = kv
}

func (b *Bank) AddService(nc *nats.Conn) micro.Service {
	conf := b.ServiceConfig()
	srv, err := micro.AddService(nc, conf)
	if err != nil {
		log.Fatalln(err)
	}
	b.service = srv
	b.AddEndpoints()
	return srv
}

func respondParsingError(req micro.Request) {
	resp := responses.Response{
		Status: responses.Error,
	}
	if err := req.RespondJSON(resp); err != nil {
		fmt.Println(err)
	}
}

func (b *Bank) addUserRequest(req micro.Request) {
	r := requests.NewBankAccount{}
	if err := json.Unmarshal(req.Data(), &r); err != nil {
		fmt.Println(1, r)
		respondParsingError(req)
		return
	}
	id, err := uuid.Parse(r.ID)
	if id == uuid.Nil || err != nil {
		fmt.Println(2)
		respondParsingError(req)
		return
	}
	if err := b.AddUser(id); err != nil {
		fmt.Println(3)
		fmt.Println(err)
		respondParsingError(req)
		return
	}
	response := responses.NewBankAccount{Status: responses.OK}
	if err := req.RespondJSON(response); err != nil {
		fmt.Println(err)
	}
}

func (b *Bank) getUserFundsRequest(req micro.Request) {
	_ = req.Respond([]byte("swag 2"))
}

func (b *Bank) AddEndpoints() {
	if b.service == nil {
		log.Fatalln("err: service isn't set")
	}
	group := b.service.AddGroup(b.ServiceGroup())
	if err := group.AddEndpoint("addUser", micro.HandlerFunc(b.addUserRequest)); err != nil {
		log.Fatalln(err)
	}
	if err := group.AddEndpoint("getFunds", micro.HandlerFunc(b.getUserFundsRequest)); err != nil {
		log.Fatalln(err)
	}
	admin := b.service.AddGroup(b.AdminGroup())
	if err := admin.AddEndpoint("setFunds", micro.HandlerFunc(b.adminSetFundsRequest)); err != nil {
		log.Fatalln()
	}
}

type InitializeCompanyBankAccount struct {
	ID       uuid.UUID
	Currency CurrencyCode
	Unit     UnitType
	Sum      int
}

func (b Bank) adminSetFundsRequest(req micro.Request) {
	r := InitializeCompanyBankAccount{}
	if err := json.Unmarshal(req.Data(), &r); err != nil {
		log.Println(err)
	}
	if err := b.AddUser(r.ID); err != nil {
		log.Println(err)
	}
	minorSum, err := b.ConvertCurrency(r.Currency, r.Unit, Minor, r.Sum)
	if err != nil {
		log.Println(err)
	}
	if err := b.initialDeposit(r.ID, r.Currency, minorSum); err != nil {
		log.Println(err)
	}
	funds, err := b.GetFunds(r.ID)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(funds)
	req.RespondJSON(funds)
}

func (b Bank) initialDeposit(id uuid.UUID, currency CurrencyCode, sum int) error {
	return b.Put(id, currency, Available, sum)
}

func (b Bank) Init() {
	js := config.JetStream()
	_, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: b.Bucket(),
	})
	if err != nil {
		log.Fatalln(err)
	}
}

type Availability string

const (
	Available Availability = "available"
	OnHold    Availability = "on_hold"
)

func (b Bank) Put(userID uuid.UUID, currency CurrencyCode, status Availability, value int) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, ok := b.currencyMap[currency]; !ok {
		return fmt.Errorf("err: unknown currency: %s", currency)
	}
	key := fmt.Sprintf("%s.%s.%s", userID.String(), string(currency), string(status))
	_, err = b.accounts.Put(key, v)
	return err
}

func (b Bank) Get(userID uuid.UUID, currency CurrencyCode, status Availability) (int, error) {
	key := fmt.Sprintf("%s.%s.%s", userID.String(), string(currency), string(status))
	v, err := b.accounts.Get(key)
	if err != nil {
		return 0, err
	}
	var i int
	err = json.Unmarshal(v.Value(), &i)
	return i, err
}

func (b Bank) AddUser(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("error adding user: cannot have nil user id")
	}
	for _, c := range b.currencies {
		if err := b.Put(id, c.Code, Available, 0); err != nil {
			return err
		}
		if err := b.Put(id, c.Code, OnHold, 0); err != nil {
			return nil
		}
	}
	return nil
}

func (b Bank) Transfer(give, recv uuid.UUID, code CurrencyCode, sum int, fromOnHold bool) error {
	status := Available
	if fromOnHold {
		status = OnHold
	}
	current, err := b.Get(give, code, status)
	remainder := current - sum
	if err != nil {
		return err
	}
	if remainder < 0 {
		return fmt.Errorf("transfer failed: insufficient funds")
	}

	if err := b.Put(give, code, status, remainder); err != nil {
		return err
	}
	if err := b.Put(recv, code, Available, sum); err != nil {
		return err
	}
	return nil
}

func (b Bank) Hold(user uuid.UUID, code CurrencyCode, sum int) error {
	current, err := b.Get(user, code, Available)
	remainder := current - sum
	if err != nil {
		return err
	}
	if remainder < 0 {
		return fmt.Errorf("hold failed: insufficient available securities to put on hold")
	}

	if err := b.Put(user, code, Available, remainder); err != nil {
		return err
	}
	if err := b.Put(user, code, OnHold, sum); err != nil {
		return err
	}
	return nil
}

func (b Bank) GetFunds(id uuid.UUID) (Account, error) {
	fundMap := map[CurrencyCode]Funds{}
	for _, c := range b.currencies {
		available, err := b.Get(id, c.Code, Available)
		if err != nil {
			return Account{}, err
		}
		onHold, err := b.Get(id, c.Code, OnHold)
		if err != nil {
			return Account{}, err
		}
		fundMap[c.Code] = Funds{
			TotalMinor:     available + onHold,
			AvailableMinor: available,
			OnHoldMinor:    onHold,
			TotalMajor:     (available + onHold) / c.MajorUnit.MinorRatio,
			AvailableMajor: (available) / c.MajorUnit.MinorRatio,
			OnHoldMajor:    (onHold) / c.MajorUnit.MinorRatio,
			Currency:       c.Code,
		}
	}
	account := Account{
		UserID: id,
		Funds:  fundMap,
	}
	return account, nil
}

func (b Bank) ConvertCurrency(code CurrencyCode, from, to UnitType, sum int) (int, error) {
	currency, ok := b.currencyMap[code]
	if !ok {
		return 0, fmt.Errorf("err: unknown currency: %s", code)
	}
	minor, err := currency.ConvertToMinor(from, sum)
	if err != nil {
		return 0, err
	}
	return currency.ConvertFromMinor(to, minor)
}

func (c *Currency) ConvertToMinor(from UnitType, sum int) (int, error) {
	unit, ok := c.UnitMap[from]
	if !ok {
		return 0, fmt.Errorf("err: unknown currency unit: %s", from)
	}
	if from == Micro {
		return sum / unit.MinorRatio, nil
	}

	return sum * unit.MinorRatio, nil
}

func (c *Currency) ConvertFromMinor(to UnitType, sum int) (int, error) {
	unit, ok := c.UnitMap[to]
	if !ok {
		return 0, fmt.Errorf("err: unknown currency unit: %s", to)
	}
	if to == Micro {
		return sum * unit.MinorRatio, nil
	}

	return sum / unit.MinorRatio, nil
}
