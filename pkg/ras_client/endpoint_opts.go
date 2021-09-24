package client

import "github.com/lestrrat-go/option"

type EndpointOption interface {
	Option
	GlobalOption
	endpointOption()
}

type endpointOption struct {
	Option
}

func (*endpointOption) endpointOption() {}
func (*endpointOption) globalOption()   {}

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
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}

func Version(version int32) EndpointOption {
	return newEndpointOption(versionIdent{}, version)
}
func AutosaveAuth(save bool) EndpointOption {
	return newEndpointOption(saveAuthIdent{}, save)
}
func InitChannel(init bool) EndpointOption {
	return newEndpointOption(initChannelIdent{}, init)
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
