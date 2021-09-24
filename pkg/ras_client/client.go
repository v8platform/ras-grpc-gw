package client

import (
	"context"
	"encoding/binary"
	"github.com/spf13/cast"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type DialFunc func(addr string) (net.Conn, error)

type RequestInfo struct {
	Method     string
	FullMethod string
	Request    interface{}
	Reply      interface{}
	Endpoint   interface{}
}

type InvokeHandler func(ctx context.Context, channel *Channel, opts ...interface{}) error

func GetClusters(ctx context.Context, req *messagesv1.GetClustersRequest, opts ...interface{}) (*messagesv1.GetClustersResponse, error) {
	var cc Client
	var channel *Channel
	endpoint, err := cc.GetEndpoint(ctx)
	if err != nil {
		return nil, err
	}
	return GetClustersHandler(ctx, cc, channel, endpoint, req, opts...)
}

func GetClustersHandler(ctx context.Context, cc Client, channel *Channel, endpoint Endpoint, req *messagesv1.GetClustersRequest, opts ...interface{}) (*messagesv1.GetClustersResponse, error) {

	reply := new(messagesv1.GetClustersResponse)

	handler := func(ctx context.Context, cn *Channel, o ...interface{}) error {
		return endpointRequest(ctx, cn, endpoint, req, reply, o...)
	}

	requestInfo := &RequestInfo{
		Method:     "GetClusters",
		FullMethod: "/ras.api.v1.ClustersService/GetClusters",
		Request:    req,
		Reply:      reply,
		Endpoint:   endpoint,
	}

	err := cc.Invoke(ctx, requestInfo, channel, handler, opts...)
	if err != nil {
		return nil, err
	}

	return reply, nil

}

func endpointRequest(ctx context.Context, channel *Channel, endpoint Endpoint, req *messagesv1.GetClustersRequest, reply *messagesv1.GetClustersResponse, opts ...interface{}) error {

}

type Client interface {
	Host() string

	Close() error
	RemoveChannel(ctx context.Context, cn *Channel)
	Stats() *Stats

	CloseChannel(channel *Channel) error

	PutChannel(ctx context.Context, cn interface{})

	Invoke(ctx context.Context, info *RequestInfo, channel *Channel, handler InvokeHandler, opts ...interface{}) error

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

type EndpointConfig struct {
	Options []RequestOption

	DefaultAgentAuth    Auth
	DefaultClusterAuth  Auth
	DefaultInfobaseAuth Auth

	SaveAuthRequests bool
	Auths            map[string]Auth
}

// Stats contains pool state information and accumulated stats.
type Stats struct {
	Hits     uint32 // number of times free connection was found in the pool
	Misses   uint32 // number of times free connection was NOT found in the pool
	Timeouts uint32 // number of times a wait timeout occurred

	TotalChannels uint32 // number of total connections in the pool
	IdleChannels  uint32 // number of idle connections in the pool
	StaleChannels uint32 // number of stale connections removed from the pool
}

type client struct {
	addr string
	dial DialFunc

	useAutoReconnect bool

	poolSize           int
	poolTimeout        time.Duration
	idleTimeout        time.Duration
	maxChannelAge      time.Duration
	idleCheckFrequency time.Duration
	minIdleChannels    int

	channelsMu      sync.Mutex
	poolSizeLen     int
	idleChannelsLen int
	idleChannels    []*Channel
	channels        []*Channel

	dialErrorsNum   uint32 // atomic
	lastDialErrorMu sync.RWMutex
	lastDialError   error

	queue chan struct{}

	connectOptions  map[interface{}]ConnectOption
	endpointOptions map[interface{}]EndpointOption

	mu      sync.Mutex
	_closed uint32 // atomic

	endpoints map[string]*Endpoint

	endpointConfig map[string]*EndpointConfig

	clientService clientv1.ClientService

	stats Stats

	clientv1.ClustersService
	clientv1.AdminService
	clientv1.AuthService
	clientv1.ConnectionsService
	clientv1.InfobasesService
	clientv1.LocksService
	clientv1.SessionsService
}

func (c *client) AuthenticateCluster(ctx context.Context, req *messagesv1.ClusterAuthenticateRequest, opts ...interface{}) (*emptypb.Empty, error) {

	reply, err := c.AuthService.AuthenticateCluster(ctx, req, opts...)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (c *client) AuthenticateInfobase(ctx context.Context, req *messagesv1.AuthenticateInfobaseRequest, opts ...interface{}) (*emptypb.Empty, error) {

	reply, err := c.AuthService.AuthenticateInfobase(ctx, req, opts...)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (c *client) AuthenticateServer(ctx context.Context, req *messagesv1.ServerAuthenticateRequest, opts ...interface{}) (*emptypb.Empty, error) {

	reply, err := c.AuthService.AuthenticateServer(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (c *client) CloseChannel(channel *Channel) error {

	c.removeChannelWithLock(channel)
	c.freeTurn()
	return c.closeChannel(channel)
}

func (c *client) GetChannel(ctx context.Context, opts ...interface{}) (clientv1.Channel, error) {

	channel := ChannelFromContext(ctx)

	if channel != nil && !channel.Closed() {
		return channel, nil
	}

	var (
		opt         Option
		ok          bool
		initChannel bool
	)

	for _, val := range opts {

		if opt, ok = val.(Option); !ok {
			continue
		}

		switch opt.Ident() {
		case initChannelIdent{}:
			initChannel = opt.Value().(bool)
		}
	}

	if initChannel {
		return c.getChannel(ctx)
	}

	return nil, ErrNoChannel

}

func (c *client) PutChannel(ctx context.Context, cn interface{}) {

	var (
		channel *Channel
		ok      bool
	)

	if channel, ok = cn.(*Channel); !ok {
		return
	}

	if !channel.pooled {
		c.RemoveChannel(ctx, channel)
		return
	}

	c.channelsMu.Lock()
	c.idleChannels = append(c.idleChannels, channel)
	c.idleChannelsLen++
	c.channelsMu.Unlock()
	c.freeTurn()
}

func (c *client) GetEndpoint(ctx context.Context, opts ...interface{}) (clientv1.Channel, clientv1.Endpoint, error) {

	if c.closed() {
		return nil, nil, ErrClosed
	}

	uuid := EndpointFromContext(ctx)
	endpoint, ok := c.endpoints[uuid]

	if ok {

	}

	channel := ChannelFromContext(ctx)

	if channel != nil && !channel.Closed() {
		c.removeChannel(channel)
		// TODO Удалить отключенный канал
		channel = nil
	}

	channel, err := c.getChannel(ctx)
	if err != nil {
		return nil, nil, err
	}

	idxEndpoints := channel.Endpoints()

	if len(idxEndpoints) > 0 {
		if id, ok := idxEndpoints[endpoint.UUID]; ok {
			return channel, newChannelEndpoint(id, endpoint.Ver), err
		}
	}

	id, err := c.initEndpoint(ctx, channel, endpoint)

	if err != nil {
		return nil, nil, err
	}

	return channel, newChannelEndpoint(id, endpoint.Ver), err
}

func (c *client) RemoveChannel(ctx context.Context, cn *Channel) {
	c.removeChannelWithLock(cn)
	c.freeTurn()
	_ = c.closeChannel(cn)
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

	c.queue = make(chan struct{}, c.poolSize)
	c.idleChannels = make([]*Channel, c.poolSize)
	c.channels = make([]*Channel, c.poolSize)

	c.channelsMu.Lock()
	c.checkMinIdleChannels()
	c.channelsMu.Unlock()

	if c.idleTimeout > 0 && c.idleCheckFrequency > 0 {
		go c.reaper(c.idleCheckFrequency)
	}

	return nil
}

func (c *client) CloseAllChannels() error {
	return nil
}

func (c *client) Stats() *Stats {
	return &Stats{
		Hits:     atomic.LoadUint32(&c.stats.Hits),
		Misses:   atomic.LoadUint32(&c.stats.Misses),
		Timeouts: atomic.LoadUint32(&c.stats.Timeouts),

		TotalChannels: uint32(c.channelsLen()),
		IdleChannels:  uint32(c.idleChannelsLen),
		StaleChannels: atomic.LoadUint32(&c.stats.StaleChannels),
	}
}

func (c *client) Close() error {
	if !atomic.CompareAndSwapUint32(&c._closed, 0, 1) {
		return ErrClosed
	}

	var firstErr error
	c.channelsMu.Lock()
	for _, ch := range c.idleChannels {
		if err := ch.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	c.idleChannels = nil
	c.poolSizeLen = 0
	c.idleChannelsLen = 0
	c.channelsMu.Unlock()

	return firstErr
}

func (c *client) getChannel(ctx context.Context) (*Channel, error) {

	if c.closed() {
		return nil, ErrClosed
	}

	err := c.waitTurn(ctx)
	if err != nil {
		return nil, err
	}

	for {

		c.channelsMu.Lock()
		channel := c.getIdleChannel()
		c.channelsMu.Unlock()

		if channel == nil {
			break
		}

		if c.isStaleChannel(channel) {
			_ = c.closeChannel(channel)
			continue
		}

		atomic.AddUint32(&c.stats.Hits, 1)

		return channel, nil
	}

	atomic.AddUint32(&c.stats.Misses, 1)

	channel, err := c.newChannel(true)

	if err != nil {
		c.freeTurn()
		return nil, err
	}

	err = c.initChannel(ctx, channel)

	if err != nil {
		return nil, err
	}

	return channel, err

}

func (c *client) initEndpoint(ctx context.Context, channel *Channel, endpoint *Endpoint) (int32, error) {

	err := c.initChannel(ctx, channel)

	if err != nil {
		return 0, err
	}

	reply, err := clientv1.EndpointOpenHandler(ctx, channel, &protocolv1.EndpointOpen{
		Version: endpoint.version,
		Service: protocolv1.ServiceName,
	})

	if err != nil {
		version := clientv1.DetectSupportedVersion(err)
		if len(version) == 0 {
			return 0, err
		}

		endpoint.version = version
		endpoint.Ver = cast.ToInt32(version)
	}

	channel.SetEndpoint(endpoint.UUID, reply.EndpointId)

	// TODO init auth

	return reply.EndpointId, nil
}

func (c *client) removeChannelWithLock(cn *Channel) {
	c.channelsMu.Lock()
	c.removeChannel(cn)
	c.channelsMu.Unlock()
}

// Len returns total number of connections.
func (c *client) channelsLen() int {
	c.channelsMu.Lock()
	n := len(c.channels)
	c.channelsMu.Unlock()
	return n
}

func (c *client) getIdleChannel() *Channel {

	if len(c.idleChannels) == 0 {
		return nil
	}

	idx := len(c.idleChannels) - 1
	cn := c.idleChannels[idx]
	c.idleChannels = c.idleChannels[:idx]

	c.idleChannelsLen--
	c.checkMinIdleChannels()

	return cn
}

func (c *client) checkMinIdleChannels() {
	if c.minIdleChannels == 0 {
		return
	}
	for c.poolSizeLen < c.poolSize && c.idleChannelsLen < c.minIdleChannels {
		c.poolSizeLen++
		c.idleChannelsLen++
		go func() {
			err := c.addIdleChannel()
			if err != nil {
				c.channelsMu.Lock()
				c.poolSizeLen--
				c.idleChannelsLen--
				c.channelsMu.Unlock()
			}
		}()
	}
}

func (c *client) applyOptions(opts ...ClientOption) {

	for _, opt := range opts {
		switch opt.Ident() {
		case dialFuncIdent{}:
			c.dial = opt.Value().(DialFunc)
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

func (c *client) restoreEndpoints(opts ...Option) error {

	return nil
}

func (c *client) reaper(frequency time.Duration) {

	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for range ticker.C {
		if c.closed() {
			break
		}
		n, err := c.reapStaleChannels()
		if err != nil {
			continue
		}
		atomic.AddUint32(&c.stats.StaleChannels, uint32(n))
	}
}

func (c *client) freeTurn() {
	<-c.queue
}

func (c *client) getTurn() {
	c.queue <- struct{}{}
}

func (c *client) waitTurn(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	select {
	case c.queue <- struct{}{}:
		return nil
	default:
	}

	timer := timers.Get().(*time.Timer)
	timer.Reset(c.poolTimeout)

	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		timers.Put(timer)
		return ctx.Err()
	case c.queue <- struct{}{}:
		if !timer.Stop() {
			<-timer.C
		}
		timers.Put(timer)
		return nil
	case <-timer.C:
		timers.Put(timer)
		atomic.AddUint32(&c.stats.Timeouts, 1)
		return ErrTimeout
	}
}

func (c *client) reapStaleChannels() (int, error) {
	var n int
	for {
		c.getTurn()

		c.channelsMu.Lock()
		cn := c.reapStaleChannel()
		c.channelsMu.Unlock()

		c.freeTurn()

		if cn != nil {
			_ = c.closeChannel(cn)
			n++
		} else {
			break
		}
	}
	return n, nil
}

func (c *client) closeChannel(cn *Channel) error {
	return cn.Close()
}

func (c *client) initChannel(ctx context.Context, cn *Channel) error {

	if cn.inited {
		return nil
	}

	_, err := clientv1.NegotiateHandler(ctx, cn, protocolv1.NewNegotiateMessage())
	if err != nil {
		return err
	}

	var connectParam = map[string]*protocolv1.Param{}

	for _, option := range c.connectOptions {
		switch option.Ident() {
		case timeoutIdent{}:
			var b []byte
			binary.BigEndian.PutUint32(b, option.Value().(uint32))
			connectParam["connect.timeout"] = &protocolv1.Param{
				Type:  protocolv1.ParamType_PARAM_INT,
				Value: b,
			}
		}
	}

	_, err = clientv1.ConnectHandler(ctx, cn, &protocolv1.ConnectMessage{
		Params: connectParam,
	})

	if err != nil {
		return err
	}

	return nil

}

func (c *client) newChannel(pooled bool) (*Channel, error) {
	cn, err := c.dialChannel(pooled)
	if err != nil {
		return nil, err
	}

	c.channelsMu.Lock()

	c.channels = append(c.channels, cn)
	if pooled {
		// If pool is full remove the cn on next Put.
		if c.poolSize >= c.poolSize {
			cn.pooled = false
		} else {
			c.poolSize++
		}
	}

	c.channelsMu.Unlock()
	return cn, nil
}

func (c *client) addIdleChannel() error {

	cn, err := c.dialChannel(true)
	if err != nil {
		return err
	}

	c.channelsMu.Lock()
	c.channels = append(c.channels, cn)
	c.idleChannels = append(c.idleChannels, cn)
	c.channelsMu.Unlock()
	return nil
}

func (c *client) reapStaleChannel() *Channel {
	if len(c.idleChannels) == 0 {
		return nil
	}

	cn := c.idleChannels[0]
	if !c.isStaleChannel(cn) {
		return nil
	}

	c.idleChannels = append(c.idleChannels[:0], c.idleChannels[1:]...)
	c.idleChannelsLen--
	c.removeChannel(cn)

	return cn
}

func (c *client) removeChannel(cn *Channel) {
	for i, channel := range c.channels {
		if channel == cn {
			c.channels = append(c.channels[:i], c.channels[i+1:]...)
			if cn.pooled {
				c.poolSize--
				c.checkMinIdleChannels()
			}

			return
		}
	}
}

func (c *client) dialChannel(pooled bool) (*Channel, error) {
	if c.closed() {
		return nil, ErrClosed
	}

	if atomic.LoadUint32(&c.dialErrorsNum) >= uint32(c.poolSize) {
		return nil, c.getLastDialError()
	}

	netConn, err := c.dial(c.addr)
	if err != nil {
		c.setLastDialError(err)
		if atomic.AddUint32(&c.dialErrorsNum, 1) == uint32(c.poolSize) {
			go c.tryDial()
		}
		return nil, err
	}

	cn := newChannel(netConn)
	cn.pooled = pooled
	return cn, nil
}

func (c *client) tryDial() {
	for {
		if c.closed() {
			return
		}

		conn, err := c.dial(c.addr)
		if err != nil {
			c.setLastDialError(err)
			time.Sleep(time.Second)
			continue
		}

		atomic.StoreUint32(&c.dialErrorsNum, 0)
		_ = conn.Close()
		return
	}
}

func (c *client) closed() bool {
	return atomic.LoadUint32(&c._closed) == 1
}

func (c *client) setLastDialError(err error) {
	c.lastDialErrorMu.Lock()
	c.lastDialError = err
	c.lastDialErrorMu.Unlock()
}

func (c *client) getLastDialError() error {
	c.lastDialErrorMu.RLock()
	err := c.lastDialError
	c.lastDialErrorMu.RUnlock()
	return err
}

func (c *client) isStaleChannel(cn *Channel) bool {

	if cn.Closed() {
		return true
	}

	if c.idleTimeout == 0 && c.maxChannelAge == 0 {
		return false
	}

	now := time.Now()
	if c.idleTimeout > 0 && now.Sub(cn.UsedAt()) >= c.idleTimeout {
		return true
	}

	if c.maxChannelAge > 0 && now.Sub(cn.CreatedAt()) >= c.maxChannelAge {
		return true
	}

	return false
}

func newClient(addr string, opts ...ClientOption) *client {
	c := &client{
		addr:               addr,
		endpoints:          map[string]*Endpoint{},
		endpointConfig:     map[string]*EndpointConfig{},
		endpointOptions:    map[interface{}]EndpointOption{},
		connectOptions:     map[interface{}]ConnectOption{},
		dial:               defaultDial,
		poolSize:           10,
		poolTimeout:        30 * time.Second,
		idleTimeout:        time.Hour,
		idleCheckFrequency: 5 * time.Minute,
		maxChannelAge:      time.Hour,
		minIdleChannels:    0,
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

	c.queue = make(chan struct{}, c.poolSize)

	c.idleChannels = make([]*Channel, c.poolSize)
	c.channels = make([]*Channel, c.poolSize)

	c.channelsMu.Lock()
	c.checkMinIdleChannels()
	c.channelsMu.Unlock()

	if c.idleTimeout > 0 && c.idleCheckFrequency > 0 {
		go c.reaper(c.idleCheckFrequency)
	}

	return c
}

func defaultDial(addr string) (net.Conn, error) {

	return net.Dial("tcp", addr)

}
