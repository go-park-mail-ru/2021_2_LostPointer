FROM golang:1.17-buster AS build

WORKDIR /app

ADD ../. .

#RUN make TARGET="authorization" build
ENV CGO_ENABLED=0
COPY ../go.mod .
COPY ../go.sum .
RUN go mod download
RUN go build ./cmd/authorization/authorization.go

ENV AUTH_PORT=":3080"

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/ .
ADD cmd/authorization .

CMD ["./authorization"]

EXPOSE 3080