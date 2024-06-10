package login

import (
	"context"
	"fmt"

	"github.com/benjamin-wright/auth-server/cmd/forms/internal/sut"
	tokenClient "github.com/benjamin-wright/auth-server/cmd/tokens/pkg/client"
	userClient "github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type LoginForm struct {
	SUT      string `form:"sut" binding:"required"`
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Post(prefix string, domain string, suts *sut.Client, tokens *tokenClient.Client, users *userClient.Client) api.Handler {
	return api.Handler{
		Method: "POST",
		Path:   fmt.Sprintf("%s/login", prefix),
		Handler: func(c *gin.Context) {
			data := LoginForm{}
			err := c.Bind(&data)
			if err != nil {
				log.Error().Err(err).Msg("failed to bind form data")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/login?error=server")
				return
			}

			ok, err := suts.Check(data.SUT)
			if err != nil {
				log.Error().Err(err).Msg("failed to check SUT")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/login?error=server")
				return
			}

			if !ok {
				log.Warn().Msg("invalid SUT")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/login?error=sut")
				return
			}

			res, ok, err := users.CheckPassword(context.Background(), data.Username, data.Password)
			if err != nil {
				log.Error().Err(err).Msg("failed to check password")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/login?error=server")
				return
			}

			if !ok {
				log.Warn().Msg("incorrect user or password")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/login?error=fail")
				return
			}

			if res.Admin {
				token, err := tokens.GetAdminToken(res.ID)
				if err != nil {
					log.Error().Err(err).Msg("failed to get admin token")
					c.Redirect(302, "http://"+domain+"/"+prefix+"/login?error=server")
					return
				}
				c.SetCookie("ponglehub.login", token.Token, token.MaxAge, "", domain, false, true)
				c.Redirect(302, "http://"+domain+"/auth/admin")
			} else {
				token, err := tokens.GetLoginToken(res.ID)
				if err != nil {
					log.Error().Err(err).Msg("failed to get login token")
					c.Redirect(302, "http://"+domain+"/"+prefix+"/login?error=server")
					return
				}
				c.SetCookie("ponglehub.login", token.Token, token.MaxAge, "", domain, false, true)
				c.Redirect(302, "http://"+domain+"/home")
			}
		},
	}
}
