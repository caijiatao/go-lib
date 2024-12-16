#!/bin/bash
trap "rm c_cache;kill 0" EXIT

go build -o c_cache
./c_cache -port=8081 &
./c_cache -port=8082 &
./c_cache -port=8083 &
