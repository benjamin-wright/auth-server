package user

import (
	"net/http"

	"github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/gin-gonic/gin"
)

func List(c *users.Client) api.Handler {
	return api.Handler{
		Method: "GET",
		Path:   "/user",
		Handler: func(ctx *gin.Context) {
			users, err := c.ListUsers()
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			response := client.ListUsersResponse{
				Users: make([]client.GetUserResponse, len(users)),
			}

			for i, user := range users {
				response.Users[i] = client.GetUserResponse{
					ID:       user.ID,
					Username: user.Name,
					Admin:    user.Admin,
				}
			}

			ctx.JSON(http.StatusOK, response)
		},
	}
}
