#Build
FROM golang:1.17.3-alpine3.14 AS build

WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
RUN mkdir users
RUN mkdir playlists
RUN apk add --no-cache --update libwebp-tools
COPY ./. .
RUN mkdir .bin && mkdir .bin/webp && cp $(which cwebp) .bin/webp/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ./cmd/main/main.go

#Environment
FROM alpine:latest

WORKDIR /app
COPY --from=build /app .

CMD ["./main"]