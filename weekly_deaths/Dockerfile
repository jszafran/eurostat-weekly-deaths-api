# build frontend
FROM node:20-alpine as FrontendBuildStage

WORKDIR /usr/src/app
COPY frontend/package*.json .
RUN pwd
RUN ls -la
RUN npm install
COPY frontend/ .
RUN npm run build

# build Go app
FROM golang:1.20-alpine3.16 AS BuildStage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=FrontendBuildStage /usr/src/app/dist ./frontend/
RUN CGO_ENABLED=0 GOOS=linux go build -o /weekly_deaths_api

EXPOSE 8080

# final image
FROM alpine:latest

WORKDIR /
COPY --from=BuildStage /weekly_deaths_api /weekly_deaths_api

EXPOSE 8080

ARG COMMIT
ENV COMMIT=$COMMIT

ENTRYPOINT ["/weekly_deaths_api"]
