package login

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/common"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

//go:embed get.html
var loginContent string

type GetLoginData struct {
	Nonce string
	Error string
}

func Get(prefix string, domain string, client *redis.Client) api.Handler {
	t, err := common.New(loginContent)
	if err != nil {
		panic(fmt.Errorf("failed to create login template: %+v", err))
	}

	return api.Handler{
		Method: "GET",
		Path:   fmt.Sprintf("%s/login", prefix),
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

			err = t.Execute(c.Writer, common.RenderData{
				Common: common.CommonData{
					Prefix:       prefix,
					Domain:       domain,
					Title:        "Login",
					RegisterLink: true,
				},
				Context: GetLoginData{
					Nonce: uuid.String(),
					Error: "",
				},
			})
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to render login page: %+v", err))
			}
		},
	}
}
