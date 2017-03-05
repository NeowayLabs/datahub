all: build

build:
	go build ./cmd/datahub

run: build
	./datahub
