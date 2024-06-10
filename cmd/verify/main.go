package main

import (
	"os"

	"github.com/benjamin-wright/auth-server/cmd/verify/internal/handlers"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/tokens"
)

func main() {
	loginURL := os.Getenv("LOGIN_URL")
	keyfile := tokens.Keyfile(os.Getenv("TOKEN_KEYFILE"))

	api.Run(api.Router(api.RunOptions{
		Handlers: []api.Handler{
			handlers.Auth(keyfile, loginURL),
		},
	}))
}
