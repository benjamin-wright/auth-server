package server

import (
	"net/http"

	"github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/gin-gonic/gin"
)

func getUser(c *users.Client) api.Handler {
	return api.Handler{
		Method: "GET",
		Path:   "/:name",
		Handler: func(ctx *gin.Context) {
			name := ctx.Param("name")
			user, err := c.GetUser(name)
			if err == users.ErrNoUser {
				ctx.AbortWithStatus(http.StatusNotFound)
				return
			} else if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			ctx.JSON(http.StatusOK, client.GetUserResponse{
				ID:       user.ID,
				Username: user.Name,
				Admin:    user.Admin,
			})
		},
	}
}
