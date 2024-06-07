package main

import (
	"os"

	"github.com/benjamin-wright/auth-server/cmd/forms/internal/server"
	"github.com/benjamin-wright/auth-server/cmd/forms/internal/sut"
	tokenClient "github.com/benjamin-wright/auth-server/cmd/tokens/pkg/client"
	usersClient "github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/rs/zerolog/log"
)

func main() {
	api.Init()

	domain := os.Getenv("AUTH_DOMAIN")
	prefix := os.Getenv("FORMS_PREFIX")
	tokens := tokenClient.New(os.Getenv("TOKENS_URL"))
	users := usersClient.New(os.Getenv("USERS_URL"))
	suts, err := sut.New()
	if err != nil {
		log.Error().Err(err).Msg("Error getting suts client")
		os.Exit(1)
	}

	api.Run(server.Router(prefix, domain, tokens, users, suts))
}
