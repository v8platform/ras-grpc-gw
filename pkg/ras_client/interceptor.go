package client

import (
	"context"
	"fmt"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	messagesv1 "github.com/v8platform/protos/gen/ras/messages/v1"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client/cond"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client/md"
	"log"
	"reflect"
)

type Interceptor clientv1.Interceptor
type InterceptorHandler clientv1.InterceptorHandler

const (
	InfobaseUserKeys = "infobase-user"
	InfobasePwdKeys  = "infobase-password infobase-pwd"
)

const (
	OverwriteClusterIdKey = "overwrite-cluster-id"
	ClusterIdKeys         = OverwriteClusterIdKey + " cluster-id"
	ClusterUserKeys       = "cluster-user"
	ClusterPwdKeys        = "cluster-password cluster-pwd"
)

func ChainInterceptor(interceptors ...Interceptor) Interceptor {
	n := len(interceptors)

	return func(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {
		chainer := func(currentInter clientv1.Interceptor, currentHandler clientv1.InterceptorHandler) clientv1.InterceptorHandler {
			return func(currentCtx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, currentReq interface{}) (interface{}, error) {
				return currentInter(currentCtx, channel, endpoint, info, currentReq, currentHandler)
			}
		}

		chainedHandler := handler
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(clientv1.Interceptor(interceptors[i]), chainedHandler)
		}

		return chainedHandler(ctx, channel, endpoint, req)
	}
}
func OverwriteClusterIdInterceptor() Interceptor {
	return overwriteClusterIdInterceptor
}
func AddInfobaseAuthInterceptor() Interceptor {
	return addInfobaseAuthInterceptor
}

func AddClusterAuthInterceptor() Interceptor {
	return addClusterAuthInterceptor
}

func overwriteClusterIdInterceptor(ctx context.Context, channel clientv1.Channel,
	endpoint clientv1.Endpoint, info *clientv1.RequestInfo,
	req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {

	condition := cond.And{
		IsEndpointRequest,
		HasMdValues(OverwriteClusterIdKey),
		HasClusterId,
	}

	next := func() (interface{}, error) {
		return handler(ctx, channel, endpoint, req)
	}

	if condition.IsFalse(ctx, endpoint, info, req) {
		return next()
	}

	reqMd := md.ExtractMetadata(ctx)
	clusterId := reqMd.Get(OverwriteClusterIdKey)

	rValue := reflect.ValueOf(req)

	rValue.FieldByName("ClusterId").SetString(clusterId)

	return handler(ctx, channel, endpoint, req)
}

func addClusterAuthInterceptor(ctx context.Context, channel clientv1.Channel, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {

	condition := cond.And{
		IsEndpointRequest,
		HasMdValues(ClusterUserKeys, ClusterPwdKeys),
		NeedClusterAuth,
	}

	next := func() (interface{}, error) {
		return handler(ctx, channel, endpoint, req)
	}

	if condition.IsFalse(ctx, endpoint, info, req) {
		return next()
	}

	type getClusterId interface {
		GetClusterID() string
	}

	reqMd := md.ExtractMetadata(ctx)
	clusterId := reqMd.Get(ClusterIdKeys)

	if len(clusterId) == 0 {
		tReq := req.(getClusterId)
		clusterId = tReq.GetClusterID()
	}

	if len(clusterId) == 0 {
		log.Printf("can`t add cluster auth before request <%s> cluster is not set\n", reflect.TypeOf(req))
		return next()
	}

	_, err := clientv1.AuthenticateClusterHandler(ctx, channel, endpoint, &messagesv1.ClusterAuthenticateRequest{
		ClusterId: clusterId,
		User:      reqMd.Get(ClusterUserKeys),
		Password:  reqMd.Get(ClusterPwdKeys),
	}, nil)

	if err != nil {
		log.Println(err)
	}

	return handler(ctx, channel, endpoint, req)

}

func addInfobaseAuthInterceptor(ctx context.Context, channel clientv1.Channel,
	endpoint clientv1.Endpoint, info *clientv1.RequestInfo,
	req interface{}, handler clientv1.InterceptorHandler) (interface{}, error) {

	condition := cond.And{
		IsEndpointRequest,
		HasMdValues(InfobaseUserKeys, InfobasePwdKeys),
		NeedInfobaseAuth,
	}

	next := func() (interface{}, error) {
		return handler(ctx, channel, endpoint, req)
	}

	if condition.IsFalse(ctx, endpoint, info, req) {
		return next()
	}

	type getClusterId interface {
		GetClusterID() string
	}

	reqMd := md.ExtractMetadata(ctx)

	clusterId := reqMd.Get(ClusterIdKeys)
	if len(clusterId) == 0 {
		tReq := req.(getClusterId)
		clusterId = tReq.GetClusterID()
	}
	if len(clusterId) == 0 {
		log.Printf("can`t add infobase auth before request <%s> cluster is not set\n", reflect.TypeOf(req))
		return next()
	}

	user := reqMd.Get(InfobaseUserKeys)
	password := reqMd.Get(InfobasePwdKeys)

	// TODO Для тестов
	//user := os.Getenv("IB_USER")
	//password := os.Getenv("IB_PWD")

	_, err := clientv1.AuthenticateInfobaseHandler(ctx, channel, endpoint, &messagesv1.AuthenticateInfobaseRequest{
		ClusterId: clusterId,
		User:      user,
		Password:  password,
	}, nil)

	if err != nil {
		fmt.Println(err)
	}

	return next()
}
