FROM golang:1.17.3-alpine3.14 AS build

WORKDIR /app

ADD ../. .

ENV CGO_ENABLED=0
COPY ../go.mod .
COPY ../go.sum .
RUN go mod download
RUN go build ./cmd/profile/profile.go

WORKDIR /app
COPY --from=build /app/ .
ADD cmd/profile .

CMD ["./profile"]
