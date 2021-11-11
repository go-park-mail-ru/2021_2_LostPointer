FROM golang:1.17.3-buster AS build

WORKDIR /app

ENV CGO_ENABLED=0
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY ./. .

RUN apt update && apt install -y build-essential
RUN go build ./cmd/main/main.go


#ADD cmd/main .
#
#CMD ["./main"]
