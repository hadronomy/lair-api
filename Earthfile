VERSION 0.8
FROM golang:1.23

COPY . /lair-api
WORKDIR /lair-api

deps:
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

build:
    FROM +deps
    RUN go build -o build/lair-api cmd/api/main.go
    SAVE ARTIFACT build/lair-api /lair-api AS LOCAL build/lair-api

docker:
    COPY +build/lair-api .
    ENTRYPOINT ["/bin/sh", "-c", "./lair-api"]
    SAVE IMAGE hadronomy/lair-api:latest

all:
    BUILD +build
    BUILD +docker