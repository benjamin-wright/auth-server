package admin

import (
	_ "embed"
	"fmt"

	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/common"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/sut"
	userClient "github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
)

//go:embed get.html
var registerContent string

type User struct {
	ID       string
	Username string
	Admin    bool
	Self     bool
}

type GetListUsersData struct {
	Caller string
	SUT    string
	Users  []User
	Error  string
}

func Get(prefix string, domain string, suts *sut.Client, usersClient *userClient.Client) api.Handler {
	t, err := common.New(registerContent)
	if err != nil {
		panic(fmt.Errorf("failed to create register template: %+v", err))
	}

	return api.Handler{
		Method: "GET",
		Path:   fmt.Sprintf("%s/admin", prefix),
		Handler: func(c *gin.Context) {
			callingUser := c.Request.Header.Get("x-auth-user")
			if callingUser == "" {
				c.AbortWithStatus(401)
				return
			}

			uuid, err := suts.Get()
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to generate SUT: %+v", err))
				return
			}

			resp, err := usersClient.ListUsers(c.Request.Context())
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to list users: %+v", err))
				return
			}

			users := make([]User, len(resp.Users))
			for i, user := range resp.Users {
				users[i] = User{
					ID:       user.ID,
					Username: user.Username,
					Admin:    user.Admin,
					Self:     false,
				}
			}

			err = t.Execute(c.Writer, common.RenderData{
				Common: common.CommonData{
					Prefix: prefix,
					Domain: domain,
					Title:  "Admin",
					Logout: true,
				},
				Context: GetListUsersData{
					Caller: callingUser,
					SUT:    uuid,
					Users:  users,
					Error:  "",
				},
			})
			if err != nil {
				c.AbortWithError(500, fmt.Errorf("failed to render admin page: %+v", err))
			}
		},
	}
}
