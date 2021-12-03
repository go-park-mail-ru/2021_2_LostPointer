FROM golang:1.17.3-alpine3.14 AS build

WORKDIR /app

ENV CGO_ENABLED=0
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

RUN mkdir users
RUN mkdir playlists

RUN apk add --no-cache --update libwebp-tools
RUN mkdir .bin && mkdir .bin/webp && cp $(which cwebp) .bin/webp/

COPY ./. .

RUN GOARCH=arm64 GOOS=linux go build ./cmd/main/main.go

ADD cmd/main .

CMD ["./main"]
