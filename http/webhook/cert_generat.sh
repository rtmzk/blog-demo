#!/usr/bin/env bash

[ ! -d  certs ] && mkdir certs

openssl genrsa -out certs/server.key 2048
openssl req -new -out certs/server.csr -key certs/server.key -subj /CN=localhost
openssl x509 -req -days 3600 -in certs/server.csr -signkey certs/server.key -out certs/server.crt