FROM golang:1.13-alpine

COPY bin/translationApi /go/bin/translationApi

ENV GIN_MODE=release

ENTRYPOINT /go/bin/translationApi
EXPOSE 8080
