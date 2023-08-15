package register

import (
	"context"
	"fmt"

	tokenClient "github.com/benjamin-wright/auth-server/cmd/tokens/pkg/client"
	userClient "github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	usersLib "github.com/benjamin-wright/auth-server/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type RegisterForm struct {
	Nonce           string `form:"nonce" binding:"required"`
	Username        string `form:"username" binding:"required"`
	Password        string `form:"password" binding:"required"`
	ConfirmPassword string `form:"confirm-password" binding:"required"`
}

func Post(prefix string, domain string, rdb *redis.Client, tokens *tokenClient.Client, users *userClient.Client) api.Handler {
	return api.Handler{
		Method: "POST",
		Path:   fmt.Sprintf("%s/register", prefix),
		Handler: func(c *gin.Context) {
			data := RegisterForm{}
			err := c.Bind(&data)
			if err != nil {
				log.Error().Err(err).Msg("failed to bind form data")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/register?error=server")
				return
			}

			cmd := rdb.Get(c, data.Nonce)
			if cmd.Err() != nil {
				log.Warn().Err(cmd.Err()).Msgf("failed to get nonce: %s", data.Nonce)
				c.Redirect(302, "http://"+domain+"/"+prefix+"/register?error=nonce")
				return
			}

			defer func() {
				delCmd := rdb.Del(c, data.Nonce)
				if delCmd.Err() != nil {
					log.Error().Err(delCmd.Err()).Msg("failed to delete nonce")
				}
			}()

			if !usersLib.CheckPasswordComplexity(data.Password) {
				log.Warn().Msg("password does not meet complexity requirements")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/register?error=complexity")
				return
			}

			if data.Password != data.ConfirmPassword {
				log.Warn().Msg("passwords do not match")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/register?error=password")
				return
			}

			response, err := users.AddUser(context.Background(), data.Username, data.Password)
			if err != nil {
				if err == userClient.ErrUserExists {
					log.Warn().Err(err).Msg("user already exists")
					c.Redirect(302, "http://"+domain+"/"+prefix+"/register?error=exists")
					return
				}

				log.Error().Err(err).Msg("failed to add user")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/register?error=server")
				return
			}

			token, err := tokens.GetLoginToken(response.ID)
			if err != nil {
				log.Error().Err(err).Msg("failed to get login token")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/register?error=server")
				return
			}

			c.SetCookie("ponglehub.login", token.Token, token.MaxAge, "", domain, false, true)
			c.Redirect(302, "http://"+domain+"/home")
		},
	}
}
