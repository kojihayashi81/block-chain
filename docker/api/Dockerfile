FROM golang:1.18-alpine

WORKDIR /go/src/api

COPY ./api/go.mod .
COPY ./api/go.sum .
RUN go mod download

COPY ./api .

RUN go install github.com/cosmtrek/air@latest
