package bank

import (
	"fmt"
	"math/rand"
)

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

func initCurrencies() ([]Currency, error) {
	return []Currency{
		{
			Name: "canadian dollars",
			Code: "CAD",
			MicroUnit: CurrencyUnit{
				UnitType:     Micro,
				NameSingular: "hundredth of a cent",
				NamePlural:   "hundredths of a cent",
				Symbol:       "µ",
				MinorRatio:   100,
			},
			MinorUnit: CurrencyUnit{
				UnitType:     Minor,
				NameSingular: "penny",
				NamePlural:   "pennies",
				Symbol:       "¢",
				MinorRatio:   1,
			},
			MajorUnit: CurrencyUnit{
				UnitType:     Major,
				NameSingular: "dollar",
				NamePlural:   "dollars",
				Symbol:       "$",
				MinorRatio:   100,
			},
			HighVolumeUnit: CurrencyUnit{
				UnitType:     Millions,
				NameSingular: "Million",
				NamePlural:   "Millions",
				Symbol:       "M",
				MinorRatio:   100000000,
			},
			MegaVolumeUnit: CurrencyUnit{
				UnitType:     Billions,
				NameSingular: "Billion",
				NamePlural:   "Billions",
				Symbol:       "B",
				MinorRatio:   100000000000,
			},
			HighestVolumeUnit: CurrencyUnit{
				UnitType:     Trillions,
				NameSingular: "Trillion",
				NamePlural:   "Trillions",
				Symbol:       "Tr",
				MinorRatio:   100000000000000,
			},
		},
		{
			Name: "american dollars",
			Code: "USD",
			MicroUnit: CurrencyUnit{
				UnitType:     Micro,
				NameSingular: "hundredth of a cent",
				NamePlural:   "hundredths of a cent",
				Symbol:       "µ",
				MinorRatio:   100,
			},
			MinorUnit: CurrencyUnit{
				UnitType:     Minor,
				NameSingular: "penny",
				NamePlural:   "pennies",
				Symbol:       "¢",
				MinorRatio:   1,
			},
			MajorUnit: CurrencyUnit{
				UnitType:     Major,
				NameSingular: "dollar",
				NamePlural:   "dollars",
				Symbol:       "$",
				MinorRatio:   100,
			},
			HighVolumeUnit: CurrencyUnit{
				UnitType:     Millions,
				NameSingular: "million",
				NamePlural:   "millions",
				Symbol:       "M",
				MinorRatio:   100000000,
			},
			MegaVolumeUnit: CurrencyUnit{
				UnitType:     Billions,
				NameSingular: "billion",
				NamePlural:   "billions",
				Symbol:       "B",
				MinorRatio:   100000000000,
			},
			HighestVolumeUnit: CurrencyUnit{
				UnitType:     Trillions,
				NameSingular: "trillion",
				NamePlural:   "trillions",
				Symbol:       "Tr",
				MinorRatio:   100000000000000,
			},
		},
		{
			Name: "british pound sterling",
			Code: "GBP",
			MicroUnit: CurrencyUnit{
				UnitType:     Micro,
				NameSingular: "hundredth of a penny",
				NamePlural:   "hundredths of a penny",
				Symbol:       "µ",
				MinorRatio:   100,
			},
			MinorUnit: CurrencyUnit{
				UnitType:     Minor,
				NameSingular: "penny",
				NamePlural:   "pence",
				Symbol:       "p",
				MinorRatio:   1,
			},
			MajorUnit: CurrencyUnit{
				UnitType:     Major,
				NameSingular: "pound",
				NamePlural:   "pounds",
				Symbol:       "£",
				MinorRatio:   100,
			},
			HighVolumeUnit: CurrencyUnit{
				UnitType:     Millions,
				NameSingular: "million",
				NamePlural:   "millions",
				Symbol:       "M",
				MinorRatio:   100000000,
			},
			MegaVolumeUnit: CurrencyUnit{
				UnitType:     Billions,
				NameSingular: "billion",
				NamePlural:   "billions",
				Symbol:       "B",
				MinorRatio:   100000000000,
			},
			HighestVolumeUnit: CurrencyUnit{
				UnitType:     Trillions,
				NameSingular: "trillion",
				NamePlural:   "trillions",
				Symbol:       "tr",
				MinorRatio:   100000000000000,
			},
		},
		{
			Name: "euro",
			Code: "EUR",
			MicroUnit: CurrencyUnit{
				UnitType:     Micro,
				NameSingular: "hundredth of a euro cent",
				NamePlural:   "hundredths of a euro cent",
				Symbol:       "µ",
				MinorRatio:   100,
			},
			MinorUnit: CurrencyUnit{
				UnitType:     Minor,
				NameSingular: "euro penny",
				NamePlural:   "euro pennies",
				Symbol:       "c",
				MinorRatio:   1,
			},
			MajorUnit: CurrencyUnit{
				UnitType:     Major,
				NameSingular: "euro",
				NamePlural:   "euros",
				Symbol:       "€",
				MinorRatio:   100,
			},
			HighVolumeUnit: CurrencyUnit{
				UnitType:     Millions,
				NameSingular: "million",
				NamePlural:   "millions",
				Symbol:       "M",
				MinorRatio:   100000000,
			},
			MegaVolumeUnit: CurrencyUnit{
				UnitType:     Billions,
				NameSingular: "billion",
				NamePlural:   "billions",
				Symbol:       "B",
				MinorRatio:   100000000000,
			},
			HighestVolumeUnit: CurrencyUnit{
				UnitType:     Trillions,
				NameSingular: "trillion",
				NamePlural:   "trillions",
				Symbol:       "Tr",
				MinorRatio:   100000000000000,
			},
		},
		{
			Name: "japanese yen",
			Code: "JPY",
			MicroUnit: CurrencyUnit{
				UnitType:     Micro,
				NameSingular: "tenth of a yen",
				NamePlural:   "tenth of a yen",
				Symbol:       "µ",
				MinorRatio:   10,
			},
			MinorUnit: CurrencyUnit{
				UnitType:     Minor,
				NameSingular: "yen",
				NamePlural:   "yen",
				Symbol:       "¥",
				MinorRatio:   1,
			},
			MajorUnit: CurrencyUnit{
				UnitType:     Major,
				NameSingular: "thousand yen",
				NamePlural:   "thousand yen",
				Symbol:       "¥K",
				MinorRatio:   100000,
			},
			HighVolumeUnit: CurrencyUnit{
				UnitType:     Millions,
				NameSingular: "million",
				NamePlural:   "millions",
				Symbol:       "M",
				MinorRatio:   100000000,
			},
			MegaVolumeUnit: CurrencyUnit{
				UnitType:     Billions,
				NameSingular: "billion",
				NamePlural:   "billions",
				Symbol:       "B",
				MinorRatio:   100000000000,
			},
			HighestVolumeUnit: CurrencyUnit{
				UnitType:     Trillions,
				NameSingular: "trillion",
				NamePlural:   "trillions",
				Symbol:       "Tr",
				MinorRatio:   100000000000000,
			},
		},
	}, nil
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

func (c Currency) ConvertToMinor(from UnitType, sum int) (int, error) {
	unit, ok := c.UnitMap[from]
	if !ok {
		return 0, fmt.Errorf("err: unknown currency unit: %s", from)
	}
	if from == Micro {
		return sum / unit.MinorRatio, nil
	}

	return sum * unit.MinorRatio, nil
}

func (c Currency) ConvertFromMinor(to UnitType, sum int) (int, error) {
	unit, ok := c.UnitMap[to]
	if !ok {
		return 0, fmt.Errorf("err: unknown currency unit: %s", to)
	}
	if to == Micro {
		return sum * unit.MinorRatio, nil
	}

	return sum / unit.MinorRatio, nil
}
