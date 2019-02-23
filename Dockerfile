FROM golang:1.10-alpine AS build_env

WORKDIR $GOPATH/src/github.com/ah8ad3/gateway

COPY . .

RUN apk add --no-cache bash git openssh

RUN go get -t .
RUN sh build/build.sh

FROM alpine
# RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN apk add --no-cache bash

WORKDIR /app
COPY --from=build_env $GOPATH/src/github.com/ah8ad3/gateway/dist/gateway /app
COPY --from=build_env $GOPATH/src/github.com/ah8ad3/gateway/entrypoint.sh /app

EXPOSE 3000

ENTRYPOINT ["./entrypoint.sh"]
