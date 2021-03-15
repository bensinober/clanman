PHONY: build all

help:  	## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

all: clean build-pi

clean:
	rm -rf ./build/*
	go vet ./src/...

run: ## Run amd64 version
	go run src/*.go

test: ## Run amd64 test
	go run src/*.go --test
build:  ## Build for linux amd64
	go build -o ./build/clanman src/*.go

build-pi:  ## build for RaspberryPi
	GOOS=linux GOARCH=arm GOARM=6 go build -o ./build/clanman src/*.go
