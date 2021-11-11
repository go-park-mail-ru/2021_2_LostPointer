FROM golang:1.17.3-alpine3.14 AS build

WORKDIR /app

ENV CGO_ENABLED=0
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY ./. .

RUN apk add --no-cache --update libpng-dev libjpeg-turbo-dev giflib-dev tiff-dev autoconf automake make gcc g++ wget && \
    wget https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-0.6.0.tar.gz && \
    tar -xvzf libwebp-0.6.0.tar.gz && \
    mv libwebp-0.6.0 libwebp && \
    rm libwebp-0.6.0.tar.gz && \
    cd /app/libwebp && \
    ./configure && \
    make && \
    make install && \
    rm -rf libwebp

RUN go build ./cmd/main/main.go


ADD cmd/main .

CMD ["./main"]
