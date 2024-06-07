package server

import (
	"net/http"

	"github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/gin-gonic/gin"
)

func putPassword(c *users.Client) api.Handler {
	return api.Handler{
		Method: "PUT",
		Path:   "/:name/password",
		Handler: func(ctx *gin.Context) {
			name := ctx.Param("name")
			var body client.CheckPasswordRequest
			err := ctx.BindJSON(&body)
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}

			user, err := c.CheckPassword(users.User{Name: name, Password: body.Password})
			if err == users.ErrPasswordMismatch || err == users.ErrNoUser {
				ctx.JSON(http.StatusUnauthorized, err)
				return
			} else if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			ctx.JSON(http.StatusOK, client.CheckPasswordResponse{
				Username: user.Name,
				ID:       user.ID,
				Admin:    user.Admin,
			})
		},
	}
}
