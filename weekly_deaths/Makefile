BINARY_NAME=weekly_deaths_backend

build:
	go build -o ./bin/${BINARY_NAME} && chmod +x ./bin/${BINARY_NAME}

run: build
	./bin/${BINARY_NAME}

test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage-report:
	go test ./... -coverprofile=coverage.out

test-check-coverage-report:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

run-sandbox:
	cd cmd/sandbox && go run .

run-frontend:
	cd frontend && npm run serve -- --port 9000

build-frontend:
	cd frontend && npm run build
