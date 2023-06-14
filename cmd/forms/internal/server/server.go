package server

import (
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func Router(prefix string, client *redis.Client) *gin.Engine {
	return api.Router(api.RunOptions{
		Handlers: []api.Handler{
			getLogin(prefix, client),
		},
	})
}
