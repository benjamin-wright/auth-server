package integration

import (
	"context"
	"testing"

	"github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/stretchr/testify/assert"
)

func TestCheckPasswordIntegration(t *testing.T) {
	type checkPasswordSpec struct {
		username string
		password string
		response *client.CheckPasswordResponse
		ok       bool
		err      string
	}

	test(t, []testSpec[checkPasswordSpec]{
		{
			name: "Missing",
			spec: checkPasswordSpec{
				username: "myuser",
				password: "Password1!",
			},
		},
		{
			name: "Bad password",
			existing: []users.User{
				{Name: "my-user", Password: "Password3$"},
			},
			spec: checkPasswordSpec{
				username: "my-user",
				password: "Password1!",
			},
		},
		{
			name: "Success",
			existing: []users.User{
				{Name: "my-user", Password: "Password3$"},
			},
			spec: checkPasswordSpec{
				username: "my-user",
				password: "Password3$",
				response: &client.CheckPasswordResponse{
					ID:       "a uuid",
					Username: "my-user",
				},
				ok: true,
			},
		},
	}, func(u *testing.T, c *client.Client, spec checkPasswordSpec) {
		response, ok, err := c.CheckPassword(context.TODO(), spec.username, spec.password)

		if spec.err != "" {
			assert.EqualError(u, err, spec.err)
		} else {
			assert.NoError(u, err)
		}

		assert.Equal(u, spec.ok, ok)

		if spec.response != nil {
			assert.True(u, IS_UUID.MatchString(response.ID), "response.ID is not a UUID")
			assert.Equal(u, spec.response.Username, response.Username)
		} else {
			assert.Nil(u, response)
		}
	})
}
