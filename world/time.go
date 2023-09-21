package world

import (
	"log"
	"time"

	"github.com/jxlxx/GreenIsland/payloads"
	"github.com/jxlxx/GreenIsland/subjects"
)

func (w *World) Run() {

	w.connect()

	for {
		tick := w.Tick()

		if err := w.nc.Publish(subjects.TickHour.String(), tick); err != nil {
			log.Println(err)
		}

		if err := w.nc.Publish(subjects.DailyTick(tick.Quarter, tick.Day, tick.Hour), tick); err != nil {
			log.Println(err)
		}

		if tick.Day != w.current.Day {
			if err := w.nc.Publish(subjects.TickDay.String(), tick); err != nil {
				log.Println(err)
			}
		}

		if tick.Quarter != w.current.Quarter {
			if err := w.nc.Publish(subjects.TickQuarter.String(), tick); err != nil {
				log.Println(err)
			}
		}
		w.current = tick
		time.Sleep(w.HourDuration)
	}
}

func (w *World) connect() {

}

func (w *World) Tick() payloads.WorldTick {
	w.totalHours += 1
	w.elaspsedRealTime += w.HourDuration

	hour := w.totalHours % 24
	day := (w.totalHours / 24) % 90
	quarter := ((w.totalHours/24)/90)%4 + 1

	return payloads.WorldTick{
		Quarter: quarter,
		Day:     day,
		Hour:    hour,
		EGT:     w.totalHours,
		ERT:     w.elaspsedRealTime,
	}
}
