#!/bin/bash

mkdir -p example
openssl genrsa -out example/server.key 4096
openssl req -nodes -new -x509 -sha256 -days 365 -config scripts/cert/localhost.cnf -extensions 'req_ext' -key example/server.key -out example/server.crt
