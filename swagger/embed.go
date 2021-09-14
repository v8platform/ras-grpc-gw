package swagger

import (
	"embed"
)

//go:embed docs.swagger.json
var SwaggerDocs embed.FS
