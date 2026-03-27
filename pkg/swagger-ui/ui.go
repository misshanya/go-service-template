package swaggerui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed *
var dist embed.FS

func Handler() (http.Handler, error) {
	f, err := fs.Sub(dist, ".")
	if err != nil {
		return nil, err
	}

	return http.FileServer(http.FS(f)), nil
}
