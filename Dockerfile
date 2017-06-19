FROM golang:1.8.3-alpine

WORKDIR /app

COPY bin/keyvault /usr/local/bin/keyvault
