package md

import (
	"context"
)

func AnnotateContext(ctx context.Context, metadataAnnotators []AnnotationHandler, req interface{}, options ...AnnotateContextOption) (context.Context, error) {
	ctx, md, err := annotateContext(ctx, metadataAnnotators, req, options...)
	if err != nil {
		return nil, err
	}
	if md == nil {
		return ctx, nil
	}

	return NewRequestMetadataContext(ctx, md), nil
}

func annotateContext(ctx context.Context, metadataAnnotators []AnnotationHandler, req interface{}, options ...AnnotateContextOption) (context.Context, RequestMetadata, error) {

	for _, o := range options {
		ctx = o(ctx)
	}

	var pairs []string

	type getClusterId interface {
		GetClusterID() string
	}

	if tReq, ok := req.(getClusterId); ok {
		clusterId := tReq.GetClusterID()
		if len(clusterId) > 0 {
			pairs = append(pairs, "cluster-id", clusterId)
		}
	}

	// timeout := DefaultContextTimeout
	// if tm := req.Header.Get(metadataGrpcTimeout); tm != "" {
	// 	var err error
	// 	timeout, err = timeoutDecode(tm)
	// 	if err != nil {
	// 		return nil, nil, status.Errorf(codes.InvalidArgument, "invalid grpc-timeout: %s", tm)
	// 	}
	// }
	//

	md := Pairs(pairs...)
	for _, mda := range metadataAnnotators {
		md = Join(md, mda(ctx, req))
	}
	return ctx, md, nil
}

type requestMetadataKey struct{}

// NewRequestMetadataContext creates a new context with ServerMetadata
func NewRequestMetadataContext(ctx context.Context, md RequestMetadata) context.Context {
	return context.WithValue(ctx, requestMetadataKey{}, md)
}

// RequestMetadataFromContext returns the ServerMetadata in ctx
func RequestMetadataFromContext(ctx context.Context) (md RequestMetadata, ok bool) {
	md, ok = ctx.Value(requestMetadataKey{}).(RequestMetadata)
	return
}

func ExtractMetadata(ctx context.Context) RequestMetadata {
	md, ok := RequestMetadataFromContext(ctx)
	if !ok {
		return Pairs()
	}
	return md
}
