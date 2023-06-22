FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go generate && go build -o /weekly_deaths_api ./cmd/api

EXPOSE 8080

CMD ["/weekly_deaths_api"]
