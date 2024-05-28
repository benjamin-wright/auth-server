package main

import (
	"os"

	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server"
	tokenClient "github.com/benjamin-wright/auth-server/cmd/tokens/pkg/client"
	usersClient "github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/db-operator/v2/pkg/redis"
	"github.com/rs/zerolog/log"
)

func main() {
	api.Init()

	domain := os.Getenv("AUTH_DOMAIN")
	prefix := os.Getenv("FORMS_PREFIX")
	tokens := tokenClient.New(os.Getenv("TOKENS_URL"))
	users := usersClient.New(os.Getenv("USERS_URL"))

	cfg, err := redis.ConfigFromEnv()
	if err != nil {
		log.Error().Err(err).Msg("Error getting redis config")
		return
	}

	client, err := redis.Connect(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to redis")
		return
	}

	api.Run(server.Router(prefix, domain, client, tokens, users))
}
