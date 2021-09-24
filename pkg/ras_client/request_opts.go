package client

import (
	"context"
	"github.com/lestrrat-go/option"
)

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
func (*requestOption) globalOption()   {}

type endpointIdent struct{}
type getEndpointFromContextIdent struct{}

func newRequestOption(n interface{}, v interface{}) RequestOption {
	return &requestOption{option.New(n, v)}
}

func EndpointUUID(uuid string) RequestOption {
	return newRequestOption(endpointIdent{}, uuid)
}

func GetEndpoint(f func(ctx context.Context) string) RequestOption {
	return newRequestOption(getEndpointFromContextIdent{}, f)
}
