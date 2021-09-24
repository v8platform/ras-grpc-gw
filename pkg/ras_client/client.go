package client

import (
	"context"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type DialFunc func(addr string) (net.Conn, error)

type Client interface {
	Host() string

	Connect(opts ...ConnectOption) error
	Close() error
	RemoveChannel(ctx context.Context, cn *Channel)
	Stats() *Stats

	NewEndpoint(opts ...EndpointOption) (*Channel, *Endpoint, error)
	CloseChannel(channel *Channel) error

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
	maxConnAge         time.Duration
	idleCheckFrequency time.Duration
	minIdleChannels    int

	channelsMu      sync.Mutex
	poolSizeLen     int
	idleChannelsLen int
	idleChannels    []*Channel

	dialErrorsNum   uint32 // atomic
	lastDialErrorMu sync.RWMutex
	lastDialError   error

	_closed uint32 // atomic
	queue   chan struct{}

	connectOptions  map[interface{}]ConnectOption
	endpointOptions map[interface{}]EndpointOption
	channels        []*Channel

	mu sync.Mutex

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

func (c *client) NewEndpoint(opts ...EndpointOption) (*Channel, *Endpoint, error) {
	panic("implement me")
}

func (c *client) CloseChannel(channel *Channel) error {

	return nil
}

func (c *client) GetChannel(ctx context.Context, opts ...interface{}) (clientv1.Channel, error) {

	channel := ChannelFromContext(ctx)

	if channel != nil && !channel.Closed() {
		return channel, nil
	}

	return c.initNewChannel(opts...)
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

	channel := ChannelFromContext(ctx)

	if channel != nil && !channel.Closed() {
		c.removeChannel(channel)
		// TODO Удалить отключенный канал
		channel = nil
	}

	err := c.waitTurn(ctx)
	if err != nil {
		return nil, nil, err
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

		// TODO открытие точки в данном канале
		// endpoint, err := p.openEndpoint(ctx, newConn)

		if !endpoint.Inited {
			endpoint, err = p.openEndpoint(ctx, endpoint.conn)

			if err != nil {
				return nil, nil, err
			}

		}

		atomic.AddUint32(&c.stats.Hits, 1)

		return channel, endpoint, nil
	}

	atomic.AddUint32(&c.stats.Misses, 1)

	channel, err = c.newChannel(true)

	if err != nil {
		c.freeTurn()
		return nil, nil, err
	}

	// TODO открытие точки в данном канале
	// endpoint, err := p.openEndpoint(ctx, newConn)

	return channel, endpoint, err

	// return c.NewEndpoint()

}

func (p *client) RemoveChannel(ctx context.Context, cn *Channel) {
	p.removeChannelWithLock(cn)
	p.freeTurn()
	_ = p.closeChannel(cn)
}

func (p *client) removeChannelWithLock(cn *Channel) {
	p.channelsMu.Lock()
	p.removeChannel(cn)
	p.channelsMu.Unlock()
}

// Len returns total number of connections.
func (p *client) channelsLen() int {
	p.channelsMu.Lock()
	n := len(p.channels)
	p.channelsMu.Unlock()
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

func (c *client) Channels() []*Channel {
	return c.channels
}

func (c *client) CloseAllChannels() error {
	return nil
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

	// c.cc = conn

	if restoreConnect || restoreEndpoints {
		// if err := c.initConnect(); err != nil {
		// 	return err
		// }
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

func (c *client) initNewChannel(opts ...interface{}) (*Channel, error) {
	return nil, nil
}

func (c *client) getFreeChannel(endpoint *Endpoint) (*Channel, error) {
	return c.initNewChannel()
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

	cn := NewChannel(netConn)
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

	if c.idleTimeout == 0 && c.maxConnAge == 0 {
		return false
	}

	now := time.Now()
	if c.idleTimeout > 0 && now.Sub(cn.UsedAt()) >= c.idleTimeout {
		return true
	}

	if c.maxConnAge > 0 && now.Sub(cn.CreatedAt()) >= c.maxConnAge {
		return true
	}

	return false
}

func newClient(addr string, opts ...ClientOption) *client {
	c := &client{
		addr:             addr,
		endpoints:        map[string]*Endpoint{},
		endpointConfig:   map[string]*EndpointConfig{},
		endpointOptions:  map[interface{}]EndpointOption{},
		connectOptions:   map[interface{}]ConnectOption{},
		channels:         []*Channel{},
		endpointChannels: map[string][]int32{},
		dial:             defaultDial,
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
