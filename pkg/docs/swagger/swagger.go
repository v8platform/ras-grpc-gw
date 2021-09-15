package swagger

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type SwaggerOpts struct {
	Up      bool
	Path    string
	SpecURL string
}

func (o *SwaggerOpts) ServeHTTP(w http.ResponseWriter, r *http.Request, params map[string]string) {

	docs := httpSwagger.Handler(httpSwagger.URL(o.SpecURL))

	// r.RequestURI = strings.ReplaceAll(r.RequestURI, "/docs", "/docs")
	docs.ServeHTTP(w, r)

}
