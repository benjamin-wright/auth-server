package main

import (
	"fmt"

	"github.com/benjamin-wright/auth-server/cmd/tokens/internal/handlers"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/tokens"
)

func main() {
	issuer, err := tokens.NewIssuer("/etc/auth-server/certs/signing.key")
	if err != nil {
		panic(fmt.Errorf("failed to create token issuer: %+v", err))
	}

	api.Run(api.Router(api.RunOptions{
		Handlers: []api.Handler{
			handlers.GetLoginToken(issuer),
			handlers.GetAdminToken(issuer),
		},
	}))
}
