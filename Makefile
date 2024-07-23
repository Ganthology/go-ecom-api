build:
	@go build -o bin/go-ecom-api cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/go-ecom-api