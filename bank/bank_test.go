package bank

import (
	"testing"
)

// func (b Bank) ConvertCurrency(code CurrencyCode, from, to UnitType, sum int) (int, error)
// func (c Currency) ConvertToMinor(from UnitType, sum int) (int, error) {

func TestMatchingOrder(t *testing.T) {
	b := &Bank{}
	b.Setup()
	tests := []struct {
		name     string
		code     CurrencyCode
		from     UnitType
		to       UnitType
		input    int
		expected int
	}{
		{
			name:     "simple",
			code:     b.currencyMap["USD"].Code,
			from:     Minor,
			to:       Minor,
			input:    0,
			expected: 0,
		},
		{
			name:     "100c to $1",
			code:     b.currencyMap["USD"].Code,
			from:     Minor,
			to:       Major,
			input:    100,
			expected: 1,
		},
		{
			name:     "1000000$ to 1 mil",
			code:     b.currencyMap["USD"].Code,
			from:     Major,
			to:       Millions,
			input:    1000000,
			expected: 1,
		},
		{
			name:     "Not enough",
			code:     b.currencyMap["USD"].Code,
			from:     Major,
			to:       Trillions,
			input:    1000000,
			expected: 0,
		},
		{
			name:     "major to micro",
			code:     b.currencyMap["USD"].Code,
			from:     Major,
			to:       Micro,
			input:    1,
			expected: 10000,
		},
		{
			name:     "minor to micro",
			code:     b.currencyMap["USD"].Code,
			from:     Minor,
			to:       Micro,
			input:    1,
			expected: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bank := &Bank{}
			bank.Setup()
			result, err := bank.ConvertCurrency(tt.code, tt.from, tt.to, tt.input)
			if err != nil {
				t.Fatalf(err.Error())
			}
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
