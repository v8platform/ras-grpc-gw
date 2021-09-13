package swagger

import (
	"embed"
)

//go:embed swagger/*
var SwaggerDocs embed.FS
