package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/macaroon.v1"
	"gopkg.in/mgo.v2"
	"encoding/base64"
	"log"
	"encoding/json"
)

func main() {
	router := httprouter.New()

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	router.POST("/deployments", func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		userParams := struct { UserID, Commit string } {}
		json.NewDecoder(req.Body).Decode(&userParams)
		

		log.Printf("#### %+v", userParams)
		macaroon, _ := macaroon.New([]byte("lol"), "macarrão", "localhost:8080")
		macaroon.AddFirstPartyCaveat("user:" + userParams.UserID)
		macaroon.AddFirstPartyCaveat("commit:" + userParams.Commit)

		macaroon.AddThirdPartyCaveat([]byte("Ygor approval key"), "Ygor seal of approval", "http://localhost:6060/approvals?from=Ygor")

		j, _ := macaroon.MarshalJSON()
		b, _ := macaroon.MarshalBinary()
		log.Println("####", string(j))
		token := base64.URLEncoding.EncodeToString(b)
		w.Write([]byte(token))
	})

	http.ListenAndServe(":8080", router)
}