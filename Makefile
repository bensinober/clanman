PHONY:	build all clean

help:  	## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

all: clean build-pi

clean:
	rm -rf ./dist/*
	go vet src/clanman.go src/server.go src/controls.go src/menu.go src/display.go src/fluid.go

run: ## Run amd64 version
	go run src/clanman.go src/server.go src/controls.go src/menu.go src/display.go src/fluid.go

test: ## Run amd64 test
	cd src && go run clanman.go server.go controls.go menu.go display.go fluid.go --test

build: ## Build for linux amd64
	go build -o ./dist/clanman src/clanman.go src/server.go src/controls.go src/menu.go src/display.go src/fluid.go
	cp src/*.json dist/

build-pi:  ## build for RaspberryPi
	GOOS=linux GOARCH=arm GOARM=6 go build -o ./dist/clanman src/*.go
	cp src/*.json dist/

push:	## push dist to patchbox
	rsync -avz dist/* patch@patchbox:./clanman/
