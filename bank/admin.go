package bank

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/micro"
)

type Deposit struct {
	AccountID uuid.UUID
	Currency  CurrencyCode
	Unit      UnitType
	Sum       int
}

func (b Bank) adminDepositRequest(req micro.Request) {
	r := Deposit{}
	if err := json.Unmarshal(req.Data(), &r); err != nil {
		log.Println(err)
	}
	minorSum, err := b.ConvertCurrency(r.Currency, r.Unit, Minor, r.Sum)
	if err != nil {
		log.Println(err)
	}
	if err := b.Put(r.AccountID, r.Currency, Available, minorSum); err != nil {
		log.Println(err)
	}
	account, err := b.getAccount(r.AccountID)
	if err != nil {
		log.Println(err)
	}
	if err := req.RespondJSON(account); err != nil {
		log.Println(err)
	}
}

func (b Bank) adminTransferRequest(req micro.Request) {
	_ = req.Respond([]byte("unimplemented"))
}

func (b Bank) adminHoldRequest(req micro.Request) {
	_ = req.Respond([]byte("unimplemented"))
}
