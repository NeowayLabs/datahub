version ?= "latest"

all: build

build:
	go build ./cmd/datahub

run: build
	./datahub

image:
	docker build . -t neowaylabs/datahub

publish: image
	docker push neowaylabs/datahub:$(version)
