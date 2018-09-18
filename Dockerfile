FROM golang:alpine

ADD ./bin/translationApi /go/bin/translationApi
ENTRYPOINT /go/bin/translationApi
EXPOSE 8080