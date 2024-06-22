package user

import (
	"net/http"

	"github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/gin-gonic/gin"
)

func Get(c *users.Client) api.Handler {
	return api.Handler{
		Method: "GET",
		Path:   "/user/:id",
		Handler: func(ctx *gin.Context) {
			id := ctx.Param("id")
			user, err := c.GetUser(id)
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
