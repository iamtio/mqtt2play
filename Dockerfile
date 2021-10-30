FROM golang:1.16-alpine3.14 as build
WORKDIR /builddir
RUN apk add --no-cache 	build-base alsa-lib-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build ./cmd/mqtt2play-server

FROM alpine:3.14
WORKDIR /opt/mqtt2play
RUN apk add --no-cache alsa-utils
COPY --from=build /builddir/mqtt2play-server .
CMD ./mqtt2play