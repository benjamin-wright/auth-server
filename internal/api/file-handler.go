package api

import (
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	Path   string
	FSPath string
	Files  []string
}

func (f *FileHandler) addRoutes(router *gin.Engine) {
	if f.Files != nil {
		for _, file := range f.Files {
			router.StaticFile(f.Path+"/"+file, f.FSPath+"/"+file)
		}
	} else {
		router.Static(f.Path, f.FSPath)
	}
}
