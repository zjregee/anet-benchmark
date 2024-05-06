#!/usr/bin/env bash

rm -rf output/
mkdir -p output/bin/
mkdir -p output/log/

go build -v -o output/bin/net_server ./net

go build -v -o output/bin/net_client ./net/client
