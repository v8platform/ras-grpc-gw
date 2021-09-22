package ras_client

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

/*
Описание работы клиента
1. Получаем точку обмена
 - Есть ID
 - Нет ID



Состав клиента
1. Пул соединений - варианты: размерный или одиночный
2. Индекс точек работы и соединений - map[string]*conn

Состав Endpoint2
+ ID     string
+ Config EndpointConfig
+ usedAt uint32 // atomic

Состав conn
... поля пула
+ idxEndpoints [string]*protocolv1.Endpoint2



*/

var defaultVersion = "10.0"

var _ clientv1.ClientImpl = (*ClientConn)(nil)
var _ clientv1.ClientServiceImpl = (*ClientConn)(nil)

type Poller interface {
	Get() net.Conn
	Put(conn net.Conn)
}

type Client struct {
	pool Poller

	idxEndpoints map[string]*Endpoint
}

type EndpointConfig struct {
	negotiate *protocolv1.NegotiateMessage
	connect   *protocolv1.ConnectMessage
}

type Endpoint struct {
	UUID string

	EndpointInfo

	conn   net.Conn
	config EndpointConfig
}

func (e *Endpoint) GetVersion() int32 {
	panic("implement me")
}

func (e *Endpoint) GetId() int32 {
	panic("implement me")
}

func (e *Endpoint) GetService() string {
	panic("implement me")
}

func (e *Endpoint) GetFormat() int32 {
	panic("implement me")
}

type EndpointInfo struct {
	service string
	version int32
	id      int32
	format  int32
}

func (e *Endpoint) stale() bool {
	return e.conn == nil || e.conn.Closed()
}

func (e *Endpoint) init() error {
	err := e.config.negotiate.Formatter(e.conn, 0)
	if err != nil {
		return
	}

	_, err = c.connect(ctx, c.ConnectMessage)
	if err != nil {
		return
	}
}

func (e *Endpoint) setConn(conn net.Conn) error {
	e.conn = conn

	if err := e.init(); err != nil {
		return err
	}

}

func (x *Endpoint) NewMessage(data interface{}) (*protocolv1.EndpointMessage, error) {
	switch typed := data.(type) {
	case io.Reader:
		packet, err := protocolv1.NewPacket(data)
		if err != nil {
			return nil, err
		}
		var message protocolv1.EndpointMessage
		if err := packet.Unpack(&message); err != nil {
			return nil, err
		}
		return &message, nil
	case protocolv1.EndpointMessageFormatter:
		return protocolv1.NewEndpointMessage(x, typed)
	default:
		return nil, fmt.Errorf("unknown type <%T> to create new message", typed)
	}
}

func (x *Endpoint) UnpackMessage(data interface{}, into protocolv1.EndpointMessageParser) error {
	switch typed := data.(type) {
	case io.Reader:
		packet, err := protocolv1.NewPacket(data)
		if err != nil {
			return err
		}
		var message protocolv1.EndpointMessage
		if err := packet.Unpack(&message); err != nil {
			return err
		}
		return message.Unpack(x, into)
	case protocolv1.Packet:
		var message protocolv1.EndpointMessage
		if err := typed.Unpack(&message); err != nil {
			return err
		}
		return message.Unpack(x, into)
	case *protocolv1.EndpointMessage:
		return typed.Unpack(x, into)
	default:
		return fmt.Errorf("unknown type <%T> to create unpack message", typed)
	}
}

func (x *Endpoint) Request(ctx context.Context, req *EndpointRequest) (*anypb.Any, error) {
	message, err := anypb.UnmarshalNew(req.GetRequest(), proto.UnmarshalOptions{})
	if err != nil {
		return nil, err
	}

	reqMessage, err := protocolv1.NewEndpointMessage(x, message)
	if err != nil {
		return nil, err
	}

	respMessage, err := x.client.EndpointMessage(ctx, reqMessage)
	if err != nil {
		return nil, err
	}

	respProtoMessage, err := anypb.UnmarshalNew(req.GetRespond(), proto.UnmarshalOptions{})
	if err != nil {
		return nil, err
	}

	if _, ok := respProtoMessage.(*emptypb.Empty); ok {
		if err := x.UnpackMessage(respMessage, nil); err != nil {
			return nil, err
		}
		return anypb.New(respProtoMessage)
	}

	messageParser, ok := respProtoMessage.(v1.EndpointMessageParser)
	if !ok {
		return nil, fmt.Errorf("not parser interface")
	}
	if err := x.UnpackMessage(respMessage, messageParser); err != nil {
		return nil, err
	}
	return anypb.New(respProtoMessage)
}

