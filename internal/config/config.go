package config

import (
	"time"
)

const DefaultConfig = `
secure: false
log_level: info 
secret: "move_me_to_env RAS_GRPC_GW_SECRET"
grpc:
  port: 3000
  webui: 
    up: false
    path: /grpc-ui
gateway:
  up: false
  port: 3001
  docs:
   up: false
   path: /docs
database:
  pudge:
cache:
  ttl: 10m
  memory:
sha1_salt: salt_omm
ras:
  host: ras:1545
  version: "10.0"
  conn_pool_size: "15"
  conn_idle_timeout: "30m"
  conn_idle_check_timer: 5m
  conn_max_endpoints_count: 20
`

type Config struct {
	LogLevel string `config:"log_level"`

	// Secure признак использования защищенного подключения
	Secure bool `json:"secure" config:"secure"`
	// ? Добавить путь к ключу

	// Grpc конфигурация сервера GRPC
	Grpc GrpcConfig `json:"grpc"`

	// Gateway конфигурация сервера REST API
	Gateway GatewayConfig `json:"gateway"`

	// Секрет используется для подписи токенов и
	// Шифрования хеша паролей а базе данных
	Secret string `json:"secret"`

	// Настройка базы данных
	Database DatabaseConfig `json:"database"`

	// Настройка механизма кеширвания
	Cache CacheConfig `json:"cache"`

	// Соль для добавления в хеши
	SHA1Salt string `json:"sha1_salt" yaml:"sha1_salt" config:"sha1_salt"`

	// Настройка подключения к службе RAS
	RAS RASConfig `json:"ras"`
}

type RASConfig struct {

	// Адрес службы RAS сервера 1С
	Host string `json:"host"`
	// Версия подключенной клиента
	Version string `json:"version"`
	// Размер пула соединений клиента
	ConnPoolSize int `json:"conn_pool_size" yaml:"conn_pool_size" config:"conn_pool_size"`
	// Таймаут простоя соединения
	ConnIdleTimeout time.Duration `json:"conn_idle_timeout" yaml:"conn_idle_timeout" config:"conn_idle_timeout"`
	// Период проверки простоя соединения
	ConnIdleCheckTimer time.Duration `json:"conn_idle_check_timer" yaml:"conn_idle_check_timer" config:"conn_idle_check_timer"`
	// Максимальное число открытых точек обмена на подключение
	ConnMaxEndpointsCount int `json:"conn_max_endpoints_count" yaml:"conn_max_endpoints_count" config:"conn_max_endpoints_count"`

	// Настройка точки обмена по умолчанию
	DefaultEndpointConfig EndpointConfig `json:"default_endpoint_config" config:"default_endpoint_config"`
	// Известная авторизация на информационных базах
	DefaultAuthorizations map[string]UserAuth `json:"default_authorizations" config:"default_authorizations"`
}

// Настройка доступов к по умолчанию для точки обмена
type EndpointConfig struct {
	// Доступ к серверу 1С
	Agent UserAuth `json:"agent"`
	// Доступ к локальному кластеру сервера 1С
	Cluster UserAuth `json:"cluster"`
	// Доступ к информационной базе
	Infobase UserAuth `json:"infobase"`
}

type UserAuth struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type DatabaseConfig struct {
	Engine Namespace `config:",inline,replace"`
}

type CacheConfig struct {
	TTL    time.Duration `config:"ttl"`
	Engine Namespace     `config:",inline,replace"`
}

type GrpcConfig struct {
	Port  int32
	Webui GrpcUIConfig
}

type GrpcUIConfig struct {
	Up   bool
	Path string
}

type GatewayConfig struct {
	Up   bool
	Port int32
	Docs SwaggerConfig
}

type SwaggerConfig struct {
	Up   bool
	Path string
}
