package bank

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	"github.com/jxlxx/GreenIsland/config"
)

func (b *Bank) serviceConfig() micro.Config {
	conf := micro.Config{
		Name:        b.serviceName(),
		Version:     b.version(),
		Description: b.description(),
	}
	return conf
}

func (b *Bank) AddService(nc *nats.Conn) micro.Service {
	conf := b.serviceConfig()
	srv, err := micro.AddService(nc, conf)
	if err != nil {
		log.Fatalln(err)
	}
	b.service = srv
	b.addEndpoints()
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

func (b *Bank) createAccountRequest(req micro.Request) {
	r := NewAccountRequest{}
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

func (b Bank) getAccountRequest(req micro.Request) {
	_ = req.Respond([]byte("unimplemented"))
}

func (b *Bank) getCustomerAccountsRequest(req micro.Request) {
	_ = req.Respond([]byte("unimplemented"))
}

func (b *Bank) addEndpoints() {
	if b.service == nil {
		log.Fatalln("err: service isn't set")
	}
	g := b.service.AddGroup(b.serviceGroup())
	if err := g.AddEndpoint("create", micro.HandlerFunc(b.createAccountRequest)); err != nil {
		log.Fatalln(err)
	}
	if err := g.AddEndpoint("account", micro.HandlerFunc(b.getAccountRequest)); err != nil {
		log.Fatalln(err)
	}
	if err := g.AddEndpoint("accounts", micro.HandlerFunc(b.getCustomerAccountsRequest)); err != nil {
		log.Fatalln(err)
	}

	admin := b.service.AddGroup(b.adminGroup())
	if err := admin.AddEndpoint("deposit", micro.HandlerFunc(b.adminDepositRequest)); err != nil {
		log.Fatalln()
	}
	if err := admin.AddEndpoint("transfer", micro.HandlerFunc(b.adminTransferRequest)); err != nil {
		log.Fatalln()
	}
	if err := admin.AddEndpoint("hold", micro.HandlerFunc(b.adminHoldRequest)); err != nil {
		log.Fatalln()
	}
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
