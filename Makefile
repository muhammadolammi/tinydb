build:
	@go build -o bin/tinydb .

run: build
	@./bin/tinydb

test:
	go test ./... -v