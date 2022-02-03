DOCKER_TAG = secrets-manager:latest
MODULE = github.com/daniilty/secrets-manager
BRANCH ?=$(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_TIME ?= $(shell date +%FT%T%z)
LDFLAGS += -X ${MODULE}/internal/healthcheck.Branch=${BRANCH} -X ${MODULE}/internal/healthcheck.CommitHash=${COMMIT_HASH} -X ${MODULE}/internal/healthcheck.BuildTime=${BUILD_TIME}

build:
	go build -ldflags "${LDFLAGS}" -o service ${MODULE}/cmd/main
build_docker:
	docker build -t ${DOCKER_TAG} -f ./docker/Dockerfile .

