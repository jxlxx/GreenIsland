package world

import (
	"time"

	"github.com/nats-io/nats.go"

	"github.com/jxlxx/GreenIsland/config"
	"github.com/jxlxx/GreenIsland/payloads"
	"github.com/jxlxx/GreenIsland/subjects"
)

type World struct {
	Countries    []Country
	HourDuration time.Duration
	nc           *nats.Conn
	current      payloads.WorldTick
	totalHours   int
}

func New() *World {
	return Init()
}

func Init() *World {

	// get countries

	// get companies

	world := &World{
		HourDuration: time.Millisecond,
		nc:           config.Connect(),
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
	quarter := (w.totalHours / 24) / 90

	return payloads.WorldTick{
		Quarter: quarter,
		Day:     day,
		Hour:    hour,
	}
}
