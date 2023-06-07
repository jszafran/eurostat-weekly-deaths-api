build:
	go build -o bin/api cmd/api/*.go && chmod +x ./bin/api

run:
	go run cmd/api/*.go

test:
	go test ./...

testverbose:
	go test -v ./...

coveragereport:
	go test ./... -coverprofile=coverage.out

coveragehtml:
		go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
