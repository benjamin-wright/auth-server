package main

import (
	"github.com/benjamin-wright/auth-server/cmd/users/internal/server"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/rs/zerolog/log"
)

func main() {
	api.Init()

	cli, err := users.New()
	if err != nil {
		log.Error().Err(err).Msg("Error getting database client")
	}

	api.Run(server.Router(cli))
}
