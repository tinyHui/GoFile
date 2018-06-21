.PHONY: vet build clean install

# Prepend our _vendor directory to the system GOPATH
# so that import path resolution will prioritize
# our third party snapshots.
GOPATH := ${PWD}/_vendor:${GOPATH}
ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export GOPATH

default: build

vet:
	go vet ./cmd/main.go

build: vet
	go build -v -o ./bin/main ./cmd/main.go

deploy-build: vet
	docker build . --tag emotech

deploy:
	mkdir -p ./storage
	docker run -d -e "GIN_MODE=release" -p 8080:8080 -v $(ROOT_DIR)/storage:/storage emotech

clean:
	rm -rf ./vendor

install:
	glide up

