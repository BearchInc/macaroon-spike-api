package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		w.Write([]byte("Success!"))
	})

	http.ListenAndServe(":8080", router)
}