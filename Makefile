build-ginzanso:
	@go build -o bin/ginzanso cmd/scraper/ginzanso/main.go

run-ginzanso: build-ginzanso
	@./bin/ginzanso -rod="show,trace"

build:
	@go build -o bin/main cmd/main.go

run: build
	@./bin/main

test:
	@go test -v ./...

lint:
	@golangci-lint run -v