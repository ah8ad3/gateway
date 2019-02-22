FROM golang:1.10

WORKDIR $GOPATH/src/github.com/ah8ad3/gateway

COPY . .

RUN go get -t ./...

EXPOSE 3000

ENTRYPOINT ["./entrypoint.sh"]
