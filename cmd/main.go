package main

import (
	"training/api"
	"training/configuration"
	"training/internal/propose"
)
func main() {

	db := configuration.PrepareDatabase()
	p := propose.NewMain(db)
	a := api.Api{
		ProposeMain: p,
	}
	router := a.NewApi()
	a.Start(router)
}
