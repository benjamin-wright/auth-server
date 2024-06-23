package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/benjamin-wright/auth-server/cmd/tokens/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/tokens"
	"github.com/gin-gonic/gin"
)

func GetLoginToken(t *tokens.TokenIssuer) api.Handler {
	return api.Handler{
		Method: "GET",
		Path:   "/:subject/login",
		Handler: func(c *gin.Context) {
			subject := c.Param("subject")

			loginToken, err := t.New(subject, []string{"login"}, time.Hour)
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to create login token: %+v", err))
				return
			}

			c.JSON(http.StatusOK, client.GetLoginTokenResponse{
				Token:  loginToken,
				MaxAge: int(time.Hour / time.Second),
			})
		},
	}
}
