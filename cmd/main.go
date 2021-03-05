package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"training/api"
	"training/configuration"
	"training/internal/propose"
)
func main() {
	configuration.PrepareDatabase()
	dsn := "host=localhost user=keycloak password=password dbname=keycloak port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil {
		panic(err)
	}
	p := propose.NewMain(db)
	a := api.Api{
		ProposeMain: p,
	}
	router := a.NewApi()
	a.Start(router)
}
