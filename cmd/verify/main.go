package main

import (
	"fmt"
	"os"

	"github.com/benjamin-wright/auth-server/cmd/verify/internal/handlers"
	"github.com/benjamin-wright/auth-server/internal/api"
	"github.com/benjamin-wright/auth-server/internal/tokens"
)

func main() {
	loginURL := os.Getenv("LOGIN_URL")

	verifier, err := tokens.NewVerifier("/etc/auth-server/certs/signing.crt")
	if err != nil {
		panic(fmt.Errorf("failed to create token issuer: %+v", err))
	}

	api.Run(api.Router(api.RunOptions{
		Handlers: []api.Handler{
			handlers.Auth(verifier, loginURL),
		},
	}))
}
