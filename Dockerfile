FROM golang:1.8.0

COPY ./cmd /go/src/github.com/NeowayLabs/datahub/cmd
COPY ./*.go /go/src/github.com/NeowayLabs/datahub

WORKDIR /go/src/github.com/NeowayLabs/datahub

RUN go build ./cmd/datahub
