package client

import (
	"context"
	"fmt"
	uuid2 "github.com/google/uuid"
	"github.com/lestrrat-go/option"
	"time"
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
type interceptorsIdent struct{}
type requestTimeoutIdent struct{}

func newRequestOption(n interface{}, v interface{}) RequestOption {
	return &requestOption{option.New(n, v)}
}

func EndpointUUID(uuid func(ctx context.Context) (uuid2.UUID, bool)) RequestOption {
	return newRequestOption(endpointIdent{}, uuid)
}

func RequestTimeout(d time.Duration) RequestOption {

	return newRequestOption(requestTimeoutIdent{}, d)

}

func RequestInterceptor(interceptor ...Interceptor) RequestOption {

	if len(interceptor) == 0 {
		panic(fmt.Errorf("need 1 or more intercentors"))
	}
	if len(interceptor) > 1 {
		return newRequestOption(interceptorsIdent{}, ChainInterceptor(interceptor...))
	}

	return newRequestOption(interceptorsIdent{}, interceptor[0])

}
