package bank

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	"github.com/jxlxx/GreenIsland/config"
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
	customers   nats.KeyValue
	currencies  []Currency
	currencyMap map[CurrencyCode]Currency
}

type AccountStatus string

const (
	Active    AccountStatus = "active"
	Cancelled AccountStatus = "cancelled"
)

type Account struct {
	UserID    uuid.UUID
	AccountID uuid.UUID
	Status    AccountStatus
	Funds     map[CurrencyCode]Funds
}

type Funds struct {
	TotalMinor     int
	AvailableMinor int
	OnHoldMinor    int
	TotalMajor     int
	AvailableMajor int
	OnHoldMajor    int
	Currency       CurrencyCode
	MajorUnit      UnitType
	MinorUnit      UnitType
}

type Availability string

const (
	Available Availability = "available"
	OnHold    Availability = "on_hold"
)

func (b *Bank) Setup() {
	currencies, _ := initCurrencies()
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

func (b *Bank) Connect() {
	js := config.JetStream()
	accounts, err := js.KeyValue(b.accountBucket())
	if err != nil {
		log.Fatalln(err)
	}
	customers, err := js.KeyValue(b.customerBucket())
	if err != nil {
		log.Fatalln(err)
	}
	b.js = js
	b.accounts = accounts
	b.customers = customers
}

func (b Bank) Init() {
	js := config.JetStream()
	_, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: b.accountBucket(),
	})
	if err != nil {
		log.Fatalln(err)
	}
	_, err = js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: b.customerBucket(),
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func (b Bank) newAccount(id uuid.UUID) (Account, error) {
	if id == uuid.Nil {
		return Account{}, fmt.Errorf("error adding user: cannot have nil user id")
	}
	account := Account{}
	for _, c := range b.currencies {
		if err := b.put(id, c.Code, Available, 0); err != nil {
			return Account{}, err
		}
		if err := b.put(id, c.Code, OnHold, 0); err != nil {
			return Account{}, nil
		}
	}
	return account, nil
}

func (b Bank) getAccount(id uuid.UUID) (Account, error) {
	fundMap := map[CurrencyCode]Funds{}
	for _, c := range b.currencies {
		available, err := b.get(id, c.Code, Available)
		if err != nil {
			return Account{}, err
		}
		onHold, err := b.get(id, c.Code, OnHold)
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
		AccountID: id,
		Funds:     fundMap,
	}
	return account, nil
}

func (b Bank) transfer(give, recv uuid.UUID, code CurrencyCode, sum int, fromOnHold bool) error {
	status := Available
	if fromOnHold {
		status = OnHold
	}
	current, err := b.get(give, code, status)
	remainder := current - sum
	if err != nil {
		return err
	}
	if remainder < 0 {
		return fmt.Errorf("transfer failed: insufficient funds")
	}

	if err := b.put(give, code, status, remainder); err != nil {
		return err
	}
	if err := b.put(recv, code, Available, sum); err != nil {
		return err
	}
	return nil
}

func (b Bank) hold(user uuid.UUID, code CurrencyCode, sum int) error {
	current, err := b.get(user, code, Available)
	remainder := current - sum
	if err != nil {
		return err
	}
	if remainder < 0 {
		return fmt.Errorf("hold failed: insufficient available securities to put on hold")
	}

	if err := b.put(user, code, Available, remainder); err != nil {
		return err
	}
	if err := b.put(user, code, OnHold, sum); err != nil {
		return err
	}
	return nil
}
