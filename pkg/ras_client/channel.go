package client

import (
	"context"
	"fmt"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	"io"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

var _ clientv1.Channel = (*Channel)(nil)

func NewChannel(conn net.Conn) *Channel {
	return &Channel{
		id:          0, // TODO Сделать увеличение номеров
		conn:        conn,
		_closed:     0,
		_usedAt:     0,
		createdAt:   time.Now(),
		idxEndpoint: map[string]int32{},
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

	idxEndpoint map[string]int32
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

func (c *Channel) SetEndpoint(endpoint string, id int32) {

	c.Lock()
	defer c.Unlock()

	c.idxEndpoint[endpoint] = id

}

func (c *Channel) Endpoints() map[string]int32 {

	c.Lock()
	defer c.Unlock()

	endpoints := make(map[string]int32, len(c.idxEndpoint))
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

func (c *Channel) RecvMsg(ctx context.Context, msg interface{}, opts ...interface{}) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.Lock()
	defer func() { c.SetUsedAt(time.Now()) }()
	defer c.Unlock()

	timeout := 3 * time.Second

	switch m := msg.(type) {
	case io.ReaderFrom:

		if err := c.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
			return err
		}
		_, err := m.ReadFrom(c.conn)
		if err != nil {
			return err
		}

	case []byte:
		if err := c.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
			return err
		}
		_, err := c.conn.Read(m)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown msg type %s", reflect.TypeOf(msg))
	}

	return nil
}

func (c *Channel) CreatedAt() time.Time {
	return c.createdAt
}
