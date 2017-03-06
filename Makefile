version ?= "latest"

all: build

build:
	go get ./...
	go build ./cmd/datahub

run: image
	docker run -ti --rm --net=host neowaylabs/datahub

image:
	docker build . -t neowaylabs/datahub

shell: image
	docker run -ti neowaylabs/datahub /bin/bash

publish: image
	docker push neowaylabs/datahub:$(version)
