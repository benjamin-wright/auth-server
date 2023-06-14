package server

import (
	_ "embed"
	"fmt"
	"html/template"
	"time"

	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

//go:embed login.html
var loginContent string

type GetLoginData struct {
	LoginNonce string
}

func getLogin(prefix string, client *redis.Client) api.Handler {
	t := template.New("login")
	template.Must(t.Parse(loginContent))

	return api.Handler{
		Method: "GET",
		Path:   fmt.Sprintf("%s/login", prefix),
		Handler: func(c *gin.Context) {
			uuid, err := uuid.NewRandom()
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to generate uuid: %+v", err))
				return
			}

			client.Set(c, uuid.String(), struct{}{}, 5*time.Minute)

			err = t.Execute(c.Writer, GetLoginData{
				LoginNonce: uuid.String(),
			})
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to render login page: %+v", err))
			}
		},
	}
}
