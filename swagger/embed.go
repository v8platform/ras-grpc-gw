package swagger

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed docs.swagger.json
var Docs embed.FS

func Handler() http.Handler {

	subFS, err := fs.Sub(Docs, ".")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}
	return http.FileServer(http.FS(subFS))
}

//go:generate C:/Users/khorevaa/go/go1.17.1/bin/go.exe run -ldflags "-X 'main.File=./docs.swagger.json'" ./cmd/gen.go
