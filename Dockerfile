FROM golang:1.10-alpine AS build_env

WORKDIR $GOPATH/src/github.com/ah8ad3/gateway

COPY . .

RUN apk add --no-cache bash git openssh

RUN go get -t .
RUN sh build/build.sh

FROM alpine
RUN apk update && apk add ca-certificates bash && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=build_env /go/src/github.com/ah8ad3/gateway/dist/gateway .
COPY --from=build_env /go/src/github.com/ah8ad3/gateway/entrypoint.sh .
COPY --from=build_env /go/src/github.com/ah8ad3/gateway/services.json .
COPY --from=build_env /go/src/github.com/ah8ad3/gateway/integrates.json .

EXPOSE 3000

RUN ./gateway secret
RUN ./gateway load

CMD ["./gateway", "run"]
