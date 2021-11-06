.PHONY: all clean install uninstall

DOCKER_COMPOSE=./tests/docker/docker-compose.yaml
GO_MODULES=github.com/bufbuild/buf/cmd/buf \
           		github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking \
           		github.com/bufbuild/buf/cmd/protoc-gen-buf-lint \
           		github.com/v8platform/protoc-gen-go-ras \
           		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
           		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
           		google.golang.org/protobuf/cmd/protoc-gen-go \
           		google.golang.org/grpc/cmd/protoc-gen-go-grpc

include .secret.env

login:
	echo "${GITHUB_TOKEN}" | docker login ghcr.io --username ${GITHUB_USER} --password-stdin

all: hello

deps:
	go get ${GO_MODULES}
	go install ${GO_MODULES}
proto-gen:
	cd ./rasgrpcgwapis && buf mod update
	buf generate

dc-up:
	docker-compose -f ${DOCKER_COMPOSE} up -d

dc-down:
	docker-compose -f ${DOCKER_COMPOSE} down

generate:
	go generate ./...

test:
	go test ./...