package reset

import (
	_ "embed"
	"fmt"

	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/common"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/sut"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/gin-gonic/gin"
)

//go:embed get.html
var resetContent string

type GetResetData struct {
	SUT   string
	Error string
}

func Get(prefix string, domain string, suts *sut.Client, users *users.Client) api.Handler {
	t, err := common.New(resetContent)
	if err != nil {
		panic(fmt.Errorf("failed to create login template: %+v", err))
	}

	return api.Handler{
		Method: "GET",
		Path:   fmt.Sprintf("%s/login", prefix),
		Handler: func(c *gin.Context) {
			uuid, err := suts.Get()
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to generate SUT: %+v", err))
				return
			}

			err = t.Execute(c.Writer, common.RenderData{
				Common: common.CommonData{
					Prefix: prefix,
					Domain: domain,
					Title:  "Reset Password",
				},
				Context: GetResetData{
					SUT:   uuid,
					Error: "",
				},
			})
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to render login page: %+v", err))
			}
		},
	}
}
