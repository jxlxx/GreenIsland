package bank

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/jxlxx/GreenIsland/config"
	"github.com/jxlxx/GreenIsland/types"
)

type Bank struct {
	js          nats.JetStreamContext
	accounts    nats.KeyValue
	currencies  []Currency
	currencyMap map[CurrencyCode]Currency
}

type CurrencyCode string

type Currency struct {
	Name           string       `yaml:"name"`
	Code           CurrencyCode `yaml:"code"`
	MicroUnit      CurrencyUnit `yaml:"micro_unit"`
	MinorUnit      CurrencyUnit `yaml:"minor_unit"`
	MajorUnit      CurrencyUnit `yaml:"major_unit"`
	HighVolumeUnit CurrencyUnit `yaml:"high_volume_unit"`
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
	Currency       Currency
}

func New() *Bank {
	js := config.JetStream()
	kv, err := js.KeyValue("bank")
	if err != nil {
		log.Fatalln(err)
	}
	currencies, _ := InitCurrencies()
	cm := make(map[CurrencyCode]Currency)
	for _, c := range currencies {
		cm[c.Code] = c
	}
	return &Bank{
		js:          js,
		accounts:    kv,
		currencies:  currencies,
		currencyMap: cm,
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
			Currency:       c,
		}
	}
	account := Account{
		UserID: u.ID,
		Funds:  fundMap,
	}
	return account, nil
}

func Initialize() {
	js := config.JetStream()
	_, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "bank",
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func InitCurrencies() ([]Currency, error) {
	return []Currency{
		{
			Name: "canadian dollars",
			Code: "CAD",
			MajorUnit: CurrencyUnit{
				NameSingular: "dollar",
				NamePlural:   "dollars",
				Symbol:       "$",
				MinorRatio:   100,
			},
			MicroUnit: CurrencyUnit{
				NameSingular: "hundredth of a cent",
				NamePlural:   "hundredths of a cent",
				Symbol:       "µ",
				MinorRatio:   100,
			},
			HighVolumeUnit: CurrencyUnit{
				NameSingular: "Million",
				NamePlural:   "Millions",
				Symbol:       "M",
				MinorRatio:   100000000,
			},
			MinorUnit: CurrencyUnit{
				NameSingular: "penny",
				NamePlural:   "pennies",
				Symbol:       "¢",
				MinorRatio:   1,
			},
		},
		{
			Name: "american dollars",
			Code: "USD",
			MicroUnit: CurrencyUnit{
				NameSingular: "hundredth of a cent",
				NamePlural:   "hundredths of a cent",
				Symbol:       "µ",
				MinorRatio:   100,
			},
			MajorUnit: CurrencyUnit{
				NameSingular: "dollar",
				NamePlural:   "dollars",
				Symbol:       "$",
			},
			HighVolumeUnit: CurrencyUnit{
				NameSingular: "Million",
				NamePlural:   "Millions",
				Symbol:       "M",
				MinorRatio:   100000000,
			},
			MinorUnit: CurrencyUnit{
				NameSingular: "penny",
				NamePlural:   "pennies",
				Symbol:       "¢",
			},
		},
		{
			Name: "british pound sterling",
			Code: "GBP",
			MicroUnit: CurrencyUnit{
				NameSingular: "hundredth of a penny",
				NamePlural:   "hundredths of a penny",
				Symbol:       "µ",
				MinorRatio:   100,
			},
			MajorUnit: CurrencyUnit{
				NameSingular: "pound",
				NamePlural:   "pounds",
				Symbol:       "£",
			},
			HighVolumeUnit: CurrencyUnit{
				NameSingular: "Million",
				NamePlural:   "Millions",
				Symbol:       "M",
				MinorRatio:   100000000,
			},
			MinorUnit: CurrencyUnit{
				NameSingular: "penny",
				NamePlural:   "pence",
				Symbol:       "p",
			},
		},
		{
			Name: "euro",
			Code: "EUR",
			MicroUnit: CurrencyUnit{
				NameSingular: "hundredth of a euro",
				NamePlural:   "hundredths of a euro",
				Symbol:       "µ",
				MinorRatio:   100,
			},
			MajorUnit: CurrencyUnit{
				NameSingular: "euro",
				NamePlural:   "euros",
				Symbol:       "€",
			},
			HighVolumeUnit: CurrencyUnit{
				NameSingular: "Million",
				NamePlural:   "Millions",
				Symbol:       "M",
				MinorRatio:   100000000,
			},
			MinorUnit: CurrencyUnit{
				NameSingular: "euro cent",
				NamePlural:   "euro cents",
				Symbol:       "c",
			},
		},
		{
			Name: "japanese yen",
			Code: "JPY",
			MicroUnit: CurrencyUnit{
				NameSingular: "tenth of a yen",
				NamePlural:   "tenth of a yen",
				Symbol:       "µ",
				MinorRatio:   10,
			},
			MajorUnit: CurrencyUnit{
				NameSingular: "yen",
				NamePlural:   "yen",
				Symbol:       "¥",
				MinorRatio:   1,
			},
			HighVolumeUnit: CurrencyUnit{
				NameSingular: "yen",
				NamePlural:   "yen",
				Symbol:       "¥",
				MinorRatio:   1,
			},
			MinorUnit: CurrencyUnit{
				NameSingular: "yen",
				NamePlural:   "yen",
				Symbol:       "¥",
				MinorRatio:   1,
			},
		},
	}, nil
}
