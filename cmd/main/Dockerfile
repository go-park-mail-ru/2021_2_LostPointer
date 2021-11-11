FROM golang:1.17.3-alpine3.14 AS build

WORKDIR /app

ENV CGO_ENABLED=0
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY ./. .

RUN go build ./cmd/main/main.go

ADD cmd/main .

CMD ["./main"]
