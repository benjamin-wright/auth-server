package admin

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/common"
	userClient "github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

//go:embed get.html
var registerContent string

type GetListUsersData struct {
	Nonce string
	Users []userClient.GetUserResponse
	Error string
}

func Get(prefix string, domain string, client *redis.Client, usersClient *userClient.Client) api.Handler {
	t, err := common.New(registerContent)
	if err != nil {
		panic(fmt.Errorf("failed to create register template: %+v", err))
	}

	return api.Handler{
		Method: "GET",
		Path:   fmt.Sprintf("%s/admin", prefix),
		Handler: func(c *gin.Context) {
			uuid, err := uuid.NewRandom()
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to generate nonce: %+v", err))
				return
			}

			cmd := client.Set(context.Background(), uuid.String(), true, 5*time.Minute)
			if cmd.Err() != nil {
				c.AbortWithError(500, fmt.Errorf("failed to set nonce: %+v", cmd.Err()))
				return
			}

			resp, err := usersClient.ListUsers(c.Request.Context())
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to list users: %+v", err))
				return
			}

			err = t.Execute(c.Writer, common.RenderData{
				Common: common.CommonData{
					Prefix: prefix,
					Domain: domain,
					Title:  "Admin",
				},
				Context: GetListUsersData{
					Nonce: uuid.String(),
					Users: resp.Users,
					Error: "",
				},
			})
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to render admin page: %+v", err))
			}
		},
	}
}
