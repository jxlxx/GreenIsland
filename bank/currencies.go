package bank

func InitCurrencies() ([]Currency, error) {
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