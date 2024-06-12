package invite

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/common"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/sut"
	userClient "github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/password"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:embed post.html
var inviteResponseContent string

type InviteForm struct {
	SUT      string `form:"sut" binding:"required"`
	Username string `form:"username" binding:"required"`
}

type InviteResponse struct {
	ID       string
	Username string
	Password string
}

func Post(prefix string, domain string, suts *sut.Client, users *userClient.Client) api.Handler {
	t, err := common.New(inviteResponseContent)
	if err != nil {
		panic(fmt.Errorf("failed to create invite response template: %+v", err))
	}

	return api.Handler{
		Method: "POST",
		Path:   fmt.Sprintf("%s/admin/invite", prefix),
		Handler: func(c *gin.Context) {
			data := InviteForm{}
			err := c.Bind(&data)
			if err != nil {
				log.Error().Err(err).Msg("failed to bind form data")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/admin/invite?error=server")
				return
			}

			ok, err := suts.Check(data.SUT)
			if err != nil {
				log.Error().Err(err).Msg("failed to check sut")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/admin/invite?error=server")
				return
			}

			if !ok {
				log.Warn().Msg("sut is invalid")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/admin/invite?error=sut")
				return
			}

			pwd := password.Generate(32, true, true)

			response, err := users.AddUser(context.Background(), data.Username, pwd, false)
			if err != nil {
				if err == userClient.ErrUserExists {
					log.Warn().Err(err).Msg("user already exists")
					c.Redirect(302, "http://"+domain+"/"+prefix+"/admin/invite?error=exists")
					return
				}

				log.Error().Err(err).Msg("failed to add user")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/admin/invite?error=server")
				return
			}

			err = t.Execute(c.Writer, common.RenderData{
				Common: common.CommonData{
					Prefix: prefix,
					Domain: domain,
					Title:  "Invited User",
					Logout: true,
				},
				Context: InviteResponse{
					ID:       response.ID,
					Username: data.Username,
					Password: pwd,
				},
			})
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to render invite response page: %+v", err))
			}
		},
	}
}
