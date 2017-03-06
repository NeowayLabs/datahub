version ?= "latest"

all: build

build:
	go build ./cmd/datahub

run: build
	./datahub

image:
	docker build . -t neowaylabs/datahub

shell: image
	docker run -ti neowaylabs/datahub /bin/bash

publish: image
	docker push neowaylabs/datahub:$(version)
