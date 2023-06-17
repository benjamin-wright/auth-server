package api

import (
	"embed"
	"io/fs"
	"net/http"
)

type FileHandler struct {
	Path   string
	FSPath string
	FS     embed.FS
}

func (f *FileHandler) GetHttpFS() http.FileSystem {
	sub, err := fs.Sub(f.FS, f.FSPath)

	if err != nil {
		panic(err)
	}

	return http.FS(sub)
}
