package main

import (
	"context"
	"fmt"
	"os"

	"github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/rs/zerolog/log"
)

func main() {
	c := client.New(os.Getenv("USERS_URL"))

	err := run(c, os.Getenv("USERS_ADMIN_USERNAME"), os.Getenv("USERS_ADMIN_PASSWORD"))
	if err != nil {
		log.Error().Err(err).Msg("failed to run")
	}
}

func run(c *client.Client, user string, password string) error {
	// resp, err := c.ListUsers(context.TODO())
	// if err != nil {
	// 	return fmt.Errorf("failed to list users: %w", err)
	// }

	// if len(resp.Users) > 0 {
	// 	log.Info().Msg("Existing users found, aborting init process")
	// 	return nil
	// }

	log.Info().Msg("No users found, initializing database")

	_, err := c.AddUser(context.TODO(), user, password, true)
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}

	log.Info().Msg("User added successfully")

	return nil
}
