package server

import (
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/pages/admin"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/pages/admin/invite"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/pages/admin/user"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/pages/login"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/pages/logout"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/pages/user/reset"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/sut"
	tokenClient "github.com/benjamin-wright/auth-server/cmd/tokens/pkg/client"
	userClient "github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Router(prefix string, domain string, tokens *tokenClient.Client, users *userClient.Client, suts *sut.Client) *gin.Engine {
	options := api.RunOptions{
		Handlers: []api.Handler{
			admin.Get(prefix, domain, suts, users),
			user.Delete(prefix, domain, suts, users),
			user.Put(prefix, domain, suts, users),
			invite.Get(prefix, domain, suts),
			invite.Post(prefix, domain, suts, users),
			login.Get(prefix, domain, suts),
			login.Post(prefix, domain, suts, tokens, users),
			logout.Get(prefix, domain),
			reset.Get(prefix, domain, suts),
			reset.Post(prefix, domain, suts, users),
		},
	}

	log.Info().Msg("using static files from /www/static")
	options.StaticFiles = []api.FileHandler{
		{
			Path:   prefix + "/static",
			FSPath: "/www/static",
			Files:  []string{"styles.css", "favicon.ico"},
		},
	}

	return api.Router(options)
}
