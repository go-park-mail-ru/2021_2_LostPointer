FROM golang:1.17.3-buster AS build

WORKDIR /app

#RUN make TARGET="authorization" build
ENV CGO_ENABLED=0
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY ./. .

RUN go build ./cmd/authorization/authorization.go

ADD cmd/authorization .

CMD ["./authorization"]
