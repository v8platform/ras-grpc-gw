package v1

import (
	"bytes"
	"context"
	"fmt"
	"github.com/lestrrat-go/option"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client"

	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Option interface {
	option.Interface
}

type ConnectOption interface {
	Option
	ClientOption
	connectOption()
}

type connectOption struct {
	Option
}

func (*connectOption) connectOption() {}
func (*connectOption) clientOption()  {}

func newConnectOption(n interface{}, v interface{}) ConnectOption {
	return &connectOption{option.New(n, v)}
}

type negotiateMessageIdent struct{}
type connectMessageIdent struct{}
type timeoutIdent struct{}

func NegotiateMessage(m *protocolv1.NegotiateMessage) ConnectOption {
	return newConnectOption(negotiateMessageIdent{}, m)
}

func ConnectMessage(m *protocolv1.ConnectMessage) ConnectOption {
	return newConnectOption(connectMessageIdent{}, m)
}

func Timeout(duration time.Duration) ConnectOption {
	return newConnectOption(timeoutIdent{}, duration)
}

type RequestOption interface {
	Option
	EndpointOption
	requestOption()
}

type requestOption struct {
	Option
}

func (*requestOption) requestOption()  {}
func (*requestOption) endpointOption() {}
func (*requestOption) clientOption()   {}

type endpointIdent struct{}

func newRequestOption(n interface{}, v interface{}) RequestOption {
	return &requestOption{option.New(n, v)}
}
func Endpoint(id, version, format int32) RequestOption {
	return newRequestOption(endpointIdent{}, endpointData{id, version, format})
}

type endpointData struct {
	id, version, format int32
}

type EndpointOption interface {
	Option
	ClientOption
	endpointOption()
}

type endpointOption struct {
	Option
}

func (*endpointOption) endpointOption() {}
func (*endpointOption) clientOption()   {}

func newEndpointOption(n interface{}, v interface{}) EndpointOption {
	return &endpointOption{option.New(n, v)}
}

type versionIdent struct{}
type mustVersionIdent struct{}
type serviceIdent struct{}
type uuidIdent struct{}
type clusterAuthIdent struct{}
type agentAuthIdent struct{}
type infobaseAuthIdent struct{}
type endpointConfigIdent struct{}
type saveAuthIdent struct{}

type Auth struct {
	user     string
	password string
}

func Version(version int32) EndpointOption {
	return newEndpointOption(versionIdent{}, version)
}
func AutosaveAuth(save bool) EndpointOption {
	return newEndpointOption(saveAuthIdent{}, save)
}

func Config(config EndpointConfig) EndpointOption {
	return newEndpointOption(endpointConfigIdent{}, config)
}

func MustVersion(version int32) EndpointOption {
	return newEndpointOption(mustVersionIdent{}, version)
}

func Service(service string) EndpointOption {
	return newEndpointOption(serviceIdent{}, service)
}

func UUID(uuid string) EndpointOption {
	return newEndpointOption(uuidIdent{}, uuid)
}

func DefaultClusterAuth(user, password string) EndpointOption {
	return newEndpointOption(clusterAuthIdent{}, Auth{user, password})
}

func DefaultAgentAuth(user, password string) EndpointOption {
	return newEndpointOption(agentAuthIdent{}, Auth{user, password})
}

func DefaultInfobaseAuth(user, password string) EndpointOption {
	return newEndpointOption(infobaseAuthIdent{}, Auth{user, password})
}

type SetConnOption interface {
	Option
	setConnOption()
}

type restoreConnectIdent struct{}
type restoreEndpointsIdent struct{}

type EndpointContext interface {
	UUID() string
	endpointContext()
}

type ClientOption interface {
	Option
	clientOption()
}

type clientOption struct {
	Option
}

func (*clientOption) clientOption() {}

func newClientOption(n interface{}, v interface{}) ClientOption {
	return &clientOption{option.New(n, v)}
}

type reconnectIdent struct{}
type dialFuncIdent struct{}
type connIdent struct{}

func AutoReconnect(disable ...bool) ClientOption {
	if len(disable) > 0 {
		return newClientOption(reconnectIdent{}, disable[0])
	}
	return newClientOption(reconnectIdent{}, true)
}

func Dial(dialFunc DialFunc) ClientOption {
	return newClientOption(dialFuncIdent{}, dialFunc)
}

func Conn(conn net.Conn) ClientOption {
	return newClientOption(connIdent{}, conn)
}

type EndpointUUID string

func (e EndpointUUID) endpointContext()   {}
func (e EndpointUUID) UUID() EndpointUUID { return e }
func (e EndpointUUID) String() string     { return string(e) }

type Client interface {
	Host() string

	Connect(opts ...ConnectOption) error
	SetConn(conn net.Conn, opts ...SetConnOption) error
	GetConn() net.Conn
	Close() error
	Closed() bool
	UsedAt() time.Time
	SetUsedAt(tm time.Time)

	NewEndpoint(opts ...Option) (EndpointUUID, Endpoint, error)
	CloseEndpoint(endpoint EndpointContext) error
	Endpoints() []Endpoint

	// Генерациия автоматическая

	// GetClusters(ctx context.Context, endpoint EndpointContext, req *v1.GetClustersRequest) (*v1.GetClustersResponse, error)
	// GetClusterInfo(ctx context.Context, endpoint EndpointContext, req *v1.GetClusterInfoRequest) (*v1.GetClusterInfoResponse, error)
	// RegCluster(ctx context.Context, endpoint EndpointContext, req *v1.RegClusterRequest) (*v1.RegClusterResponse, error)
	// UnregCluster(ctx context.Context, endpoint EndpointContext, req *v1.UnregClusterRequest) (*emptypb.Empty, error)
	// Authenticate(ctx context.Context, endpoint EndpointContext, req *v1.ClusterAuthenticateRequest) (*emptypb.Empty, error)
	// GetManagers(ctx context.Context, endpoint EndpointContext, req *v1.GetClusterManagersRequest) (*v1.GetClusterManagersResponse, error)
	// GetManagerInfo(ctx context.Context, endpoint EndpointContext, req *v1.GetClusterManagerInfoRequest) (*v1.GetClusterManagerInfoResponse, error)
	// GetWorkingProcesses(ctx context.Context, endpoint EndpointContext, req *v1.GetWorkingProcessesRequest) (*v1.GetWorkingProcessesResponse, error)
	// GetWorkingProcessInfo(ctx context.Context, endpoint EndpointContext, req *v1.GetWorkingProcessInfoRequest) (*v1.GetWorkingProcessInfoResponse, error)
	// GetWorkingServers(ctx context.Context, endpoint EndpointContext, req *v1.GetWorkingServersRequest) (*v1.GetWorkingServersResponse, error)
	// GetWorkingServerInfo(ctx context.Context, endpoint EndpointContext, req *v1.GetWorkingServerInfoRequest) (*v1.GetWorkingServerInfoResponse, error)
	// AddWorkingServer(ctx context.Context, endpoint EndpointContext, req *v1.AddWorkingServerRequest) (*v1.AddWorkingServerResponse, error)
	// DeleteWorkingServer(ctx context.Context, endpoint EndpointContext, req *v1.DeleteWorkingServerRequest) (*emptypb.Empty, error)
	//
	negotiate(ctx context.Context, req *protocolv1.NegotiateMessage, opts ...RequestOption) error
	connect(ctx context.Context, req *protocolv1.ConnectMessage, opts ...RequestOption) (*protocolv1.ConnectMessageAck, error)
	disconnect(ctx context.Context, req *protocolv1.DisconnectMessage, opts ...RequestOption) error
	endpointOpen(ctx context.Context, req *protocolv1.EndpointOpen, opts ...RequestOption) (*protocolv1.EndpointOpenAck, error)
	endpointClose(ctx context.Context, req *protocolv1.EndpointClose, opts ...RequestOption) error
	endpointMessage(ctx context.Context, req *protocolv1.EndpointMessage, opts ...RequestOption) (*protocolv1.EndpointMessage, error)

	Request(ctx context.Context, req interface{}, reply interface{}, opts ...RequestOption) error
	EndpointRequest(ctx context.Context, endpoint EndpointContext, req protocolv1.EndpointMessageFormatter, reply protocolv1.EndpointMessageParser, opts ...RequestOption) error
}

func NewClient(addr string, opts ...ClientOption) Client {
	c := newClient(addr, opts...)

	return c
}

func newClient(addr string, opts ...ClientOption) *client {
	c := &client{
		addr:            addr,
		endpoints:       map[EndpointUUID]Endpoint{},
		endpointConfig:  map[EndpointUUID]EndpointConfig{},
		endpointOptions: map[interface{}]EndpointOption{},
		connectOptions:  map[interface{}]ConnectOption{},
		dial:            defaultDial,
	}
	c.applyOptions(opts...)
	return c
}

func defaultDial(addr string) (net.Conn, error) {

	return net.Dial("tcp", addr)

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

type DialFunc func(addr string) (net.Conn, error)

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

	endpoints      map[EndpointUUID]Endpoint
	endpointConfig map[EndpointUUID]EndpointConfig
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

func (c *client) negotiate(ctx context.Context, req *protocolv1.NegotiateMessage, opts ...RequestOption) error {

	if err := c.Request(ctx, req, nil, opts...); err != nil {
		return err
	}

	return nil
}

func (c *client) connect(ctx context.Context, req *protocolv1.ConnectMessage, opts ...RequestOption) (*protocolv1.ConnectMessageAck, error) {

	resp := new(protocolv1.ConnectMessageAck)

	err := c.Request(ctx, req, resp, opts...)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) disconnect(ctx context.Context, req *protocolv1.DisconnectMessage, opts ...RequestOption) error {

	if err := c.Request(ctx, req, nil, opts...); err != nil {
		return err
	}
	return nil
}

func (c *client) endpointOpen(ctx context.Context, req *protocolv1.EndpointOpen, opts ...RequestOption) (*protocolv1.EndpointOpenAck, error) {

	resp := new(protocolv1.EndpointOpenAck)

	err := c.Request(ctx, req, resp, opts...)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) endpointClose(ctx context.Context, req *protocolv1.EndpointClose, opts ...RequestOption) error {

	if err := c.Request(ctx, req, nil, opts...); err != nil {
		return err
	}
	return nil
}

func (c *client) endpointMessage(ctx context.Context, req *protocolv1.EndpointMessage, opts ...RequestOption) (*protocolv1.EndpointMessage, error) {

	resp := new(protocolv1.EndpointMessage)

	err := c.Request(ctx, req, resp, opts...)
	if err != nil {
		return nil, err
	}

	return resp, nil
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

func (c *client) GetConn() net.Conn {
	panic("implement me")
}

func (c *client) Closed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.closed()
}

func (c *client) NewEndpoint(opts ...EndpointOption) (EndpointUUID, Endpoint, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return "", Endpoint{}, nil
}

func (c *client) CloseEndpoint(endpoint EndpointContext) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.endpoints[EndpointUUID(endpoint.UUID())]; ok {
		err := c.endpointClose(context.Background(), &protocolv1.EndpointClose{
			// TODO Где поля????
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *client) Endpoints() []Endpoint {
	c.mu.Lock()
	defer c.mu.Unlock()

	var endpoints []Endpoint
	for _, endpoint := range c.endpoints {
		endpoints = append(endpoints, endpoint)
	}

	return endpoints
}

type Endpoint struct {
	uuid string
	id   string
}

func (c *client) Request(ctx context.Context, req protocolv1.PacketMessageFormatter, reply protocolv1.PacketMessageParser, opts ...RequestOption) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed() {
		return fmt.Errorf("client is closed")
	}

	err := streamRequest(ctx, c.cc, req, reply, opts...)
	if err != nil {
		return err
	}
	return err
}

func (c *client) EndpointRequest(ctx context.Context, endpoint EndpointContext, req protocolv1.EndpointMessageFormatter, reply protocolv1.EndpointMessageParser, opts ...RequestOption) error {

	e := c.endpoints[EndpointUUID(endpoint.UUID())]
	version := e.GetVersion()
	id := e.GetId()
	format := e.GetFormat()

	buf := &bytes.Buffer{}
	if err := req.Formatter(buf, version); err != nil {
		return err
	}
	message := &protocolv1.EndpointMessage{
		Type:       protocolv1.EndpointDataType_ENDPOINT_DATA_TYPE_MESSAGE,
		Format:     format,
		EndpointId: id,
		Data: &protocolv1.EndpointMessage_Message{
			Message: &protocolv1.EndpointDataMessage{
				Bytes: buf.Bytes(),
				Type:  req.GetMessageType(),
			},
		},
	}

	resp, err := c.endpointMessage(ctx, message, opts...)

	if err != nil {
		return err
	}

	if reply != nil {
		// Переписать функцию EndpointMessage.Unpack(version, into)
		switch resp.GetType() {
		case protocolv1.EndpointDataType_ENDPOINT_DATA_TYPE_MESSAGE:
			buf := bytes.NewBuffer(resp.GetMessage().GetBytes())
			if err := reply.Parse(buf, version); err != nil {
				return err
			}
		case protocolv1.EndpointDataType_ENDPOINT_DATA_TYPE_VOID_MESSAGE:

		case protocolv1.EndpointDataType_ENDPOINT_DATA_TYPE_EXCEPTION:
			return resp.GetFailure()
		default:
			return fmt.Errorf("unknown message type <%s>", resp.GetType())
		}
	}
	return nil
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

	ctx := context.Background()
	var err error

	if atomic.CompareAndSwapUint32(&c._connected, 0, 1) {
		err = c.disconnect(ctx, &protocolv1.DisconnectMessage{})
		if err != nil {
			return err
		}
	}

	if c.closed() {
		return nil
	}

	return c.cc.Close()
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

func streamRequest(ctx context.Context, conn net.Conn, req protocolv1.PacketMessageFormatter, reply protocolv1.PacketMessageParser, opts ...RequestOption) error {

	// opts
	// ReadTimeout
	// WriteTimeout
	// Timeout
	options := defaultStreamOptions
	options.applyOptions(opts...)

	cs := newClientStream()
	cs.Options(options)

	if err := cs.SendMsg(ctx, conn, req); err != nil {
		return err
	}

	if reply == nil {
		return nil
	}

	return cs.RecvMsg(ctx, conn, reply)
}

var defaultStreamOptions = StreamOptions{}

type streamPacket struct {
	opts StreamOptions
}

func (s streamPacket) Options(opts StreamOptions) {
	panic("implement me")
}

func newClientStream() StreamPacket {
	return streamPacket{}
}

type StreamOptions struct{}

func (s *StreamOptions) applyOptions(opts ...RequestOption) {

}

type StreamPacket interface {
	Options(opts StreamOptions)
	SendMsg(ctx context.Context, conn net.Conn, m protocolv1.PacketMessageFormatter) error
	RecvMsg(ctx context.Context, conn net.Conn, m protocolv1.PacketMessageParser) error
}

func (s streamPacket) SendMsg(ctx context.Context, conn net.Conn, m protocolv1.PacketMessageFormatter) error {

	if m.GetPacketType() == protocolv1.PacketType_PACKET_TYPE_NEGOTIATE {
		if err := m.Formatter(conn, 0); err != nil {
			return err
		}
		return nil
	}

	packet := getPacket()
	defer putPacket(packet)

	buf := &bytes.Buffer{}

	if err := m.Formatter(buf, 0); err != nil {
		return err
	}

	packet.Type = m.GetPacketType()
	packet.Data = buf.Bytes()
	packet.Size = int32(len(packet.Data))

	_, err := packet.WriteTo(conn)
	if err != nil {
		return err
	}
	return nil

}

func (s streamPacket) RecvMsg(ctx context.Context, conn net.Conn, m protocolv1.PacketMessageParser) error {

	packet := getPacket()
	defer putPacket(packet)

	if err := packet.Parse(conn, 0); err != nil {
		return err
	}
	if err := packet.Unpack(m); err != nil {
		return err
	}

	return nil

}

func getPacket() *protocolv1.Packet {
	return packetPool.Get().(*protocolv1.Packet)
}

func putPacket(packet *protocolv1.Packet) {

	packet.Type = 0
	packet.Size = 0
	packet.Data = []byte{}
	packetPool.Put(packet)
}

var packetPool = sync.Pool{New: func() interface{} {
	return &protocolv1.Packet{}
}}
