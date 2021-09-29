package client

import (
	"context"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/v8platform/ras-grpc-gw/pkg/ras_client/md"
	"google.golang.org/grpc/metadata"
	"strings"
)

func AnnotateRequestMetadataGrpc() md.AnnotationHandler {
	return annotateRequestMetadataGrpc
}

func annotateRequestMetadataGrpc(ctx context.Context, _ interface{}) md.RequestMetadata {

	grpcmd, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return md.Pairs()
	}

	var pairs []string

	isReq := func(key string) bool {
		return strings.HasPrefix(key, "x-req-")
	}

	for key, values := range grpcmd {
		if !isReq(key) {
			continue
		}
		key = strings.Trim(key, "x-req-")

		switch key {
		case "cluster-id":
			var clusterId string

			for _, value := range values {

				switch len(value) {
				case 36:

					_ = uuid.MustParse(value)
					clusterId = value

					break
				default:
					// Получить кластер по имени
				}
			}
			if len(clusterId) > 0 {
				pairs = append(pairs, OverwriteClusterIdKey, clusterId)
			}

		case "cluster-auth":
			var username, password string

			for _, value := range values {
				if len(value) == 0 {
					continue
				}
				var ok bool
				username, password, ok = parseBasicAuth(value)
				if !ok {
					continue
				}
			}

			pairs = append(pairs, strings.Fields(ClusterUserKeys)[0], username)
			pairs = append(pairs, strings.Fields(ClusterPwdKeys)[0], password)

		case "infobase-auth":
			var username, password string

			for _, value := range values {
				if len(value) == 0 {
					continue
				}
				var ok bool
				username, password, ok = parseBasicAuth(value)
				if !ok {
					continue
				}
			}

			pairs = append(pairs, strings.Fields(InfobaseUserKeys)[0], username)
			pairs = append(pairs, strings.Fields(InfobasePwdKeys)[0], password)

		}

	}

	return md.Pairs(pairs...)

}

// parseBasicAuth parses an Basic Authentication string.
// "QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {

	c, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}
