package api

import (
	"encoding/base64"
	"encoding/json"
	"github.com/bearchinc/macaroon-spike-api/database"
	"github.com/bearchinc/macaroon-spike-api/models"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/macaroon.v1"
	"log"
	"net/http"
)

func Register(router *httprouter.Router, db *database.Mongo) {
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
}
