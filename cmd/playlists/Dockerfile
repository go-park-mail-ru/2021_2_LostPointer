#Build
FROM golang:1.17.3-buster AS build

WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY ./. .
RUN CGO_ENABLED=0 go build ./cmd/playlists/playlists.go

#Environment
FROM alpine:latest

WORKDIR /app
COPY --from=build /app .

CMD ["./playlists"]
