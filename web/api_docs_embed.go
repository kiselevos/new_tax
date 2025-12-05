package web

import (
	"embed"
)

//go:embed api_docs/swagger.json
var ApiDocsFS embed.FS
