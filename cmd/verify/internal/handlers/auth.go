package handlers

import (
	"net/http"

	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/tokens"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Auth(keyfile tokens.Keyfile, loginURL string) api.Handler {
	return api.Handler{
		Method: "GET",
		Path:   "/auth",
		Handler: func(c *gin.Context) {
			role := c.Query("role")

			token, err := c.Cookie("ponglehub.login")
			if err != nil {
				log.Info().Err(err).Msg("no token found in cookie")
				c.Redirect(http.StatusTemporaryRedirect, loginURL)
				return
			}

			if token == "" {
				log.Info().Msg("empty token found in cookie")
				c.Redirect(http.StatusTemporaryRedirect, loginURL)
				return
			}

			claims, err := keyfile.Parse(token)
			if err != nil {
				log.Info().Err(err).Msg("failed to parse token")
				c.Redirect(http.StatusTemporaryRedirect, loginURL)
				return
			}

			for _, audience := range claims.Audiences {
				if audience == role {
					c.Header("x-auth-user", claims.Subject)
					c.Status(http.StatusOK)
					return
				}
			}

			log.Info().Str("wanted", role).Strs("got", claims.Audiences).Msg("role not found in token")
			c.Redirect(http.StatusTemporaryRedirect, loginURL)
		},
	}
}
