build-ginzanso:
	@go build -o bin/ginzanso cmd/ginzanso/main.go

run-ginzanso: build-ginzanso
	@./bin/ginzanso -rod="show,slow=1s,trace"

build-gmail:
	@go build -o bin/gmail cmd/gmail/main.go

run-gmail: build-gmail
	@./bin/gmail

test:
	@go test -v ./...

lint:
	@golangci-lint run -v