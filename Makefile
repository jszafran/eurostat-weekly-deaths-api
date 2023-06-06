build:
	go build -o bin/api cmd/api/*.go && chmod +x ./bin/api

run:
	go run cmd/api/*.go

test:
	go test ./...

testverbose:
	go test -v ./...
