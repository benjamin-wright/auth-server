package main

import (
	"os"

	"github.com/benjamin-wright/auth-server/cmd/tokens/internal/handlers"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/tokens"
)

func main() {
	keyfile := tokens.Keyfile(os.Getenv("TOKEN_KEYFILE"))

	api.Run(api.Router(api.RunOptions{
		Handlers: []api.Handler{
			handlers.GetLoginToken(keyfile),
			handlers.GetAdminToken(keyfile),
		},
	}))
}
