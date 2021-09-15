package server

import (
	"github.com/v8platform/ras-grpc-gw/swagger"
	"io/fs"
	"net/http"
)

// getOpenAPIHandler serves an OpenAPI UI.
// Adapted from https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L63
func getOpenAPIHandler() http.Handler {

	subFS, err := fs.Sub(swagger.Docs, ".")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}
	return http.FileServer(http.FS(subFS))
}

type SwaggerOpts struct {
	Up      bool
	Path    string
	SpecURL string
	Handler http.Handler
}

func (o *SwaggerOpts) ServeHTTP(w http.ResponseWriter, r *http.Request, _ map[string]string) {

	o.Handler.ServeHTTP(w, r)

}
