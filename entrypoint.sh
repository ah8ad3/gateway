#!/bin/bash

set -e

./gateway secret
./gateway load
./gateway run
