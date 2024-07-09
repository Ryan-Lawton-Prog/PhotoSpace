FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY ./.bin/api .
COPY ./config/docker-config.yml ./config/config.yml