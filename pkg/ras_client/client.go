package client

import (
	"context"
	"fmt"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type DialFunc func(addr string) (net.Conn, error)

type Client interface {
	Host() string

	Connect(opts ...ConnectOption) error
	SetConn(conn net.Conn, opts ...SetConnOption) error
	GetConn() net.Conn
	Close() error
	Closed() bool
	UsedAt() time.Time
	SetUsedAt(tm time.Time)

	NewEndpoint(opts ...EndpointOption) (*Endpoint, error)
	CloseEndpoint(endpoint *Endpoint) error
	Endpoints() []*Endpoint

	clientv1.Client
	clientv1.ClustersService
	clientv1.AdminService
	clientv1.AuthService
	clientv1.ConnectionsService
	clientv1.InfobasesService
	clientv1.LocksService
	clientv1.SessionsService
}

func NewClient(addr string, opts ...ClientOption) Client {
	c := newClient(addr, opts...)

	return c
}

func NewClientConn(conn net.Conn, opts ...ClientOption) Client {

	c := newClient(conn.RemoteAddr().String(), opts...)
	err := c.setConn(conn)
	if err != nil {
		panic(err)
	}

	return c
}

type EndpointConfig struct {
	Options []RequestOption

	DefaultAgentAuth    Auth
	DefaultClusterAuth  Auth
	DefaultInfobaseAuth Auth

	SaveAuthRequests bool
	Auths            map[string]Auth
}

type client struct {
	addr string
	dial DialFunc
	cc   net.Conn

	_usedAt    uint32 // atomic
	_closed    uint32 // atomic
	_connected uint32 // atomic

	useAutoReconnect bool

	connectOptions  map[interface{}]ConnectOption
	endpointOptions map[interface{}]EndpointOption

	mu sync.Mutex

	endpoints      map[string]*Endpoint
	endpointConn   map[string]net.Conn
	endpointConfig map[string]*EndpointConfig

	clientService clientv1.ClientService

	clientv1.ClustersService
	clientv1.AdminService
	clientv1.AuthService
	clientv1.ConnectionsService
	clientv1.InfobasesService
	clientv1.LocksService
	clientv1.SessionsService
}

func (c *client) CloseEndpoint(endpoint *Endpoint) error {
	_, err := c.clientService.EndpointClose(
		context.Background(),
		&protocolv1.EndpointClose{EndpointId: endpoint.GetId()},
	)

	return err

}

func (c *client) GetEndpoint(ctx context.Context) (clientv1.Endpoint, error) {
	uuid := EndpointFromContext(ctx)
	endpoint, ok := c.endpoints[uuid]

	if ok {
		return endpoint, nil
	}

	return c.NewEndpoint()

}

func (c *client) Request(ctx context.Context, handler clientv1.RequestHandler, opts ...interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed() {
		return fmt.Errorf("client is closed")
	}
	err := handler(ctx, c.cc)
	if err != nil {
		return err
	}
	return err
}

func (c *client) UsedAt() time.Time {
	unix := atomic.LoadUint32(&c._usedAt)
	return time.Unix(int64(unix), 0)
}

func (c *client) SetUsedAt(tm time.Time) {
	atomic.StoreUint32(&c._usedAt, uint32(tm.Unix()))
}

func (c *client) Close() error {

	if !atomic.CompareAndSwapUint32(&c._closed, 0, 1) {
		return nil
	}

	// ctx := context.Background()
	// var err error

	if atomic.CompareAndSwapUint32(&c._connected, 0, 1) {
		// err = c.disconnect(ctx, &protocolv1.DisconnectMessage{})
		// if err != nil {
		// 	return err
		// }
	}

	if c.closed() {
		return nil
	}

	return c.cc.Close()
}

func (c *client) Host() string {
	return c.addr
}

func (c *client) Connect(opts ...ConnectOption) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var options []Option
	for _, opt := range opts {
		options = append(options, opt)
	}

	if err := c.initConnect(options...); err != nil {
		return err
	}

	return nil
}

