build:
	@go build -o bin/main cmd/ginzanso/main.go

run: build
	@./bin/main -rod="show,slow=1s,trace"

test:
	@go test -v ./...

lint:
	@golangci-lint run -v