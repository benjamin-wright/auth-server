package server

import (
	"embed"
	"os"

	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/pages/admin"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/pages/login"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/pages/logout"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server/pages/register"
	tokenClient "github.com/benjamin-wright/auth-server/cmd/tokens/pkg/client"
	userClient "github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

//go:embed static
var staticContent embed.FS

func Router(prefix string, domain string, rdb *redis.Client, tokens *tokenClient.Client, users *userClient.Client) *gin.Engine {
	options := api.RunOptions{
		Handlers: []api.Handler{
			admin.Get(prefix, domain, rdb, users),
			login.Get(prefix, domain, rdb),
			login.Post(prefix, domain, rdb, tokens, users),
			register.Get(prefix, domain, rdb),
			register.Post(prefix, domain, rdb, tokens, users),
			logout.Get(prefix, domain),
		},
	}

	stat, err := os.Stat("/www/static")
	if err == nil && stat.IsDir() {
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

	log.Info().Msg("using default files from embedded content")
	options.EmbeddedFileHandlers = []api.EmbeddedFileHandler{
		{
			Path:   prefix + "/static",
			FSPath: "static",
			FS:     staticContent,
		},
	}

	return api.Router(options)
}
