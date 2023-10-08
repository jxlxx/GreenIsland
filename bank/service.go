package bank

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	"github.com/jxlxx/GreenIsland/config"
)

///////////////////////////////////////////////////////////////////////////////

type Handler interface {
	CreateAccount(micro.Request, uuid.UUID)
	GetAccountByID(micro.Request, uuid.UUID)
	GetAccountsByOwnerID(micro.Request, uuid.UUID)
	AdminDeposit(micro.Request, Deposit)
	AdminTransfer(micro.Request, Transfer)
	AdminHold(micro.Request, Hold)
}

type ServiceWrapper struct {
	Handler Handler
}

type Options struct {
	Name        string
	Version     string
	Description string
	CountryCode string
	BankCode    string
}

func CreateService(nc *nats.Conn, h Handler, opts Options) (micro.Service, error) {
	conf := micro.Config{
		Name:        opts.Name,
		Version:     opts.Version,
		Description: opts.Description,
	}
	service, err := micro.AddService(nc, conf)
	if err != nil {
		return nil, err
	}
	s := ServiceWrapper{
		Handler: h,
	}

	base := service.AddGroup(fmt.Sprintf("bank.%s.%s", opts.CountryCode, opts.BankCode))
	admin := service.AddGroup(fmt.Sprintf("admin.bank.%s.%s", opts.CountryCode, opts.BankCode))

	if err := base.AddEndpoint("create", micro.HandlerFunc(s.CreateAccount)); err != nil {
		return nil, err
	}
	if err := base.AddEndpoint("account", micro.HandlerFunc(s.GetAccountByID)); err != nil {
		return nil, err
	}
	if err := base.AddEndpoint("accounts", micro.HandlerFunc(s.GetAccountsByOwnerID)); err != nil {
		return nil, err
	}

	if err := admin.AddEndpoint("deposit", micro.HandlerFunc(s.AdminDeposit)); err != nil {
		return nil, err
	}
	if err := admin.AddEndpoint("transfer", micro.HandlerFunc(s.AdminTransfer)); err != nil {
		return nil, err
	}
	if err := admin.AddEndpoint("hold", micro.HandlerFunc(s.AdminHold)); err != nil {
		return nil, err
	}
	return service, nil
}

func (s *ServiceWrapper) CreateAccount(r micro.Request) {
}

func (s *ServiceWrapper) GetAccountByID(r micro.Request) {
}

func (s *ServiceWrapper) GetAccountsByOwnerID(r micro.Request) {
}

func (s ServiceWrapper) AdminDeposit(req micro.Request) {
	subj := strings.Split(req.Subject(), ".")
	if len(subj) != 2 {
		respondError(req, "wrong number of tokens in subject")
	}
	deposit := Deposit{}
	if err := json.Unmarshal(req.Data(), &deposit); err != nil {
		log.Println(err)
	}
	s.Handler.AdminDeposit(req, deposit)
}

func (s *ServiceWrapper) AdminTransfer(r micro.Request) {
}

func (s *ServiceWrapper) AdminHold(r micro.Request) {
}

func (b *Bank) serviceConfig() micro.Config {
	conf := micro.Config{
		Name:        b.serviceName(),
		Version:     b.version(),
		Description: b.description(),
	}
	return conf
}

///////////////////////////////////////////////////////////////////////////////

func (b *Bank) AddService(nc *nats.Conn) micro.Service {
	srv, err := CreateService(nc, b, Options{})
	if err != nil {
		log.Fatalln(err)
	}
	return srv
}

func respondError(req micro.Request, errorMessage string) {
	resp := Response{
		Status:  "Error",
		Message: errorMessage,
	}
	if err := req.RespondJSON(resp); err != nil {
		fmt.Println(err)
	}
}

func (b *Bank) CreateAccount(req micro.Request, id uuid.UUID) {
	r := NewAccountPayload{}
	if err := json.Unmarshal(req.Data(), &r); err != nil {
		respondError(req, "cannot parse request")
		return
	}
	account, err := b.newAccount(r.UserID)
	if err != nil {
		fmt.Println(err)
		respondError(req, err.Error())
		return
	}
	response := AccountResponse{Status: "OK", Account: account}
	if err := req.RespondJSON(response); err != nil {
		fmt.Println(err)
	}
}

func (b Bank) GetAccountByID(req micro.Request, id uuid.UUID) {
	_ = req.Respond([]byte("unimplemented"))
}

func (b *Bank) GetAccountsByOwnerID(req micro.Request, id uuid.UUID) {
	_ = req.Respond([]byte("unimplemented"))
}

type Deposit struct {
	AccountID uuid.UUID
	Currency  CurrencyCode
	Unit      UnitType
	Sum       int
}

func (b Bank) AdminDeposit(req micro.Request, deposit Deposit) {
	minorSum, err := b.ConvertCurrency(deposit.Currency, deposit.Unit, Minor, deposit.Sum)
	if err != nil {
		log.Println(err)
	}
	if err := b.put(deposit.AccountID, deposit.Currency, Available, minorSum); err != nil {
		log.Println(err)
	}
	account, err := b.getAccount(deposit.AccountID)
	if err != nil {
		log.Println(err)
	}
	if err := req.RespondJSON(account); err != nil {
		log.Println(err)
	}
}

type Transfer struct {
}

func (b Bank) AdminTransfer(req micro.Request, t Transfer) {
	_ = req.Respond([]byte("unimplemented"))
}

type Hold struct {
}

func (b Bank) AdminHold(req micro.Request, h Hold) {
	_ = req.Respond([]byte("unimplemented"))
}

func (b Bank) accountBucket() string {
	return fmt.Sprintf("bank-accounts-%s-%s-%d", b.CountryCode, b.Code, b.ID)
}

func (b Bank) customerBucket() string {
	return fmt.Sprintf("bank-customers-%s-%s-%d", b.CountryCode, b.Code, b.ID)
}

func (b Bank) serviceName() string {
	return fmt.Sprintf("%sBankingService", b.Code)
}

func (b Bank) serviceGroup() string {
	return fmt.Sprintf("%s.%s", b.CountryCode, b.Code)
}

func (b Bank) adminGroup() string {
	return fmt.Sprintf("admin.%s.%s", b.CountryCode, b.Code)
}

func (b Bank) description() string {
	return fmt.Sprintf("This is the banking microservice for %s.", b.Name)
}

func (b Bank) version() string {
	return config.GetEnvOrDefault("VERSION", "0.0.1")
}
