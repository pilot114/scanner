FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY ./app /go/src/github.com/pilot114/micro_headers/app
WORKDIR /go/src/github.com/pilot114/micro_headers/app

RUN go get ./
RUN go build
