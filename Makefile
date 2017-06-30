GO ?= go
GOX ?= gox
DOCKER ?= docker
OS=$(shell uname)
CURRENTDIR=$(shell pwd)
DOCKER_IMAGE=nexus-registry.cn133.azure.net/tools/keyvault-config
VERSION=0.2

ifeq ($(OS),Darwin)
 	GLIDE_INSTALL=brew install glide
else
  	GLIDE_INSTALL=curl https://glide.sh/get | sh
endif

all: build

gox:
	${GO} get github.com/mitchellh/gox

glide:
	${GLIDE_INSTALL}

vendor: glide
	cd ${CURRENTDIR}/src/keyvault; \
	glide up; \
	cd -

build: vendor gox
	GOPATH=${CURRENTDIR}; \
	${GOX} -verbose -osarch="linux/amd64" -output="${CURRENTDIR}/bin/keyvault"

docker: build
	DOCKER build -t ${DOCKER_IMAGE}:${VERSION} .

docker-push:
	DOCKER push ${DOCKER_IMAGE}:${VERSION}
