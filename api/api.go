package api

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	propose2 "training/api/v1/propose"
	"training/internal/propose"
)

type Api struct {
	ProposeMain propose.ProposeUsecases
}


func (api Api) Start(router *mux.Router){
	server := &http.Server{
		Handler: router,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
func(api Api) NewApi() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/propose", propose2.Create(api.ProposeMain)).Methods("POST")
	return router
}