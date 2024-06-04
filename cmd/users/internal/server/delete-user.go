package server

import (
	"net/http"

	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/gin-gonic/gin"
)

func deleteUser(c *users.Client) api.Handler {
	return api.Handler{
		Method: "DELETE",
		Path:   "/id/:id",
		Handler: func(ctx *gin.Context) {
			id := ctx.Param("id")
			err := c.DeleteUser(id)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			ctx.Status(http.StatusNoContent)
		},
	}
}
