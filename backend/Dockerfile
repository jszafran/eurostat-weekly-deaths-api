FROM golang:1.20-alpine3.16 AS BuildStage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go generate cmd/api/main.go && go build -o /weekly_deaths_api ./cmd/api

EXPOSE 8080

FROM alpine:latest

WORKDIR /

COPY --from=BuildStage /weekly_deaths_api /weekly_deaths_api

EXPOSE 8080

ENTRYPOINT ["/weekly_deaths_api"]
