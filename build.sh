#!/bin/sh

set -o errexit

VERSION=0.1

gox -osarch=linux/amd64 -output=bin/keyvault
docker build -t keyvault-config:$VERSION .
