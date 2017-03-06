FROM golang:1.8.0

RUN apt-get update
RUN apt-get install -y r-base r-base-dev libgtk2.0

COPY ./examples /go/src/github.com/NeowayLabs/datahub/examples

WORKDIR /go/src/github.com/NeowayLabs/datahub

RUN R -f ./examples/installdeps.r

COPY ./cmd /go/src/github.com/NeowayLabs/datahub/cmd
COPY ./scientists /go/src/github.com/NeowayLabs/datahub/scientists
COPY ./company /go/src/github.com/NeowayLabs/datahub/company
COPY ./db /go/src/github.com/NeowayLabs/datahub/db
COPY ./tools /go/src/github.com/NeowayLabs/datahub/tools
COPY ./*.go /go/src/github.com/NeowayLabs/datahub


RUN go get ./...
RUN go build ./cmd/datahub

CMD ./datahub
