package client

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
	"io"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

var _ clientv1.Channel = (*Channel)(nil)

func newChannel(conn net.Conn) *Channel {
	return &Channel{
		conn:        conn,
		_closed:     0,
		_usedAt:     0,
		createdAt:   time.Now(),
		idxEndpoint: map[uuid.UUID]*ChannelEndpoint{},
	}
}

type Channel struct {
	conn net.Conn
	sync.Mutex

	_closed   uint32
	_usedAt   uint32 // atomic
	createdAt time.Time
	pooled    bool
	inited    bool

	idxEndpoint map[uuid.UUID]*ChannelEndpoint

	recvWg sync.WaitGroup
}

func (c *Channel) Close() error {

	if c.Closed() {
		return nil
	}

	return c.conn.Close()
}

func (c *Channel) Closed() bool {

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

func (c *Channel) SetEndpoint(id uuid.UUID, endpoint *ChannelEndpoint) {

	c.Lock()
	defer c.Unlock()

	c.idxEndpoint[id] = endpoint

}

func (c *Channel) Endpoints() map[uuid.UUID]*ChannelEndpoint {

	c.Lock()
	defer c.Unlock()

	endpoints := make(map[uuid.UUID]*ChannelEndpoint, len(c.idxEndpoint))
	for key, value := range c.idxEndpoint {
		endpoints[key] = value
	}
	return nil
}

func (c *Channel) UsedAt() time.Time {
	unix := atomic.LoadUint32(&c._usedAt)
	return time.Unix(int64(unix), 0)
}

func (c *Channel) SetUsedAt(tm time.Time) {
	atomic.StoreUint32(&c._usedAt, uint32(tm.Unix()))
}

func (c *Channel) SendMsg(ctx context.Context, msg interface{}, opts ...interface{}) error {
	// TODO Прочитать опции канала?

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.Lock()
	defer func() { c.SetUsedAt(time.Now()) }()
	defer c.Unlock()

	writeTimeout := 3 * time.Second

	switch m := msg.(type) {
	case io.WriterTo:

		if err := c.conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
			return err
		}
		_, err := m.WriteTo(c.conn)
		if err != nil {
			return err
		}

	case []byte:
		if err := c.conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
			return err
		}
		_, err := c.conn.Write(m)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown msg type %s", reflect.TypeOf(msg))
	}

	return nil
}

func (c *Channel) SendPacket(ctx context.Context, msg *protocolv1.Packet) error {
	// TODO Прочитать опции канала?

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.Lock()
	defer func() { c.SetUsedAt(time.Now()) }()
	defer c.Unlock()

	writeTimeout := 3 * time.Second

	if err := c.conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		return err
	}
	_, err := msg.WriteTo(c.conn)
	if err != nil {
		return err
	}

	return nil
}

func (c *Channel) recv(msg *protocolv1.Packet, err error) chan struct{} {

	done := make(chan struct{})
	c.recvWg.Add(1)

	go func() {

		defer close(done)
		defer c.recvWg.Done()

		_, err = msg.ReadFrom(c.conn)

	}()

	return done
}

type frame struct {
	err  error
	done chan struct{}
}

func (f frame) Done() chan struct{} {
	return f.done
}

func (f frame) Err() error {
	return f.err
}

func (c *Channel) recvFrame(msg *protocolv1.Packet) frame {

	f := frame{done: make(chan struct{})}
	c.recvWg.Add(1)

	go func() {

		defer close(f.done)
		defer c.recvWg.Done()

		_, err := msg.ReadFrom(c.conn)
		f.err = err

	}()

	return f
}

func (c *Channel) RecvPacket(ctx context.Context) (*protocolv1.Packet, error) {

	c.Lock()
	defer func() { c.SetUsedAt(time.Now()) }()
	defer c.Unlock()

	packet := new(protocolv1.Packet)

	frame := c.recvFrame(packet)

	select {
	case <-frame.Done():
		if frame.Err() != nil {
			return nil, frame.err
		}
		return packet, nil
	case <-ctx.Done():
		// TODO
		return nil, ctx.Err()
	}
}

func (c *Channel) RecvMsg(ctx context.Context, msg interface{}, opts ...interface{}) (err error) {

	c.Lock()
	defer func() { c.SetUsedAt(time.Now()) }()
	defer c.Unlock()

	packet := msg.(*protocolv1.Packet)

	done := c.recv(packet, err)

	select {
	case <-done:
		if err != nil {
			return err
		}
		return nil
	case <-ctx.Done():
		// TODO
		return ctx.Err()
	}
}

func (c *Channel) IsChannelEndpoint(endpoint *ChannelEndpoint) bool {

	if len(c.idxEndpoint) == 0 {
		return false
	}

	cEndpoint, ok := c.idxEndpoint[endpoint.UUID]

	if ok && cEndpoint == endpoint {
		return true
	}

	return false
}

func (c *Channel) CreatedAt() time.Time {
	return c.createdAt
}

var _ protocolv1.Endpoint = (*ChannelEndpoint)(nil)

type AuthType int

const (
	ClusterAuth AuthType = iota
	InfobaseAuth
)

type ChannelEndpoint struct {
	UUID    uuid.UUID
	ID      int32
	Version int32
	UsedAt  time.Time
	hash    [2]hashMap // Хеши авторизации
}
type hashMap map[string][32]byte

func (m hashMap) CompareAndSwap(key string, hash [32]byte) (old [32]byte, swap bool) {

	old, swap = m[key]
	if swap && compareHashes(old, hash) {
		return old, false
	}

	m[key] = hash
	return old, true

}
func (m hashMap) Compare(key string, hash [32]byte) (equal bool) {

	value, ok := m[key]
	if ok && compareHashes(value, hash) {
		return true
	}
	return false
}

func (m hashMap) Swap(key string, hash [32]byte) (old [32]byte, swap bool) {

	old, swap = m[key]
	m[key] = hash

	return
}

func hash(value string) [32]byte {
	return sha256.Sum256([]byte(value))
}

func compareHashes(hash1, hash2 [32]byte) bool {

	for i := 0; i < 32; i++ {
		if hash1[i] != hash2[i] {
			return false
		}
	}
	return true
}

func (c *ChannelEndpoint) GetUsedAt() time.Time {
	return c.UsedAt
}

func (c *ChannelEndpoint) CompareHash(auth AuthType, key, user, pwd string) bool {

	if c.hash[auth] == nil {
		return false
	}

	return c.hash[auth].Compare(key, hash(user+pwd))

}

func (c *ChannelEndpoint) ChangeAuth(req interface{}) {

}

func (c *ChannelEndpoint) SwapHash(auth AuthType, key, user, pwd string) {

	if c.hash[auth] == nil {
		c.hash[auth] = make(hashMap)
		return
	}
	c.hash[auth].Swap(key, hash(user+pwd))
}

func (c *ChannelEndpoint) SetUsedAt(tm time.Time) {
	c.UsedAt = tm
}

func (e *ChannelEndpoint) GetVersion() int32 {
	return e.Version
}

func (e *ChannelEndpoint) GetId() int32 {
	return e.ID
}
