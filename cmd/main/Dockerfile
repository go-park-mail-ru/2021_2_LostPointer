FROM golang:1.17-buster AS build

WORKDIR /app

ADD ../. .

ENV CGO_ENABLED=0
COPY ../go.mod .
COPY ../go.sum .
RUN go mod download
RUN go build ./cmd/main/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/ .
ADD cmd/main .

CMD ["./main"]
