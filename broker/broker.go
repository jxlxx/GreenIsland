package broker

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/jxlxx/GreenIsland/config"
)

type Broker struct {
	bucket   string
	js       nats.JetStreamContext
	accounts nats.KeyValue
}

func New(name string) *Broker {
	js := config.JetStream()
	kv, err := js.KeyValue(name)
	if err != nil {
		log.Fatalln(err)
	}
	return &Broker{
		js:       js,
		accounts: kv,
		bucket:   name,
	}
}

type Availability string

const (
	Available Availability = "available"
	OnHold    Availability = "on_hold"
)

func (b Broker) Put(userID uuid.UUID, securityID string, status Availability, value int) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s.%s.%s", userID.String(), securityID, string(status))
	_, err = b.accounts.Put(key, v)
	return err
}

func (b Broker) Get(userID uuid.UUID, securityID string, status Availability) (int, error) {
	key := fmt.Sprintf("%s.%s.%s", userID.String(), securityID, string(status))
	v, err := b.accounts.Get(key)
	if err != nil {
		return 0, err
	}
	var i int
	err = json.Unmarshal(v.Value(), &i)
	return i, err
}

func (b Broker) Transfer(give, recv uuid.UUID, securityID string, sum int, fromOnHold bool) error {
	status := Available
	if fromOnHold {
		status = OnHold
	}
	current, err := b.Get(give, securityID, status)
	remainder := current - sum
	if err != nil {
		return err
	}
	if remainder < 0 {
		return fmt.Errorf("transfer failed: insufficient funds")
	}

	if err := b.Put(give, securityID, status, remainder); err != nil {
		return err
	}
	if err := b.Put(recv, securityID, Available, sum); err != nil {
		return err
	}
	return nil
}

func (b Broker) Hold(user uuid.UUID, securityID string, sum int) error {
	current, err := b.Get(user, securityID, Available)
	remainder := current - sum
	if err != nil {
		return err
	}
	if remainder < 0 {
		return fmt.Errorf("hold failed: insufficient available securities to put on hold")
	}

	if err := b.Put(user, securityID, Available, remainder); err != nil {
		return err
	}
	if err := b.Put(user, securityID, OnHold, sum); err != nil {
		return err
	}
	return nil
}

func Initialize(name string) {
	js := config.JetStream()
	_, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: name,
	})
	if err != nil {
		log.Fatalln(err)
	}
}
