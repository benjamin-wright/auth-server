package user

import (
	"net/http"

	"github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/gin-gonic/gin"
)

func Post(c *users.Client) api.Handler {
	return api.Handler{
		Method: "POST",
		Path:   "/user",
		Handler: func(ctx *gin.Context) {
			var body client.AddUserRequest
			err := ctx.BindJSON(&body)
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}

			id, err := c.AddUser(users.User{Name: body.Username, Password: body.Password, Admin: body.Admin})
			if err == users.ErrUserExists {
				ctx.AbortWithError(http.StatusConflict, err)
				return
			} else if err == users.ErrComplexity {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			} else if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			ctx.JSON(http.StatusCreated, client.AddUserResponse{
				ID: id,
			})
		},
	}
}
