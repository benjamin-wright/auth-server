package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/benjamin-wright/auth-server/cmd/tokens/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/tokens"
	"github.com/gin-gonic/gin"
)

func ValidateLoginToken(t tokens.Keyfile) api.Handler {
	return api.Handler{
		Method: "GET",
		Path:   "/validate/login",
		Handler: func(c *gin.Context) {
			token := c.Query("token")

			if token == "" {
				c.AbortWithError(http.StatusBadRequest, errors.New("called without a token"))
				return
			}

			claims, err := t.Parse(token)
			if err != nil {
				c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("failed to parse token: %+v", err))
				return
			}

			if claims.Kind != "login" {
				c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("wrong kind of token, expected login and got %s", claims.Kind))
				return
			}

			c.JSON(http.StatusOK, client.ValidateLoginTokenResponse{
				Subject: claims.Subject,
			})
		},
	}
}
