FROM golang as builder

ARG app_env
ENV APP_ENV $app_env

COPY ./app /go/src/github.com/pilot114/micro_headers
WORKDIR /go/src/github.com/pilot114/micro_headers
RUN go get ./
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
WORKDIR /root
COPY --from=0 /go/src/github.com/pilot114/micro_headers/app .
CMD [ "./app"]