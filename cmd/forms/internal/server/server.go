package server

import (
	tokenClient "github.com/benjamin-wright/auth-server/cmd/tokens/pkg/client"
	userClient "github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func Router(prefix string, domain string, rdb *redis.Client, tokens *tokenClient.Client, users *userClient.Client) *gin.Engine {
	return api.Router(api.RunOptions{
		Handlers: []api.Handler{
			getLogin(prefix, domain, rdb),
			postLogin(prefix, domain, rdb, tokens, users),
		},
	})
}
