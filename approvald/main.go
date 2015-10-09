package main

import (
	"github.com/bearchinc/macaroon-spike-api"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	router := httprouter.New()
	api.Register(router)
	http.Handle("/", router)
}
