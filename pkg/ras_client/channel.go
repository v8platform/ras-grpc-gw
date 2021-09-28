package client

import (
	"context"
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
		id:          0, // TODO Сделать увеличение номеров
		conn:        conn,
		_closed:     0,
		_usedAt:     0,
		createdAt:   time.Now(),
		idxEndpoint: map[uuid.UUID]int32{},
	}
}

type Channel struct {
	id   int32
	conn net.Conn
	sync.Mutex

	_closed   uint32
	_usedAt   uint32 // atomic
	createdAt time.Time
	pooled    bool
	inited    bool

	idxEndpoint map[uuid.UUID]int32

	recvLen int
	recvMsg []interface{}

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

func (c *Channel) SetEndpoint(endpoint uuid.UUID, id int32) {

	c.Lock()
	defer c.Unlock()

	c.idxEndpoint[endpoint] = id

}

func (c *Channel) Endpoints() map[uuid.UUID]int32 {

	c.Lock()
	defer c.Unlock()

	endpoints := make(map[uuid.UUID]int32, len(c.idxEndpoint))
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

func (c *Channel) ID() int32 {
	return c.id
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

func (c *Channel) RecvPacket(ctx context.Context) (msg *protocolv1.Packet, err error) {

	c.Lock()
	defer func() { c.SetUsedAt(time.Now()) }()
	defer c.Unlock()

	msg = new(protocolv1.Packet)
	done := c.recv(msg, err)

	select {
	case <-done:
		if err != nil {
			return nil, err
		}
		return msg, nil
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

func (c *Channel) CreatedAt() time.Time {
	return c.createdAt
}
