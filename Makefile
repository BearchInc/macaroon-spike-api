update:
	go get gopkg.in/mgo.v2
	go get github.com/julienschmidt/httprouter
	go get gopkg.in/macaroon.v1

run: update
	go run approvald/main.go