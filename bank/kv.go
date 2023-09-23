package bank

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func (b Bank) Put(id uuid.UUID, currency CurrencyCode, status Availability, value int) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, ok := b.currencyMap[currency]; !ok {
		return fmt.Errorf("err: unknown currency: %s", currency)
	}
	key := fmt.Sprintf("%s.%s.%s", id.String(), string(currency), string(status))
	_, err = b.accounts.Put(key, v)
	return err
}

func (b Bank) Get(id uuid.UUID, currency CurrencyCode, status Availability) (int, error) {
	key := fmt.Sprintf("%s.%s.%s", id.String(), string(currency), string(status))
	v, err := b.accounts.Get(key)
	if err != nil {
		return 0, err
	}
	var i int
	err = json.Unmarshal(v.Value(), &i)
	return i, err
}

func (b Bank) AddCustomerAccount(id uuid.UUID, accountID uuid.UUID) error {
	key := fmt.Sprintf("%s.%s", id.String(), accountID.String())
	v, err := json.Marshal(Active)
	if err != nil {
		return err
	}
	_, err = b.accounts.Put(key, v)
	return err
}

func (b Bank) CancelCustomerAccount(id uuid.UUID, accountID uuid.UUID) error {
	key := fmt.Sprintf("%s.%s", id.String(), accountID.String())
	v, err := json.Marshal(Cancelled)
	if err != nil {
		return err
	}
	_, err = b.accounts.Put(key, v)
	return err
}
