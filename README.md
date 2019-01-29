# gateway
simple api gateway with golang

# Features
- define all services in services.json file
- dynamic routes declaration for services
- RoutesV1 to change request params like middleware and more
- logging connect to mongoDb for save all logs

# TODO
- working on crud requests
- working on jwt and auth staff
- RoutesV2 to prefix route and act like bridge, rewrite route url
- combining services results

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
