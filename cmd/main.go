package main

import (
	"log"

	"github.com/jxlxx/GreenIsland/config"
	"github.com/jxlxx/GreenIsland/world"
)

func main() {
	w := world.New()
	w.Connect()

	nc := config.Connect()
	defer func() {
		err := nc.Drain()
		if err != nil {
			log.Println(err)
		}
	}()

	srvs := w.AddServices(nc)
	for _, s := range srvs {
		s := s
		defer func() {
			err := s.Stop()
			if err != nil {
				log.Println(err)
			}
		}()
	}
	if err := w.Run(); err != nil {
		log.Fatalln(err)
	}
}
