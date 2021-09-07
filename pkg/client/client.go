package client

import (
	"context"
	"github.com/spf13/cast"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var defaultVersion = "10.0"

var _ clientv1.ClientImpl = (*ClientConn)(nil)
var _ clientv1.ClientServiceImpl = (*ClientConn)(nil)

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

func (c *ClientConn) GetEndpoint(ctx context.Context) (clientv1.EndpointServiceImpl, error) {

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.DataLoss, "Client: failed to get metadata")
	}

	if t, ok := md["endpoint_id"]; ok {

		for _, e := range t {
			if endpoint, ok := c.getEndpoint(e); ok {
				return clientv1.NewEndpointService(c, endpoint), nil

			}
		}
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
