package u

import (
	"path/filepath"
	"strings"
)

/* additions to mime */

var mimeTypes = map[string]string{
	// this is a list from go's mime package
	".css":  "text/css; charset=utf-8",
	".gif":  "image/gif",
	".htm":  "text/html; charset=utf-8",
	".html": "text/html; charset=utf-8",
	".jpg":  "image/jpeg",
	".js":   "application/javascript",
	".wasm": "application/wasm",
	".pdf":  "application/pdf",
	".png":  "image/png",
	".svg":  "image/svg+xml",
	".xml":  "text/xml; charset=utf-8",

	// those are my additions
	".txt":  "text/plain",
	".exe":  "application/octet-stream",
	".json": "application/json",
}

func MimeTypeFromFileName(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	mt := mimeTypes[ext]
	if mt != "" {
		return mt
	}
	// if not given, default to this
	return "application/octet-stream"
}
