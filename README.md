<p align="center">  
  <img height="150" src="./gateway.png"  alt="Gateway" title="Gateway">
</p>

[![Build Status](https://travis-ci.org/ah8ad3/gateway.svg?branch=master)](https://travis-ci.org/ah8ad3/gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/ah8ad3/gateway)](https://goreportcard.com/report/github.com/ah8ad3/gateway)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fah8ad3%2Fgateway.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fah8ad3%2Fgateway?ref=badge_shield)

---

**Note:** This project may not be youre ideal gateway, we do our best if you want to help, we appreciate that

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
- REST for modifying services in progress
- Api geo location info
- Aggregate requests dynamically by define in integrate.json by template

# Installation
### Manual
```bash
    git clone https://github.com/ah8ad3/gateway
    cd gateway
    make build
    ./dist/gateway run
```

### Docker
Soon...

# TODO
- Plugins separate for each service
- Replace mongoDB, .env with another toy
- ApiInfo async
- working on crud requests fo RouteV1
- RoutesV2 to prefix route and act like bridge, rewrite route url
- some plugins: intelligent load balancer, ssl support, simple monitoring,
cache for gateway, nginx communication, kubernetes communication

# Constants
- services work with HTTP/1.1
- RoutesV1
----
# Routing

`RoutesV1`
this is version 1 of routing in gateway, it create routes of services dynamic in gateway and receive 
request and change it to send to services

it just get `x-www-urlencoded` form and send it to service like json

testing to fix all bugs

-----
I dont implement api gateway with all standard features because i dont need them now.
i will implement a simple api gateway for my private project, after that i will implement more good features


# Contributing
Just pick one of the Todo works and do it, or write some plugin for it,
in github standard mode PR and ...


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fah8ad3%2Fgateway.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fah8ad3%2Fgateway?ref=badge_large)