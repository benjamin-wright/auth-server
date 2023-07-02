package api

import (
	"embed"
	"io/fs"
	"net/http"
)

type EmbeddedFileHandler struct {
	Path   string
	FSPath string
	FS     embed.FS
}

func (f *EmbeddedFileHandler) getHttpFS() http.FileSystem {
	sub, err := fs.Sub(f.FS, f.FSPath)

	if err != nil {
		panic(err)
	}

	return http.FS(sub)
}
