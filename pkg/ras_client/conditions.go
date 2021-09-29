package client

import clientv1 "github.com/v8platform/protos/gen/ras/client/v1"

type getClusterId interface {
	GetClusterID() string
}

func HasClusterId(_ clientv1.Endpoint, _ *clientv1.RequestInfo, req interface{}) bool {
	_, ok := req.(getClusterId)
	return ok
}

func IsEndpointRequest(endpoint clientv1.Endpoint, _ *clientv1.RequestInfo, _ interface{}) bool {
	return endpoint != nil
}
