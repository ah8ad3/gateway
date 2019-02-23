FROM golang:1.10 AS build_env

WORKDIR $GOPATH/src/github.com/ah8ad3/gateway

COPY . .

RUN go get -t ./...
RUN sh build/build.sh

FROM alpine
# RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR $GOPATH/src/github.com/ah8ad3/gateway
COPY --from=build_env $GOPATH/src/github.com/ah8ad3/gateway/dist/gateway .

EXPOSE 3000

ENTRYPOINT ["./entrypoint.sh"]
