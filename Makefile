.PHONY: build run test

build:
	go build -o bin/warung-online main.go

run:
	go run main.go

test:
	go test ./...

swagger:
	swag init