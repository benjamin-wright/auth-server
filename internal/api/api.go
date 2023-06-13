package api

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Init() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}

	zerolog.SetGlobalLevel(level)
}

type RunOptions struct {
	Handlers []Handler
}

type Handler struct {
	Method  string
	Path    string
	Handler func(c *gin.Context)
}

func Router(options RunOptions) *gin.Engine {
	r := gin.Default()

	for _, handler := range options.Handlers {
		if handler.Path == "" {
			handler.Path = "/"
		}

		r.Handle(handler.Method, handler.Path, handler.Handler)
	}

	return r
}

func Run(router *gin.Engine) {
	router.Run("0.0.0.0:80")
}
