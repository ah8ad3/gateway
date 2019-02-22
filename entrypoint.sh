#!/bin/bash

sh build/build.sh

./dist/gateway secret
./dist/gateway load
./dist/gateway run
