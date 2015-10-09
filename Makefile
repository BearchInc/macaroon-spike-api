update:
	goapp get github.com/drborges/appx
	goapp get gopkg.in/unrolled/render.v1
	goapp get gopkg.in/macaroon.v1

run: update
	goapp serve approvald/main.go