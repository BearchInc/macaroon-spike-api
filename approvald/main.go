package main

import (
	"github.com/bearchinc/macaroon-spike-api"
	"github.com/bearchinc/macaroon-spike-api/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	router := httprouter.New()
	db := database.New()
	defer db.Close()

	api.Register(router, db)

	http.ListenAndServe(":8080", router)
}