func NewClientConn(host string, opts ...Options) *ClientConn {

	opt := defaultClientOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	client := &ClientConn{
		host:      host,
		Options:   opt,
		mu:        &sync.Mutex{},
		connMu:    &sync.Mutex{},
		version:   defaultVersion,
		endpoints: &sync.Map{},
	}

	client.ClientServiceImpl = clientv1.NewClientService(client)

	return client
}

type connPool struct {
	opt           *Options
	dialErrorsNum uint32 // atomic
	_closed       uint32 // atomic

	lastDialErrorMu sync.RWMutex
	lastDialError   error

	queue        chan struct{}
	poolSize     int
	idleConnsLen int

	connsMu   sync.Mutex
	conns     []*Conn
	idleConns IdleConns
}

type ClientConn struct {
	Options

	host       string
	conn       net.Conn
	usedAt     uint32 // atomic
	_closed    uint32 // atomic
	_connected uint32 // atomic
	_locked    uint32 // atomic

	stats     Stats
	mu        *sync.Mutex // Блокировка всего клиента
	connMu    *sync.Mutex // Блокировка только соединения
	endpoints *sync.Map
	version   string

	clientv1.ClientServiceImpl
}

type Stats struct {
	Recv  uint32
	Send  uint32
	Wrong uint32
	Ping  uint32
}

type Options struct {
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
	Timeout            time.Duration
	NegotiateMessage   *protocolv1.NegotiateMessage
	ConnectMessage     *protocolv1.ConnectMessage
	OpenEndpoint       *protocolv1.EndpointOpen
}

var defaultClientOptions = Options{
	IdleTimeout:        30 * time.Minute,
	IdleCheckFrequency: 5 * time.Minute,
	Timeout:            3 * time.Second,
	NegotiateMessage:   protocolv1.NewNegotiateMessage(),
	ConnectMessage:     &protocolv1.ConnectMessage{},
	OpenEndpoint: &protocolv1.EndpointOpen{
		Service: "v8.service.Admin.Cluster",
		Version: defaultVersion,
	},
}

func (c *ClientConn) GetEndpoint(ctx context.Context, endpointID string) (clientv1.EndpointServiceImpl, error) {

	if endpoint, ok := c.getEndpoint(endpointID); ok {
		return clientv1.NewEndpointService(c, endpoint), nil
	}

	endpoint, err := c.turnEndpoint(ctx)
	if err != nil {
		return nil, err
	}

	return clientv1.NewEndpointService(c, endpoint), nil

}

func (c *ClientConn) getEndpoint(id string) (*protocolv1.Endpoint, bool) {

	if val, ok := c.endpoints.Load(id); ok {
		return val.(*protocolv1.Endpoint), ok
	}
	return nil, false
}

func (c *ClientConn) addEndpoint(endpoint *protocolv1.Endpoint) {

	id := cast.ToString(endpoint.GetId())
	log.Println(id)
	c.endpoints.Store(id, endpoint)

}

func (c *ClientConn) turnEndpoint(ctx context.Context) (*protocolv1.Endpoint, error) {

	EndpointOpenAck, err := c.EndpointOpen(ctx, &protocolv1.EndpointOpen{
		Service: "v8.service.Admin.Cluster",
		Version: c.version,
	})

	if err != nil {
		var version string

		if version = clientv1.DetectSupportedVersion(err); len(version) == 0 {
			return nil, err
		}
		if EndpointOpenAck, err = c.EndpointOpen(ctx, &protocolv1.EndpointOpen{
			Service: "v8.service.Admin.Cluster",
			Version: version,
		}); err != nil {
			return nil, err
		}

		c.version = version
	}

	end, err := c.NewEndpoint(ctx, EndpointOpenAck)
	if err != nil {
		return nil, err
	}

	c.addEndpoint(end)

	return end, nil
}

