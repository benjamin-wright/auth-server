package main

import (
	"os"

	"github.com/benjamin-wright/auth-server/cmd/tokens/pkg/client"
	"github.com/benjamin-wright/auth-server/cmd/verify/internal/handlers"
	"github.com/benjamin-wright/auth-server/internal/api"
)

func main() {
	loginURL := os.Getenv("LOGIN_URL")
	tokens := client.New(os.Getenv("TOKENS_API_URL"))

	api.Run(api.Router(api.RunOptions{
		Handlers: []api.Handler{
			handlers.Verify(tokens, loginURL),
		},
	}))
}