func (c *client) SetConn(conn net.Conn, opts ...SetConnOption) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.setConn(conn, opts...)
}

func (c *client) GetConn() net.Conn {
	return c.cc
}

func (c *client) Closed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.closed()
}

func (c *client) NewEndpoint(opts ...EndpointOption) (*Endpoint, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return &Endpoint{}, nil
}

func (c *client) Endpoints() []*Endpoint {
	c.mu.Lock()
	defer c.mu.Unlock()

	var endpoints []*Endpoint
	for _, endpoint := range c.endpoints {
		endpoints = append(endpoints, endpoint)
	}

	return endpoints
}

func (c *client) applyOptions(opts ...ClientOption) {

	for _, opt := range opts {
		switch opt.Ident() {
		case dialFuncIdent{}:
			c.dial = opt.Value().(DialFunc)
			continue
		case connIdent{}:
			c.cc = opt.Value().(net.Conn)
			continue
		case reconnectIdent{}:
			c.useAutoReconnect = opt.Value().(bool)
			continue
		}
		switch opt.(type) {
		case EndpointOption:
			c.endpointOptions[opt.Ident()] = opt.(EndpointOption)
		case ConnectOption:
			c.connectOptions[opt.Ident()] = opt.(ConnectOption)
		}
	}
}

func (c *client) setConn(conn net.Conn, opts ...SetConnOption) error {

	var (
		restoreEndpoints bool
		restoreConnect   bool
	)

	for _, opt := range opts {
		switch opt.Ident() {
		case restoreConnectIdent{}:
			restoreEndpoints = opt.Value().(bool)
		case restoreEndpointsIdent{}:
			restoreEndpoints = opt.Value().(bool)
		}
	}

	_ = c.cc.Close()

	atomic.StoreUint32(&c._connected, 0)
	atomic.StoreUint32(&c._closed, 0)

	c.cc = conn

	if restoreConnect || restoreEndpoints {
		if err := c.initConnect(); err != nil {
			return err
		}
	}

	if restoreEndpoints {
		if err := c.restoreEndpoints(); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) restoreEndpoints(opts ...Option) error {

	return nil
}

func (c *client) initConnect(opts ...Option) error {

	return nil
}

func (c *client) connected() bool {
	return atomic.LoadUint32(&c._connected) == 1
}

func (c *client) closed() bool {

	if atomic.LoadUint32(&c._closed) == 1 || c.cc == nil {
		return true
	}
	_ = c.cc.SetReadDeadline(time.Now())
	_, err := c.cc.Read(make([]byte, 0))
	var zero time.Time
	_ = c.cc.SetReadDeadline(zero)

	if err == nil {
		return false
	}

	netErr, _ := err.(net.Error)
	if err != io.EOF && !netErr.Timeout() {
		atomic.StoreUint32(&c._closed, 1)
		return true
	}
	return false
}

func newClient(addr string, opts ...ClientOption) *client {
	c := &client{
		addr:            addr,
		endpoints:       map[string]*Endpoint{},
		endpointConfig:  map[string]*EndpointConfig{},
		endpointOptions: map[interface{}]EndpointOption{},
		connectOptions:  map[interface{}]ConnectOption{},
		endpointConn:    map[string]net.Conn{},
		dial:            defaultDial,
	}

	c.clientService = clientv1.NewClientService(c)
	c.ClustersService = clientv1.NewClustersService(c)
	c.AdminService = clientv1.NewAdminService(c)
	c.AuthService = clientv1.NewAuthService(c)
	c.ConnectionsService = clientv1.NewConnectionsService(c)
	c.InfobasesService = clientv1.NewInfobasesService(c)
	c.LocksService = clientv1.NewLocksService(c)
	c.SessionsService = clientv1.NewSessionsService(c)

	c.applyOptions(opts...)
	return c
}

func defaultDial(addr string) (net.Conn, error) {

	return net.Dial("tcp", addr)

}
