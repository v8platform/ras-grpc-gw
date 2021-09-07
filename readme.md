# ras-grpc-gw

[![Release](https://img.shields.io/github/release/v8platform/ras-grpc-gw.svg?style=for-the-badge)](https://github.com/v8platform/ras-grpc-gw/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE)
[![Build status](https://img.shields.io/github/workflow/status/v8platform/ras-grpc-gw/goreleaser?style=for-the-badge)](https://github.com/v8platform/ras-grpc-gw/actions?workflow=goreleaser)
[![Codecov branch](https://img.shields.io/codecov/c/github/v8platform/ras-grpc-gw/master.svg?style=for-the-badge)](https://codecov.io/gh/v8platform/ras-grpc-gw)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](http://godoc.org/github.com/v8platform/ras-grpc-gw)
[![SayThanks.io](https://img.shields.io/badge/SayThanks.io-%E2%98%BC-1EAEDB.svg?style=for-the-badge)](https://saythanks.io/to/khorevaa)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)

> _*ЭТО ALFA Версия.*_ Проект в разработке и возможны ошибки в работе.

> _*КРАЙНЕ НЕ РЕКОМЕНДУЕТСЯ*_ использовать в промышленной эксплуатации

> Использование приложения на свой страх и риск, по принципу _*AS IS*_. Разработчики не несут ответственности за возможный ущерб

## Возможности

RAS GRPC gateway (`ras-grpc-gw`) - прокси сервер для службы RAS 1С Предприятие

Предполагаемое использование, разворачивание рядом со службой RAS в виде docker-контейнера или отдельной службы  

Особенности:

* Подключение к RAS только при первом запросе
* Если соединение было закрыто то делается `одна` попытка переподключения, при этом все точки сбрасываются

### Реализованная функциональность

* Сервис авторизации `AuthService` 
  * AuthenticateCluster - установка авторизации администратора на кластере
  * AuthenticateInfobase - установка авторизации на кластере для информационной базы
  * AuthenticateAgent - установка авторизации администратора на агенте
* Сервис кластера `ClustersService`
  * GetClusters - получение списка кластеров
  * GetClusterInfo - получение информации о кластере
  * RegCluster - регистрация нового кластера
  * UnregCluster - отмена регистрации на кластере
* Сервис информационных баз `InfobasesService`
  * GetShortInfobases - получение списка информационных баз на кластере
  * GetInfobaseSessions - получение списка сессий информационной базы 
* Сервис сессий кластера `SessionsService`
  * GetSessions - получение списка сессий кластера

## Как установить

* Установить из [`releases`](https://github.com/v8platform/ras-grpc-gw/releases/)
* Использовать готовый образ `docker`
  * `docker pull v8platform/ras-grpc-gw:latest`
  * `docker pull ghcr.io/v8platform/ras-grpc-gw:latest`

## Как использовать

### Работа с `endpoint` 

Для работы с точками обмена используются заголовки (метаданные) сообщений

Если в заголовке (или метаданные) не передать ключе `endpoint_id`  тогда операция выполниться в новой точке,
И в ответном сообщении будет указан `endpoint_id` - новой открытой точки, для дальнейшей работы с этом отрытой точкой обмена

Для `grpcurl` указание происходить через флаг `-H`. Например, `-H endpoint_id:1`

Для других клиентов надо передавать заголовок (метаданные)

### Запуск локально 

Запуск сервера для службы RAS локального кластера 1С по адресу `localhost`

```shell
ras-grpc-gw localhost:1545
```

### Запуск сервера в `docker` 


Запуск сервера для службы RAS кластера 1С по адресу `ras`


```shell
docker run -d --name ras-grpc-gw -p 3002:3002 v8platform/ras-grpc-gw ras:1545
```


### `CLI` клиент

#### Установка клиента `grpcurl`

Для использования клиента необходим скомпилированные файлы `proto` 
Уже собранный файл лежит [`./protos/protoset.bin`](./protos/protoset.bin)
``

*`Docker`*
```shell
# Download image
docker pull fullstorydev/grpcurl:latest
# Run the tool
docker run fullstorydev/grpcurl localhost:3002
```
`CLI`
```shell
go get github.com/fullstorydev/grpcurl/...
go install github.com/fullstorydev/grpcurl/cmd/grpcurl
```

#### Пример использования

*Получение списка кластеров*

`CLI`
```shell
grpcurl -protoset ./protos/protoset.bin -plaintext -H endpoint_id:1 -d '{}' localhost:3002 ras.service.api.v1.ClustersService/GetClusters
```
or

`Docker`
```shell
docker run -it -v $PWD/protos/protoset.bin:/protos/protoset.bin fullstorydev/grpcurl -protoset /protos/protoset.bin -plaintext -d '{}' localhost:3002 ras.service.api.v1.ClustersService/GetClusters
```

*Установка авторизации на кластере*
```shell
grpcurl -protoset ./protos/protoset.bin -plaintext -H endpoint_id:1 -d '{\"cluster_id\": \"e9261ed1-c9d0-40e5-8222-c7996493d507\"}' localhost:3002 ras.service.api.v1.AuthService/AuthenticateCluster
```

*Получение списка сессий кластера*
```shell
grpcurl -protoset ./protos/protoset.bin -plaintext -H endpoint_id:1 -d '{\"cluster_id\": \"e9261ed1-c9d0-40e5-8222-c7996493d507\"}' localhost:3002 ras.service.api.v1.SessionsService/GetSessions
```

*Получение списка информационных баз*
```shell
grpcurl -protoset ./protos/protoset.bin -plaintext -H endpoint_id:1 -d '{\"cluster_id\": \"e9261ed1-c9d0-40e5-8222-c7996493d507\"}' localhost:3002 ras.service.api.v1.InfobasesService/GetShortInfobases
```

## License

Лицензия [`LICENSE`](LICENSE)