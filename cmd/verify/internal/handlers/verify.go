package handlers

import (
	"net/http"

	"github.com/benjamin-wright/auth-server/cmd/tokens/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Verify(tokens *client.Client, loginURL string) api.Handler {
	return api.Handler{
		Method: "GET",
		Path:   "/",
		Handler: func(c *gin.Context) {
			token, err := c.Cookie("ponglehub.login")
			if err != nil {
				log.Info().Err(err).Msg("no token found in cookie")
				c.Redirect(http.StatusTemporaryRedirect, loginURL)
				return
			}

			res, err := tokens.ValidateLoginToken(token)
			if err != nil {
				log.Info().Err(err).Msg("token validation failed")
				c.Redirect(http.StatusTemporaryRedirect, loginURL)
				return
			}

			c.Header("x-auth-user", res.Subject)
			c.Status(http.StatusOK)
		},
	}
}
