package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/bearchinc/macaroon-spike-api/models"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/macaroon.v1"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()
	db := NewDB()
	defer db.Close()

	router.POST("/deployments", func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		userParams := struct{ UserID, Commit string }{}
		json.NewDecoder(req.Body).Decode(&userParams)

		err := db.Save("deployments", &models.Deployment{
			UserID: userParams.UserID,
			Commit: userParams.Commit,
			Status: models.DeploymentPending,
		})

		if err != nil {
			log.Fatal(err)
		}

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

type DB struct {
	session *mgo.Session
	mongo   *mgo.Database
}

func NewDB() *DB {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	return &DB{session, session.DB("db")}
}

func (db *DB) Save(collectionName string, data interface{}) error {
	return db.mongo.C(collectionName).Insert(data)
}

func (db *DB) Close() {
	db.session.Close()
}
