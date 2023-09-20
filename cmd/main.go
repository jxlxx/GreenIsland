package main

import (
	"runtime"

	"github.com/jxlxx/GreenIsland/world"
)

func main() {

	w := world.New()

	go w.Run()

	runtime.Goexit()

}
