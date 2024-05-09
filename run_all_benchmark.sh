#!/usr/bin/env bash

# . ./scripts/benchmark_server.sh

. ./scripts/build.sh
./output/bin/anet_server -addr=":8000"
