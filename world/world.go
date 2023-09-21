package world

import (
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
	HourDuration time.Duration
	nc           *nats.Conn
	current      payloads.WorldTick
	totalHours   int
	countries    []*Country
	companies    []*Company
}

func New() *World {
	return Init()
}

func Init() *World {
	countries := createCountries()
	companies := createCompanies()
	world := &World{
		HourDuration: time.Millisecond,
		nc:           config.Connect(),
		countries:    countries,
		companies:    companies,
	}
	return world
}

func (w *World) Run() {
	for {
		tick := w.Tick()
		w.nc.Publish(subjects.TickHour.String(), payloads.Bytes(tick))

		if tick.Day != w.current.Day {
			w.nc.Publish(subjects.TickDay.String(), payloads.Bytes(tick))
		}

		if tick.Quarter != w.current.Quarter {
			w.nc.Publish(subjects.TickQuarter.String(), payloads.Bytes(tick))
		}
		w.current = tick
		time.Sleep(w.HourDuration)
	}
}

func (w *World) Tick() payloads.WorldTick {
	w.totalHours += 1

	hour := w.totalHours % 24
	day := (w.totalHours / 24) % 90
	quarter := ((w.totalHours/24)/90)%4 + 1

	return payloads.WorldTick{
		Quarter: quarter,
		Day:     day,
		Hour:    hour,
	}
}

func createCountries() []*Country {
	files := []string{
		"data/countries/canada.yaml",
		"data/countries/usa.yaml",
	}
	return create(files, Country{})
}

func createCompanies() []*Company {
	files := []string{
		"data/company/aerospin.yaml",
	}
	return create(files, Company{})
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

func CreateTemplate() {
	country := Country{}
	data, err := yaml.Marshal(&country)
	if err != nil {
		log.Fatalln(err)
	}
	if err = os.WriteFile("templates/country.yaml", data, 0644); err != nil {
		log.Fatalln(err)
	}
}
