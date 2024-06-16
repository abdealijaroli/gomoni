build:
	@go build -o bin/gomoni

run: build
	@./bin/gomoni

test: 
	@go test -v ./...