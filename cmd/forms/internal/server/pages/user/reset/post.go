package reset

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/common"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/sut"
	userClient "github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:embed post.html
var resetResponseContent string

type ResetForm struct {
	SUT             string `form:"sut" binding:"required"`
	Password        string `form:"password" binding:"required"`
	ConfirmPassword string `form:"confirm" binding:"required"`
}

func Post(prefix string, domain string, suts *sut.Client, users *userClient.Client) api.Handler {
	t, err := common.New(resetResponseContent)
	if err != nil {
		panic(fmt.Errorf("failed to create reset response template: %+v", err))
	}

	return api.Handler{
		Method: "POST",
		Path:   fmt.Sprintf("%s/user/reset", prefix),
		Handler: func(c *gin.Context) {
			callingUser := c.Request.Header.Get("x-auth-user")
			if callingUser == "" {
				c.AbortWithStatus(401)
				return
			}

			data := ResetForm{}
			err := c.Bind(&data)
			if err != nil {
				log.Error().Err(err).Msg("failed to bind form data")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/reset?error=server")
				return
			}

			ok, err := suts.Check(data.SUT)
			if err != nil {
				log.Error().Err(err).Msg("failed to check sut")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/reset?error=server")
				return
			}

			if !ok {
				log.Warn().Msg("sut is invalid")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/reset?error=sut")
				return
			}

			if data.Password == "" {
				log.Warn().Msg("password is empty")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/reset?error=password")
				return
			}

			if data.Password != data.ConfirmPassword {
				log.Warn().Msg("passwords do not match")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/reset?error=password")
				return
			}

			err = users.UpdateUser(context.Background(), callingUser, false, data.Password)
			if err != nil {
				log.Error().Err(err).Msg("failed to update user password")
				c.Redirect(302, "http://"+domain+"/"+prefix+"/reset?error=server")
				return
			}

			err = t.Execute(c.Writer, common.RenderData{
				Common: common.CommonData{
					Prefix: prefix,
					Domain: domain,
					Title:  "Password Reset",
					Logout: true,
				},
			})
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to render invite response page: %+v", err))
			}
		},
	}
}
