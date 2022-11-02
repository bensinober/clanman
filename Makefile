PHONY:	build all clean

help:  	## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

all: clean build-pi

clean:
	rm -rf ./dist/*
	cd src && go vet ./...

run: ## Run amd64 version
	cd src && go run clanman.go server.go controls.go menu.go display.go sampler.go

test: ## Run amd64 test
	cd src && go run clanman.go server.go controls.go menu.go display.go sampler.go --test

build: ## Build for linux amd64
	cd src && go build -o ../dist/clanman clanman.go server.go controls.go menu.go display.go sampler.go
	cp src/*.json ./dist/

build-pi:  ## build for RaspberryPi
	cd src && GOOS=linux GOARCH=arm GOARM=6 go build -o ../dist/clanman
	cp src/*.json ./dist/

push:	## push dist to patchbox
	rsync -avz dist/* pi@clanbox:./clanman/
