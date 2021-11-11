FROM golang:1.17.3-buster AS build

WORKDIR /app

ENV CGO_ENABLED=0
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY ./. .

RUN go build ./cmd/music/music.go

ADD cmd/music .

CMD ["./music"]
