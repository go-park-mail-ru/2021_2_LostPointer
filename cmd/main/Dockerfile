#Build
FROM golang:1.17.3-alpine3.15 AS build

WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY ./. .

RUN mkdir {users,playlists}
RUN apk add --no-cache --update libwebp-tools

RUN mkdir .bin && mkdir .bin/webp && cp $(which cwebp) .bin/webp/

RUN go build ./cmd/main/main.go

#Environment
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/main .
COPY --from=build /app/users .
COPY --from=build /app/playlists .
COPY --from=build /app/.bin .

CMD ["./main"]
