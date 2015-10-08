FROM golang

ADD . /go/src/github.com/bearchinc/macaroons-spike-api

RUN go get github.com/julienschmidt/httprouter
RUN go get gopkg.in/macaroon.v1
RUN go install github.com/bearchinc/macaroons-spike-api/approvald

ENTRYPOINT /go/bin/approvald

EXPOSE 8080
