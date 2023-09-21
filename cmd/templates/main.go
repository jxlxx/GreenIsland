package main

import "github.com/jxlxx/GreenIsland/world"

func main() {
	world.CreateTemplate(world.Country{}, "templates/country.yaml")
	world.CreateTemplate(world.Company{}, "templates/company.yaml")
}
