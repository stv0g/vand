//go:build embed_frontend

package frontend

import (
	"embed"
)

//go:embed build
var Files embed.FS
