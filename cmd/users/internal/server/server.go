package server

import (
	"github.com/benjamin-wright/auth-server/cmd/users/internal/server/routes/user"
	"github.com/benjamin-wright/auth-server/cmd/users/internal/server/routes/validate"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/gin-gonic/gin"
)

func Router(cli *users.Client) *gin.Engine {
	return api.Router(api.RunOptions{
		Handlers: []api.Handler{
			user.Delete(cli),
			user.Get(cli),
			user.List(cli),
			user.Post(cli),
			user.Put(cli),
			validate.Put(cli),
		},
	})
}
