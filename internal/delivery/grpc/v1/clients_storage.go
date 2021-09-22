package v1

import (
	"context"
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"github.com/spf13/cast"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	context2 "github.com/v8platform/ras-grpc-gw/internal/context"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/anypb"
	"sync"
)

var _ clientv1.EndpointServiceImpl = (*Endpoint2)(nil)

type Client2 interface {
	GetEndpoint(ctx context.Context) (*Endpoint2, error)
}

type clientsStorage struct {
	idxClients   map[string]*ras_client.ClientConn
	idxEndpoints map[string]*Endpoint2

	mu sync.Mutex
}

type Endpoint2 struct {
	uuid   string
	client *ras_client.ClientConn
	impl   clientv1.EndpointServiceImpl

	context    map[endpointContext]interface{}
	clientUuid string
}

func (e Endpoint2) endpointContext() {}

type endpointContext = int

const (
// AgentAuth endpointContext = iota
// ClusterAuth
// InfobaseAuth
)

func (e *Endpoint2) SetEndpointContext(t endpointContext, req interface{}) {
	e.context[t] = req
}

func (e *Endpoint2) Request(ctx context.Context, req *clientv1.EndpointRequest) (resp *anypb.Any, err error) {

	defer func() {
		header := metadata.New(map[string]string{
			"X-Endpoint2": cast.ToString(e.uuid),
			"X-App":       cast.ToString(e.clientUuid),
		})
		_ = grpc.SendHeader(ctx, header)
	}()

	// Словил тут рекурсию
	return e.impl.Request(ctx, req)
}

func (c *clientsStorage) GetClient(ctx context.Context) (*ras_client.ClientConn, error) {
	panic("implement me")
}

func (c *clientsStorage) GetEndpoint(ctx context.Context) (*Endpoint2, error) {

	endpointId, ok := context2.EndpointFromContext(ctx)

	if !ok {
		return c.initEndpoint(ctx, endpointId)
	}

	endpoint := c.getEndpoint(endpointId)

	if endpoint == nil {
		return c.initEndpoint(ctx, endpointId)
	}

	return endpoint, nil
}

func (c *clientsStorage) initEndpoint(ctx context.Context, endpointId string) (*Endpoint2, error) {

	client, ok := context2.ClientFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no client setup, pls register client before")
	}

	conn := c.idxClients[client.UUID]

	if conn == nil {
		conn = ras_client.NewClientConn(client.Host)
		c.mu.Lock()
		c.idxClients[client.UUID] = conn
		c.mu.Unlock()
	}

	endpointImpl, err := conn.GetEndpoint(ctx, "")
	if err != nil {
		return nil, err
	}

	if len(endpointId) == 0 {
		endpointId = shortuuid.New()
	}

	// TODO init context

	endpoint := &Endpoint2{
		uuid:       endpointId,
		clientUuid: client.UUID,
		client:     conn,
		impl:       endpointImpl,
		context:    map[endpointContext]interface{}{},
	}

	c.mu.Lock()
	c.idxEndpoints[endpointId] = endpoint
	c.mu.Unlock()

	return endpoint, nil

}

func (c *clientsStorage) getEndpoint(endpointId string) *Endpoint2 {
	return c.idxEndpoints[endpointId]
}

func NewRasClientsStorage() Client {
	return &clientsStorage{
		idxEndpoints: make(map[string]*Endpoint2),
		idxClients:   make(map[string]*ras_client.ClientConn),
	}
}
