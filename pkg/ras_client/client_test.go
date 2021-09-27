package client

import (
	"context"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
	"testing"
)

func TestNewClient(t *testing.T) {

	host := "srv-uk-app-16:1545"

	type init struct {
		addr string
		opts []GlobalOption
	}

	type invoke struct {
		invoke       clientv1.InvokeHandler
		needEndpoint bool
		req          interface{}
		reply        interface{}
		ctx          context.Context
		interceptor  Interceptor
	}

	tests := []struct {
		name    string
		args    init
		invoke  invoke
		wantErr bool
	}{
		{
			name: "ConnectHandler",
			args: init{
				addr: host,
				opts: nil,
			},
			invoke: invoke{
				clientv1.ConnectHandler,
				false,
				&protocolv1.ConnectMessage{},
				&protocolv1.ConnectMessageAck{},
				context.Background(),
				nil,
			},
		},
		{
			name: "GetClustersRequest",
			args: init{
				addr: host,
				opts: nil,
			},
			invoke: invoke{
				clientv1.GetClustersHandler,
				true,
				&messagesv1.GetClustersRequest{},
				&messagesv1.GetClustersResponse{},
				context.Background(),
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.args.addr, tt.args.opts...)
			var (
				resp interface{}
				err  error
			)
			if resp, err = client.Invoke(tt.invoke.ctx, tt.invoke.needEndpoint, tt.invoke.req, tt.invoke.invoke, tt.invoke.interceptor); (err != nil) != tt.wantErr {
				t.Errorf("Invoke() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("stats %v", client.Stats())
			t.Logf("resp %s", resp)

		})
	}
}