func (c *ClientConn) EndpointMessage(ctx context.Context, req *protocolv1.EndpointMessage) (*protocolv1.EndpointMessage, error) {
	defer func() {

	}()

	return c.ClientServiceImpl.EndpointMessage(ctx, req)

}

func (c *ClientConn) Read(p []byte) (n int, err error) {

	if c.closed() {
		if err := c.reconnect(); err != nil {
			return 0, err
		}
	}

	err = c.conn.SetReadDeadline(time.Now().Add(c.Timeout))
	if err != nil {
		return 0, err
	}
	defer func() {
		c.SetUsedAt(time.Now())
	}()

	return c.conn.Read(p)

}

func (c *ClientConn) Write(p []byte) (n int, err error) {

	if c.closed() {
		if err := c.reconnect(); err != nil {
			return 0, err
		}
	}

	err = c.conn.SetWriteDeadline(time.Now().Add(c.Timeout))
	if err != nil {
		return 0, err
	}
	defer func() {
		c.SetUsedAt(time.Now())
	}()

	return c.conn.Write(p)
}

func (c *ClientConn) UsedAt() time.Time {
	unix := atomic.LoadUint32(&c.usedAt)
	return time.Unix(int64(unix), 0)
}

func (c *ClientConn) SetUsedAt(tm time.Time) {
	atomic.StoreUint32(&c.usedAt, uint32(tm.Unix()))
}

func (c *ClientConn) Close() error {

	if !atomic.CompareAndSwapUint32(&c._closed, 0, 1) {
		return nil
	}

	ctx := context.Background()
	var err error
	// c.endpoints.Range(func(key, value interface{}) bool {
	//
	// 	// err = c.request(ctx, &protocolv1.EndpointClose{EndpointId: key.(int32)}, nil)
	// 	// if err != nil {
	// 	// 	return false
	// 	// }
	// 	//
	// 	return true
	// })

	if atomic.CompareAndSwapUint32(&c._connected, 0, 1) {

		_, err = c.Disconnect(ctx, &protocolv1.DisconnectMessage{})
		if err != nil {
			return err
		}

	}

	if c.closed() {
		return nil
	}

	return c.conn.Close()
}

func (c *ClientConn) Lock() {
	c.connMu.Lock()
}

func (c *ClientConn) Unlock() {
	c.connMu.Unlock()
}

func (c *ClientConn) reconnect() (err error) {

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected() {
		return nil
	}

	c.endpoints = &sync.Map{}

	ctx := context.Background()

	err = c.populateConn()
	if err != nil {
		return err
	}

	err = c.NegotiateMessage.Formatter(c.conn, 0)
	if err != nil {
		return
	}

	_, err = c.connect(ctx, c.ConnectMessage)
	if err != nil {
		return
	}

	atomic.StoreUint32(&c._connected, 1)

	return err

}

func (x *ClientConn) connect(ctx context.Context, req *protocolv1.ConnectMessage) (*protocolv1.ConnectMessageAck, error) {

	// Check context
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	packet, err := protocolv1.NewPacket(req)
	if err != nil {
		return nil, err
	}
	if _, err := packet.WriteTo(x.conn); err != nil {
		return nil, err
	}
	ackPacket, err := protocolv1.NewPacket(x.conn)
	if err != nil {
		return nil, err
	}
	resp := new(protocolv1.ConnectMessageAck)
	return resp, ackPacket.Unpack(resp)
}

func (c *ClientConn) connected() bool {
	return atomic.LoadUint32(&c._connected) == 1
}

func (c *ClientConn) populateConn() (err error) {

	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *ClientConn) closed() bool {

	if atomic.LoadUint32(&c._closed) == 1 || c.conn == nil {
		return true
	}
	_ = c.conn.SetReadDeadline(time.Now())
	_, err := c.conn.Read(make([]byte, 0))
	var zero time.Time
	_ = c.conn.SetReadDeadline(zero)

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
