package client

import (
	"github.com/elastic/go-ucfg"
	"github.com/lestrrat-go/option"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
	"time"
)

type Option interface {
	option.Interface
}

type ConnectOption interface {
	Option
	GlobalOption
	connectOption()
}

type connectOption struct {
	Option
}

func (*connectOption) connectOption() {}
func (*connectOption) globalOption()  {}

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

type restoreConnectIdent struct{}
type restoreEndpointsIdent struct{}

type GlobalOption interface {
	Option
	globalOption()
}

type globalOption struct {
	Option
}

func (*globalOption) globalOption() {}

func newGlobalOption(n interface{}, v interface{}) GlobalOption {
	return &globalOption{option.New(n, v)}
}

type dialFuncIdent struct{}
type configIdent struct{}
type configFromIdent struct{}

func Dial(dialFunc DialFunc) GlobalOption {
	return newGlobalOption(dialFuncIdent{}, dialFunc)
}

func WithConfig(config Config) GlobalOption {
	return newGlobalOption(configIdent{}, config)
}
func ConfigFrom(cfg *ucfg.Config) GlobalOption {
	return newGlobalOption(configFromIdent{}, cfg)
}
