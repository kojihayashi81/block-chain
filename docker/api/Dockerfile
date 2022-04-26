FROM golang:1.18-alpine

WORKDIR /go/src/app

RUN apk update \ 
    && apk add git \
    && apk add alpine-sdk \
    && go install github.com/cosmtrek/air@latest \
    && go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install github.com/golang/mock/mockgen@v1.6.0
