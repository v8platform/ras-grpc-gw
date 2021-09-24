package client

import (
	"github.com/lestrrat-go/option"
	protocolv1 "github.com/v8platform/protos/gen/ras/protocol/v1"
	"net"
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

func EndpointData(id, version, format int32) RequestOption {
	return newRequestOption(endpointIdent{}, endpointData{id, version, format})
}

func EndpointUUID(uuid string) RequestOption {
	return newRequestOption(endpointIdent{}, uuid)
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
type initChannelIdent struct{}

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
func NewConn(init bool) EndpointOption {
	return newEndpointOption(initChannelIdent{}, init)
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
