package world

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"gopkg.in/yaml.v3"

	"github.com/jxlxx/GreenIsland/config"
	"github.com/jxlxx/GreenIsland/payloads"
	"github.com/jxlxx/GreenIsland/subjects"
)

type World struct {
	HourDuration     time.Duration
	elaspsedRealTime time.Duration
	nc               *nats.EncodedConn
	current          payloads.WorldTick
	totalHours       int
	countries        []*Country
	companies        []*Company
	adminService     micro.Service
}

func New() *World {
	countries := createCountries()
	companies := createCompanies()

	now := time.Now()
	world := &World{
		HourDuration:     time.Microsecond * 500,
		countries:        countries,
		companies:        companies,
		elaspsedRealTime: now.Sub(now),
	}
	return world
}

func (w *World) Connect() {
	w.nc = config.EncodedConnect()
	for _, c := range w.countries {
		if _, err := w.nc.Subscribe(subjects.TickDay.String(), c.DailySubscriber()); err != nil {
			fmt.Println(err)
		}
		if _, err := w.nc.Subscribe(subjects.TickQuarter.String(), c.QuarterlySubscriber()); err != nil {
			fmt.Println(err)
		}
		for _, b := range c.CommercialBanks {
			b.Connect()
			b.Setup()
		}
	}

	for _, c := range w.companies {
		if _, err := w.nc.Subscribe(subjects.TickDay.String(), c.DailySubscriber()); err != nil {
			fmt.Println(err)
		}
		if _, err := w.nc.Subscribe(subjects.TickQuarter.String(), c.QuarterlySubscriber()); err != nil {
			fmt.Println(err)
		}
	}
}

func (w *World) SetCompanyBankAccounts() {
	for _, c := range w.companies {
		c.InitializeCompany()
	}
}

func (w *World) AdminConfig() micro.Config {
	return micro.Config{
		Name:    "AdminService",
		Version: config.GetEnvOrDefault("VERSION", "0.0.1"),
	}
}

func (w *World) AdminService(nc *nats.Conn) micro.Service {
	conf := w.AdminConfig()
	srv, err := micro.AddService(nc, conf)
	if err != nil {
		log.Fatalln(err)
	}
	w.adminService = srv
	w.AddEndpoints()
	return w.adminService
}

func (w *World) AddEndpoints() {
}

func (w *World) AddServices(nc *nats.Conn) []micro.Service {
	w.AdminService(nc)
	services := []micro.Service{w.adminService}
	for _, c := range w.countries {
		for _, b := range c.CommercialBanks {
			services = append(services, b.AddService(nc))
		}
	}
	return services
}

func (w *World) Initialize() {
	for _, c := range w.countries {
		c.Initialize()
	}
}

func createCountries() []*Country {
	files := []string{
		"/data/countries/canada.yaml",
		"/data/countries/usa.yaml",
	}
	countries := create(files, Country{})
	for _, c := range countries {
		nc := config.EncodedConnect()
		c.nc = nc
	}
	return countries
}

func createCompanies() []*Company {
	files := []string{
		"/data/companies/aerospin.yaml",
	}
	companies := create(files, Company{})
	for _, c := range companies {
		nc := config.EncodedConnect()
		c.nc = nc
	}
	return companies
}

func create[T any](files []string, t T) []*T {
	slice := []*T{}
	for _, f := range files {
		var cc T
		config.ReadConfig(f, &cc)
		slice = append(slice, &cc)
	}
	return slice
}

func CreateTemplate[T any](t T, filename string) {
	var c T
	data, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalln(err)
	}
	if err = os.WriteFile(filename, data, 0644); err != nil {
		log.Fatalln(err)
	}
}
