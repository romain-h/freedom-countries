APP?=freedom-countries

.PHONY: build
## build: build the application
build:
	@echo "Building..."
	@go build -o ${APP} cmd/update/main.go

.PHONY: run
## run: runs go run main.go
run:
	go run -race cmd/update/main.go $(ARGS)

.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
