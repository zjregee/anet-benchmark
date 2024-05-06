#!/usr/bin/env bash

n=1000000
body=(1024)
concurrent=(10 100 500 1000)
repos=("net")
ports=(8001)

. ./scripts/env.sh
. ./scripts/build.sh

benchmark
