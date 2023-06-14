package main

import (
	"os"

	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/db-operator/pkg/redis"
	"github.com/rs/zerolog/log"
)

func main() {
	api.Init()

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

	prefix := os.Getenv("FORMS_PREFIX")
	api.Run(server.Router(prefix, client))
}
