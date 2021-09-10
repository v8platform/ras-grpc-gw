package v1

import (
	"context"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	context2 "github.com/v8platform/ras-grpc-gw/internal/context"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client"
	"google.golang.org/protobuf/types/known/anypb"
)


var _ clientv1.EndpointServiceImpl = (*Endpoint)(nil)

type ClientsStorage interface {

	GetClient(ctx context.Context) (*ras_client.ClientConn, error)
	GetEndpoint(ctx context.Context) (*Endpoint, error)


}

type clientsStorage struct {

	idxClients   map[string]*ras_client.ClientConn
	idxEndpoints map[string]*Endpoint
}

type Endpoint struct {

	uuid string
	client *ras_client.ClientConn
	impl clientv1.EndpointServiceImpl

	context map[endpointContext]interface{}

}

type endpointContext = int

const (
	AgentAuth endpointContext = iota
	ClusterAuth
	InfobaseAuth
)

func (e *Endpoint) SetEndpointContext(t endpointContext, req interface{}) {
	e.context[t] = req
}

func (e *Endpoint) Request(ctx context.Context, req *clientv1.EndpointRequest) (resp *anypb.Any, err error) {

	defer func() {

	}()

	return e.impl.Request(ctx, req)
}

func (c *clientsStorage) GetClient(ctx context.Context) (*ras_client.ClientConn, error){
	panic("implement me")
}

func (c *clientsStorage) GetEndpoint(ctx context.Context) (*Endpoint, error) {

	endpointId, ok := context2.EndpointFromContext(ctx)

	if !ok {
		return c.initEndpoint(ctx, endpointId)
	}

	endpoint := c.getEndpoint(endpointId)

	if endpoint == nil {
		return c.initEndpoint(ctx, endpointId)
	}

	return nil, nil
}

func (c *clientsStorage) initEndpoint(ctx context.Context, endpointId string) (*Endpoint, error) {

	return nil, nil

}

func (c *clientsStorage) getEndpoint(endpointId string) *Endpoint {
	return c.idxEndpoints[endpointId]
}


func NewRasClientsStorage() ClientsStorage {
	return &clientsStorage{}
}

