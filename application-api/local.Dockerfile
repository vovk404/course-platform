FROM golang:1.20-alpine as build

RUN apk update && apk add bash

ARG POSTGRESQL_HOST
ARG POSTGRESQL_USER
ARG POSTGRESQL_PASSWORD
ARG POSTGRESQL_DATABASE

ARG APP_BASE_URL


ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

ENV POSTGRESQL_HOST $POSTGRESQL_HOST
ENV POSTGRESQL_USER $POSTGRESQL_USER
ENV POSTGRESQL_PASSWORD $POSTGRESQL_PASSWORD
ENV POSTGRESQL_DATABASE $POSTGRESQL_DATABASE

ENV APP_BASE_URL $APP_BASE_URL

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o app/api cmd/main.go

EXPOSE 8082

CMD ["app/api"]