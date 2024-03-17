PROJECT?=github.com/Vladislav747/h-project
NAME?=h-project
VERSION?=0.0.1
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date '+%Y-%m-%dT%H:%M:%S')

clean:
	@rm bin/api

hproject_build: clean
	@go build \
		-ldflags "-w -s \
		-X ${PROJECT}/version.APIVersion=${VERSION} \
		-X ${PROJECT}/version.APIName=${NAME} \
		-X ${PROJECT}/version.Commit=${COMMIT} \
		-X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
		-o bin/api cmd/h-project/main.go

build:
	go build \
		-ldflags "-w -s \
		-X ${PROJECT}/version.APIVersion=${VERSION} \
		-X ${PROJECT}/version.APIName=${NAME} \
		-X ${PROJECT}/version.Commit=${COMMIT} \
		-X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
		-o bin/api cmd/h-project/main.go

run:
	. deployments/development/.env.local && go run ./cmd/h-project