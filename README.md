<p align="center">  
  <img height="150" src="./gateway.png"  alt="Gateway" title="Gateway">
</p>

[![Build Status](https://travis-ci.org/ah8ad3/gateway.svg?branch=master)](https://travis-ci.org/ah8ad3/gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/ah8ad3/gateway)](https://goreportcard.com/report/github.com/ah8ad3/gateway)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fah8ad3%2Fgateway.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fah8ad3%2Fgateway?ref=badge_shield)
[![codecov](https://codecov.io/gh/ah8ad3/gateway/branch/master/graph/badge.svg)](https://codecov.io/gh/ah8ad3/gateway)

---

**Note:** This project are in ReWrite mode for standard behavior, may not be stable

---

This is an simple gateway written in Golang

# Features
- define all services in services.json file
- dynamic routes declaration for services
- RoutesV1 to change request params like middleware and more
- logging connect to mongoDb for save all logs
- can be a simple load balancer, in services.json you can add more server for one service and can act like load balancer
- environment variable support
- auth JWT and register added with middleware
- Health check of services every 1 hours can be modified
- Simple Rate limiter plugin per visitor
- Ip block list
- Api geo location info
- Aggregate requests dynamically by define in integrate.json by template
- DB manager with encryption to manage proxy and plugins
- Middleware per proxy
- Now You can get, add, update, delete proxies, also you can add, update, delete plugin of proxies with rest Api 
- K8s support, need to figure out strategy
- RouteV2 is here, path prefix to support all https, https, tcp

# Installation
### Release
you can find first type document [here](https://github.com/ah8ad3/gateway/blob/master/doc/README.md)
for use released versions

### Manual
```bash
  go get -t -v https://github.com/ah8ad3/gateway
  cd $GOPATH/src/github.com/ah8ad3/gateway
  make build
  ./dist/gateway secret   #for generate secret key
  ./dist/gateway load     # for load proxies from template json file
  ./dist/gateway run      # to run server
```

### Docker
```bash
1.
  docker build -t ah8ad3/gateway .
  docker run -d -p 3000:3000 ah8ad3/gateway
2.
  docker pull ah8ad3/gateway:latest
  docker run -d -p 3000:3000 ah8ad3/gateway
```

# TODO
- Schedule for requests life cycle inside of gateway
- Support envoy, kubernetes
- ApiInfo async
- Tests
- working on crud requests fo RouteV1
- some plugins: intelligent load balancer, ssl support, simple monitoring,
cache for gateway

# Constants
- services work with HTTP/1.1
- RoutesV1

- RouteV2
- fix all constants in RouteV1

----
# Routing

`RoutesV1`
this is version 1 of routing in gateway, it create routes of services dynamic in gateway and receive 
request and change it to send to services

it just get `x-www-urlencoded` form and send it to service like json

testing to fix all bugs

`RoutesV2`
this is second version of routing in this gateway, it's like rewirite reverce proxy to connect requests to 
services, also this Routation fix all RouteV1 constants

-----

> This api gateway create to make developers easier to work with services, not end

# Contributing
Just pick one of the Todo works and do it, or write some plugin for it,
in github standard mode PR and ...


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fah8ad3%2Fgateway.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fah8ad3%2Fgateway?ref=badge_large)