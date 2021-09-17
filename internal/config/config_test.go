package config

import (
	"github.com/elastic/go-ucfg"
	"reflect"
	"testing"
	"time"
)

func TestNewConfigFrom(t *testing.T) {
	type args struct {
		from interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			"default_config_test",
			args{from: DefaultConfig},
			Config{
				"info",
				false,
				GrpcConfig{
					3000,
					GrpcUIConfig{
						false,
						"/grpc-ui",
					},
				},
				GatewayConfig{
					false,
					3001,
					SwaggerConfig{
						false,
						"/docs",
					},
				},
				"move_me_to_env RAS_GRPC_GW_SECRET",
				DatabaseConfig{Engine: Namespace{name: "pudge", config: &ucfg.Config{}}},
				CacheConfig{TTL: 10 * time.Minute,
					Engine: Namespace{name: "memory", config: &ucfg.Config{}},
				},
				"salt_omm",
				RASConfig{
					Host:                  "ras:1545",
					Version:               "10.0",
					ConnPoolSize:          15,
					ConnIdleTimeout:       30 * time.Minute,
					ConnIdleCheckTimer:    5 * time.Minute,
					ConnMaxEndpointsCount: 20,
					DefaultEndpointConfig: EndpointConfig{},
					DefaultAuthorizations: nil,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfigFrom(tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfigFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var config Config
			err = got.Unpack(&config)

			// delete pointer
			config.Database.Engine.config = nil
			tt.want.Database.Engine.config = nil
			config.Cache.Engine.config = nil
			tt.want.Cache.Engine.config = nil

			if !reflect.DeepEqual(config, tt.want) {
				t.Errorf("NewConfigFrom() got = %v, want %v", config, tt.want)
			}
		})
	}
}

func TestNewConfigWithYAML(t *testing.T) {
	type args struct {
		in     []byte
		source string
	}
	tests := []struct {
		name    string
		args    args
		want    *ucfg.Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfigWithYAML(tt.args.in, tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfigWithYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigWithYAML() got = %v, want %v", got, tt.want)
			}
		})
	}
}
