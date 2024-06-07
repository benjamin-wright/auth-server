package register

import (
	_ "embed"
	"fmt"

	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/common"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/sut"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
)

//go:embed get.html
var registerContent string

type GetRegisterData struct {
	SUT   string
	Error string
}

func Get(prefix string, domain string, suts *sut.Client) api.Handler {
	t, err := common.New(registerContent)
	if err != nil {
		panic(fmt.Errorf("failed to create register template: %+v", err))
	}

	return api.Handler{
		Method: "GET",
		Path:   fmt.Sprintf("%s/register", prefix),
		Handler: func(c *gin.Context) {
			uuid, err := suts.Get()
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to get sut: %+v", err))
				return
			}

			err = t.Execute(c.Writer, common.RenderData{
				Common: common.CommonData{
					Prefix:    prefix,
					Domain:    domain,
					Title:     "Register",
					LoginLink: true,
				},
				Context: GetRegisterData{
					SUT:   uuid,
					Error: "",
				},
			})
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to render register page: %+v", err))
			}
		},
	}
}
