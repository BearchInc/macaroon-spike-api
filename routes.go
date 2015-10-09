package api

import (
	"appengine"
	"appengine/urlfetch"
	"encoding/base64"
	"encoding/json"
	"github.com/bearchinc/macaroon-spike-api/gcm"
	"github.com/bearchinc/macaroon-spike-api/models"
	"github.com/drborges/appx"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/macaroon.v1"
	"gopkg.in/unrolled/render.v1"
	"log"
	"net/http"
)

var (
	GCMApiKey   = "AIzaSyD4jrcwQEsQrbHdhbkn22NWPH2tAByr-Jo"
	ApproverID  = "approver:ygor"
	ApproverGCM = "mR8qPvS7aQ0:APA91bFIlvG2e23PcrYvU-mJV7yIh-Gl3Re-Esjl0DLQxA6TYYxNWohaedos4v3Ed7JB31yb1ZBlb2je8YjWiffPraqe1GbrC3QZDLwqZvJmirfnDXCSv6RXeBYYjZIOi8SGnOBRcmu0"
	ApproverKey = []byte("Ygor's secret")
	ApprovalURL = "http://localhost:6060/approvals?from=Ygor"
)

type JSON map[string]interface{}

type MacaroonForm struct {
	Token string "json:token"
}

func DeploymentFrom(req *http.Request) *models.Deployment {
	deployment := &models.Deployment{}
	json.NewDecoder(req.Body).Decode(deployment)
	deployment.Approver = ApproverID
	deployment.Status = models.DeploymentAwaitingApproval
	return deployment
}

func CreateDeploymentMacaroon(deployment *models.Deployment) *macaroon.Macaroon {
	macaroon, _ := macaroon.New([]byte("lol"), "macarrão", "localhost:8080")
	macaroon.AddFirstPartyCaveat("requester:" + deployment.Requester)
	macaroon.AddFirstPartyCaveat("commit:" + deployment.Commit)
	//	macaroon.AddThirdPartyCaveat(ApproverKey, ApproverID, ApprovalURL)
	return macaroon
}

func RequestApproval(c appengine.Context, deployment *models.Deployment, macaroon *macaroon.Macaroon) {
	sender := gcm.NewSenderWithHttpClient("AIzaSyD4jrcwQEsQrbHdhbkn22NWPH2tAByr-Jo", urlfetch.Client(c))
	// TODO handle response...
	err := sender.Send(gcm.Message{
		To:               "mR8qPvS7aQ0:APA91bFIlvG2e23PcrYvU-mJV7yIh-Gl3Re-Esjl0DLQxA6TYYxNWohaedos4v3Ed7JB31yb1ZBlb2je8YjWiffPraqe1GbrC3QZDLwqZvJmirfnDXCSv6RXeBYYjZIOi8SGnOBRcmu0",
		ContentAvailable: true,
		Notification: gcm.Params{
			"title": "@Diego Borges",
			"body":  "Big Lolzis it vorks!",
		},
	})

	if err != nil {
		log.Print(err)
	}
}

func VerifyMacaroon(serializedMacaroon []byte) error {
	macaroon, _ := macaroon.New([]byte("lol"), "ss", "s")
	err := macaroon.UnmarshalBinary(serializedMacaroon)
	if err != nil {
		return err
	}

	err = macaroon.Verify([]byte("lol"), CaveatCheck, nil)
	return err
}

func CaveatCheck(d string) error { return nil }

func Register(router *httprouter.Router) {
	r := render.New()

	router.POST("/deployments", func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		deployment := DeploymentFrom(req)
		macaroon := CreateDeploymentMacaroon(deployment)
		RequestApproval(appengine.NewContext(req), deployment, macaroon)

		db := appx.NewDatastore(appengine.NewContext(req))
		if err := db.Save(deployment); err != nil {
			log.Panic(err)
		}

		b, _ := macaroon.MarshalBinary()
		token := base64.URLEncoding.EncodeToString(b)

		r.JSON(w, 200, JSON{
			"token": token,
		})
	})

	router.POST("/validate", func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		form := &MacaroonForm{}
		json.NewDecoder(req.Body).Decode(form)
		bytes, err := base64.URLEncoding.DecodeString(form.Token)
		if err != nil {
			r.JSON(w, 400, JSON{
				"message": "Error deserializing macaroon.",
				"error": err.Error(),
			})
			return
		}

		err = VerifyMacaroon(bytes)
		if err != nil {
			r.JSON(w, 400, JSON{
				"message": "Macaroon invalid.",
				"error": err.Error(),
			})
			return
		}

		r.JSON(w, 200, JSON{
			"message": "Macaroon valid.",
		})
	})
}
