APP?=freedom-countries

.PHONY: build
## build: build the application
build:
	@echo "Building..."
	@go build -o ${APP} *.go

.PHONY: build-windows
## build: build the application
build-windows:
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o ${APP} *.go

.PHONY: run
## run: runs go run main.go
run:
	go run -race *.go $(ARGS)

.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
