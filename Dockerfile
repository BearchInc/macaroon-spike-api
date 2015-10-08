FROM golang

ADD . /go/src/github.com/bearchinc/macaroons-spike-api

RUN go install github.com/bearchinc/macaroons-spike-api/approvald

ENTRYPOINT /go/bin/approvald

EXPOSE 8080
