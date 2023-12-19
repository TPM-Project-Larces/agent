.PHONY: default run build test docs clean

# Variables
APP_NAME=agent

# Tasks
default: run-with-docs

run:
	@go run main.go
run-with-docs:
	@go mod tidy
	@GOPATH=/home/bruno/go PATH=$(PATH):/home/bruno/go/bin swag init
	@go run main.go
build:
	@go build -o $(APP_NAME) main.go
test:
	@go test ./ ...
docs:
	@GOPATH=/home/bruno/go PATH=$(PATH):/home/bruno/go/bin swag init
clean:
	@rm -f $(APP_NAME)
	@rm -rf ./docs