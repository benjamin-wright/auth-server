package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/benjamin-wright/auth-server/cmd/forms/internal/sut"
	userClient "github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type PutUserRequest struct {
	Admin bool `json:"admin"`
}

func Put(prefix string, domain string, suts *sut.Client, users *userClient.Client) api.Handler {
	return api.Handler{
		Method: "PUT",
		Path:   fmt.Sprintf("%s/admin/user/:id", prefix),
		Handler: func(c *gin.Context) {
			callingUser := c.Request.Header.Get("x-auth-user")
			if callingUser == "" {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			userId := c.Param("id")
			if userId == "" {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			if callingUser == userId {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			ok, err := suts.Check(c.Query("sut"))
			if err != nil {
				log.Error().Err(err).Msg("failed to check SUT")
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			if !ok {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			var req PutUserRequest
			if err := c.BindJSON(&req); err != nil {
				log.Error().Err(err).Msg("failed to bind request")
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			err = users.UpdateUser(context.Background(), userId, req.Admin, "")
			if err != nil {
				log.Error().Err(err).Msg("failed to update user")
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.Status(http.StatusNoContent)
		},
	}
}
