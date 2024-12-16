#!/bin/bash

go build -o c_cache
./c_cache -port=8001 &
./c_cache -port=8002 &
./c_cache -port=8003 &
