package server

import (
	"context"
	"fmt"

	tokenClient "github.com/benjamin-wright/auth-server/cmd/tokens/pkg/client"
	userClient "github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type LoginForm struct {
	LoginNonce string `form:"login-nonce"`
	Username   string `form:"username"`
	Password   string `form:"password"`
	Confirm    string `form:"confirm"`
}

func postLogin(prefix string, domain string, rdb *redis.Client, tokens *tokenClient.Client, users *userClient.Client) api.Handler {
	return api.Handler{
		Method: "POST",
		Path:   fmt.Sprintf("%s/login", prefix),
		Handler: func(c *gin.Context) {
			data := LoginForm{}
			err := c.Bind(&data)
			if err != nil {
				log.Error().Err(err).Msg("failed to bind form data")
				c.Redirect(302, "http://"+domain+prefix+"/login?error=server")
				return
			}

			cmd := rdb.Get(c, data.LoginNonce)
			if cmd.Err() != nil {
				log.Warn().Err(cmd.Err()).Msg("failed to get nonce")
				c.Redirect(302, "http://"+domain+prefix+"/login?error=nonce")
				return
			}

			if data.Password != data.Confirm {
				log.Warn().Msg("passwords do not match")
				c.Redirect(302, "http://"+domain+prefix+"/login?error=fail")
				return
			}

			res, ok, err := users.CheckPassword(context.Background(), data.Username, data.Password)
			if err != nil {
				log.Error().Err(err).Msg("failed to check password")
				c.Redirect(302, "http://"+domain+prefix+"/login?error=server")
				return
			}

			if !ok {
				log.Warn().Msg("incorrect password")
				c.Redirect(302, "http://"+domain+prefix+"/login?error=fail")
				return
			}

			token, err := tokens.GetLoginToken(res.ID)
			if err != nil {
				log.Error().Err(err).Msg("failed to get login token")
				c.Redirect(302, "http://"+domain+prefix+"/login?error=server")
				return
			}

			delCmd := rdb.Del(c, data.LoginNonce)
			if delCmd.Err() != nil {
				log.Error().Err(delCmd.Err()).Msg("failed to delete nonce")
				c.Redirect(302, "http://"+domain+prefix+"/login?error=server")
				return
			}

			c.SetCookie("ponglehub.login", token.Token, token.MaxAge, "", domain, false, true)
			c.Redirect(302, fmt.Sprintf("%s/login", prefix))
		},
	}
}
