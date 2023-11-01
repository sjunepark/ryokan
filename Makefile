build:
	@go build -o bin/main cmd/ginzanso/main.go

run: build
	@./bin/main -rod=show

test:
	@go test -v ./...