package server

import (
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/gin-gonic/gin"
)

func Router(cli *users.Client) *gin.Engine {
	return api.Router(api.RunOptions{
		Handlers: []api.Handler{
			postUser(cli),
			getUsers(cli),
			getUser(cli),
			putUser(cli),
			putPassword(cli),
			deleteUser(cli),
		},
	})
}
