package world

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"gopkg.in/yaml.v3"

	"github.com/jxlxx/GreenIsland/config"
	"github.com/jxlxx/GreenIsland/payloads"
	"github.com/jxlxx/GreenIsland/subjects"
)

type World struct {
	HourDuration     time.Duration
	elaspsedRealTime time.Duration
	nc               *nats.Conn
	current          payloads.WorldTick
	totalHours       int
	countries        []*Country
	companies        []*Company
}

func New() *World {
	return Init()
}

func Init() *World {
	nc := config.Connect()
	conn, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	countries := createCountries()
	companies := createCompanies()

	for _, c := range countries {
		if _, err := conn.Subscribe(subjects.TickDay.String(), c.DailySubscriber()); err != nil {
			fmt.Println(err)
		}
		if _, err := conn.Subscribe(subjects.TickQuarter.String(), c.QuarterlySubscriber()); err != nil {
			fmt.Println(err)
		}
	}

	for _, c := range companies {
		if _, err := conn.Subscribe(subjects.TickDay.String(), c.DailySubscriber()); err != nil {
			fmt.Println(err)
		}
		if _, err := conn.Subscribe(subjects.TickQuarter.String(), c.QuarterlySubscriber()); err != nil {
			fmt.Println(err)
		}
	}

	now := time.Now()
	world := &World{
		HourDuration:     time.Microsecond * 500,
		nc:               nc,
		countries:        countries,
		companies:        companies,
		elaspsedRealTime: now.Sub(now),
	}
	return world
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
