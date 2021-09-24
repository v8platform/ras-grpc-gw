package client

import (
	"github.com/elastic/go-ucfg"
	"time"
)

type Config struct {

	// Адрес службы RAS сервера 1С
	Host string `json:"host" config:"host" json:"host,omitempty"`

	// Версия подключенной клиента
	Version string `json:"version" config:"version" json:"version,omitempty"`

	ChannelsConfig ChannelPoolConfig `config:"channels" json:"channels"`

	ConnectConfig map[string][]*ucfg.Config `config:"connect_config"`

	// Настройка точки обмена по умолчанию
	DefaultEndpointConfig EndpointConfig `json:"default_endpoint_config" config:"default_endpoint_config" config:"default_endpoint_config" json:"default_endpoint_config"`
}

type ChannelPoolConfig struct {
	PoolSize           int           `config:"pool_size" json:"pool_size,omitempty"`
	PoolTimeout        time.Duration `config:"pool_timeout" json:"pool_timeout,omitempty"`
	IdleTimeout        time.Duration `config:"idle_timeout" json:"idle_timeout,omitempty"`
	MaxChannelAge      time.Duration `config:"max_channel_age" json:"max_channel_age,omitempty"`
	IdleCheckFrequency time.Duration `config:"idle_check_frequency" json:"idle_check_frequency,omitempty"`
	MinIdleChannels    int           `config:"min_idle_channels" json:"min_idle_channels,omitempty"`
}

type ConnectMessageConfig struct {
	Params map[string]interface{} `config:"params"`
}

type NegotiateMessageConfig struct {
	Magic    int32 `config:"magic"`
	Protocol int32 `config:"protocol"`
	Version  int32 `config:"version"`
}

// EndpointConfig Настройка доступов к по умолчанию для точки обмена
type EndpointConfig struct {

	// Доступ к серверу 1С
	DefaultAgentAuth Auth `config:"default_agent_auth" json:"default_agent_auth"`
	// Доступ к локальному кластеру сервера 1С
	DefaultClusterAuth Auth `config:"default_cluster_auth" json:"default_cluster_auth"`
	// Доступ к информационной базе
	DefaultInfobaseAuth Auth `config:"default_infobase_auth" json:"default_infobase_auth"`

	// Сохранение запросов на авторитизацию
	SaveAuthRequests bool `config:"save_auth_requests" json:"save_auth_requests,omitempty"`
	// Известная авторизация на информационных базах
	Auths map[string]Auth `config:"auths" json:"auths,omitempty"`
}
