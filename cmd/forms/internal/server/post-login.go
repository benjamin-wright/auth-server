package server

import (
	"fmt"

	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type LoginForm struct {
	LoginNonce string `form:"login-nonce"`
	Username   string `form:"username"`
	Password   string `form:"password"`
	Confirm    string `form:"confirm"`
}

func postLogin(prefix string, client *redis.Client) api.Handler {
	return api.Handler{
		Method: "POST",
		Path:   fmt.Sprintf("%s/login", prefix),
		Handler: func(c *gin.Context) {
			data := LoginForm{}
			c.Bind(&data)
		},
	}
}
