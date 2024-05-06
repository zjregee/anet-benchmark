#!/usr/bin/env bash

n=1000000
body=(24 1024 2048 4096)
concurrent=(2 8 16 32 128 256)
repos=("anet")
ports=(8001)

. ./scripts/env.sh
. ./scripts/build.sh

benchmark
