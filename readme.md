# ras-grpc-gw

[![Release](https://img.shields.io/github/release/khorevaa/ras-grpc-gw.svg?style=for-the-badge)](https://github.com/khorevaa/ras-grpc-gw/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](/LICENSE.md)
[![Build status](https://img.shields.io/github/workflow/status/khorevaa/ras-grpc-gw/build?style=for-the-badge)](https://github.com/khorevaa/ras-grpc-gw/actions?workflow=releaser)
[![Codecov branch](https://img.shields.io/codecov/c/github/khorevaa/ras-grpc-gw/master.svg?style=for-the-badge)](https://codecov.io/gh/khorevaa/ras-grpc-gw)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](http://godoc.org/github.com/khorevaa/ras-grpc-gw)
[![SayThanks.io](https://img.shields.io/badge/SayThanks.io-%E2%98%BC-1EAEDB.svg?style=for-the-badge)](https://saythanks.io/to/khorevaa)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)

## Features


## Installation



## How to use


### Client

#### Install
*Docker*
```shell
# Download image
docker pull fullstorydev/grpcurl:latest
# Run the tool
docker run fullstorydev/grpcurl localhost:3002 list
```
*CLI*
```shell
go get github.com/fullstorydev/grpcurl/...
go install github.com/fullstorydev/grpcurl/cmd/grpcurl
```

#### Usage

*Get clusters*
```shell
set PROTO_DECR=../protos/decr.bin
grpcurl -protoset $PROTO_DECR -plaintext -d '{}' localhost:3002 ras.service.api.v1.RASService/GetClusters
```

### Docker

### Github Release

## License