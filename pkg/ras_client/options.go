package client

import (
	"fmt"
	"github.com/elastic/go-ucfg"
	"github.com/lestrrat-go/option"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client/md"
	"reflect"
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
type contextAnnotatorIdent struct{}

func Dial(dialFunc DialFunc) GlobalOption {
	return newGlobalOption(dialFuncIdent{}, dialFunc)
}

func WithConfig(config Config) GlobalOption {
	return newGlobalOption(configIdent{}, config)
}
func ConfigFrom(cfg *ucfg.Config) GlobalOption {
	return newGlobalOption(configFromIdent{}, cfg)
}

func AnnotateContext(handler md.AnnotationHandler) GlobalOption {
	return newGlobalOption(contextAnnotatorIdent{}, handler)
}

func combine(opts1 interface{}, o2 []interface{}) []Option {
	// we don't use append because o1 could have extra capacity whose
	// elements would be overwritten, which could cause inadvertent
	// sharing (and race conditions) between concurrent calls
	var o1 []Option

	switch typed_o1 := opts1.(type) {
	case []interface{}:
		o1 = toOptions(typed_o1)
	case []Option:
		o1 = typed_o1
	case []GlobalOption:
		for _, i := range typed_o1 {
			if io, ok := i.(Option); ok {
				o1 = append(o1, io)
			}
		}
	case map[interface{}]GlobalOption:
		for _, i := range typed_o1 {
			if io, ok := i.(Option); ok {
				o1 = append(o1, io)
			}
		}
	case []EndpointOption:
		for _, i := range typed_o1 {
			if io, ok := i.(Option); ok {
				o1 = append(o1, io)
			}
		}
	case map[interface{}]EndpointOption:
		for _, i := range typed_o1 {
			if io, ok := i.(Option); ok {
				o1 = append(o1, io)
			}
		}
	case []RequestOption:
		for _, i := range typed_o1 {
			if io, ok := i.(Option); ok {
				o1 = append(o1, io)
			}
		}
	case nil:
		// nothing
	default:
		panic(fmt.Errorf("combine: unknown slice type: %s", reflect.TypeOf(typed_o1)))
	}

	if len(o1) == 0 {
		return toOptions(o2)
	} else if len(o2) == 0 {
		return o1
	}
	ret := make([]Option, len(o1)+len(o2))
	copy(ret, o1)
	copy(ret[len(o1):], toOptions(o2))
	return ret
}

func toOptions(o2 []interface{}) []Option {
	if len(o2) == 0 {
		return []Option{}
	}
	var opts []Option
	for _, i := range o2 {
		if io, ok := i.(Option); ok {
			opts = append(opts, io)
		}
	}

	return opts
}
