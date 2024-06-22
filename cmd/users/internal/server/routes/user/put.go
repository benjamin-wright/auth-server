package user

import (
	"net/http"

	"github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/gin-gonic/gin"
)

func Put(c *users.Client) api.Handler {
	return api.Handler{
		Method: "PUT",
		Path:   "/user/:id",
		Handler: func(ctx *gin.Context) {
			id := ctx.Param("id")
			if id == "" {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}

			var user client.UpdateUserRequest
			if err := ctx.BindJSON(&user); err != nil {
				ctx.AbortWithStatus(http.StatusBadRequest)
				return
			}

			if user.Password != "" {
				err := c.SetPassword(id, user.Password)
				if err != nil {
					ctx.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				ctx.Status(http.StatusNoContent)
				return
			}

			err := c.SetAdmin(id, user.Admin)
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			ctx.Status(http.StatusNoContent)
		},
	}
}
