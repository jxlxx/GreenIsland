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
	"github.com/jxlxx/GreenIsland/types"
)

type Bank struct {
	ID             int            `yaml:"id"`
	Name           string         `yaml:"name"`
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
	return fmt.Sprintf("service-bank-%s-%d", b.CountryCode, b.ID)
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

func (b *Bank) Handle(req micro.Request) {
	_ = req.Respond([]byte("swag"))
}

func (b *Bank) AddEndpoints() {
	if b.service == nil {
		log.Fatalln("err: service isn't set")
	}
	_ = b.service.AddEndpoint("echo", b)
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

func (b Bank) AddUser(u types.User) error {
	if u.ID == uuid.Nil {
		return fmt.Errorf("error adding user: cannot have nil user id")
	}
	for _, c := range b.currencies {
		if err := b.Put(u.ID, c.Code, Available, 0); err != nil {
			return err
		}
		if err := b.Put(u.ID, c.Code, OnHold, 0); err != nil {
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

func (b Bank) GetUserFunds(u types.User) (Account, error) {
	fundMap := map[CurrencyCode]Funds{}
	for _, c := range b.currencies {
		available, err := b.Get(u.ID, c.Code, Available)
		if err != nil {
			return Account{}, err
		}
		onHold, err := b.Get(u.ID, c.Code, OnHold)
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
		UserID: u.ID,
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
