package client

import (
	"context"
	clientv1 "github.com/v8platform/protos/gen/ras/client/v1"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client/cond"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client/md"
	"reflect"
	"strings"
)

type getClusterId interface {
	GetClusterId() string
}

type getInfobaseId interface {
	GetInfobaseId() string
}

const methodsNeedInfobaseAuth = "GetInfobase DropInfobase CreateInfobase UpdateInfobase" +
	"GetInfobaseSessions GetInfobaseConnections GetInfobaseLocks"

//goland:noinspection ALL
var (
	HasClusterId = cond.Func(
		func(_ context.Context, _ clientv1.Endpoint, _ *clientv1.RequestInfo, req interface{}) bool {
			_, ok := req.(getClusterId)
			return ok
		})

	HasInfobaseId = cond.Func(
		func(_ context.Context, _ clientv1.Endpoint, _ *clientv1.RequestInfo, req interface{}) bool {
			_, ok := req.(getInfobaseId)
			return ok
		})
	RequestTypeEqual = func(t interface{}) cond.Func {
		return func(_ context.Context, _ clientv1.Endpoint, _ *clientv1.RequestInfo, req interface{}) bool {
			return reflect.TypeOf(req) == reflect.TypeOf(t)
		}
	}

	IsMethodName = func(name ...string) cond.Func {
		return func(_ context.Context, _ clientv1.Endpoint, info *clientv1.RequestInfo, _ interface{}) bool {
			methodName := strings.ToLower(info.Method)
			for _, s := range name {
				if strings.EqualFold(s, methodName) {
					return true
				}
			}
			return false
		}
	}
	IsMethodNames = func(name string) cond.Func {
		names := strings.Fields(name)

		return func(_ context.Context, _ clientv1.Endpoint, info *clientv1.RequestInfo, _ interface{}) bool {
			methodName := strings.ToLower(info.Method)
			for _, s := range names {
				if strings.EqualFold(s, methodName) {
					return true
				}
			}
			return false
		}
	}

	IsEndpointRequest = cond.Func(
		func(_ context.Context, endpoint clientv1.Endpoint, _ *clientv1.RequestInfo, _ interface{}) bool {
			return endpoint != nil
		})

	NeedClusterAuth  = cond.Not{IsMethodName("GetClusters")}
	NeedInfobaseAuth = IsMethodNames(methodsNeedInfobaseAuth)
)

func HasMdValues(keys ...string) cond.Func {
	return func(ctx context.Context, endpoint clientv1.Endpoint, info *clientv1.RequestInfo, req interface{}) bool {
		if len(keys) == 0 {
			return true
		}

		reqMd := md.ExtractMetadata(ctx)
		for _, key := range keys {

			ks := strings.Fields(key)
			var has bool

			for _, s := range ks {
				if reqMd.Has(s) {
					has = true
					break
				}
			}

			if !has {
				return false
			}
		}
		return true
	}
}
