package api

import (
	"appengine"
	"encoding/base64"
	"encoding/json"
	"github.com/bearchinc/macaroon-spike-api/models"
	"github.com/drborges/appx"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/macaroon.v1"
	"gopkg.in/unrolled/render.v1"
	"log"
	"net/http"
)

var (
	ApproverID  = "approver:ygor"
	ApproverKey = []byte("Ygor's secret")
	ApprovalURL = "http://localhost:6060/approvals?from=Ygor"
)

type JSON map[string]interface{}

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
	macaroon.AddThirdPartyCaveat(ApproverKey, ApproverID, ApprovalURL)
	return macaroon
}

func Register(router *httprouter.Router) {
	r := render.New()

	router.POST("/deployments", func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		deployment := DeploymentFrom(req)
		macaroon := CreateDeploymentMacaroon(deployment)

		db := appx.NewDatastore(appengine.NewContext(req))
		if err := db.Save(deployment); err != nil {
			log.Panic(err)
		}

		b, _ := macaroon.MarshalBinary()
		r.JSON(w, 200, JSON{
			"token": base64.URLEncoding.EncodeToString(b),
		})
	})
}
