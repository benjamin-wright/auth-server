package logout

import (
	"fmt"

	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
)

func Get(prefix string, domain string) api.Handler {
	return api.Handler{
		Method: "GET",
		Path:   fmt.Sprintf("%s/logout", prefix),
		Handler: func(c *gin.Context) {
			c.SetCookie("ponglehub.login", "", -1, "", domain, false, true)
			c.Redirect(302, "http://"+domain+"/"+prefix+"/login")
		},
	}
}
