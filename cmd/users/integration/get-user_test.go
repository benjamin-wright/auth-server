package integration

import (
	"context"
	"testing"

	"github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/stretchr/testify/assert"
)

func TestGetUserIntegration(t *testing.T) {
	type getUserSpec struct {
		user     string
		response *client.GetUserResponse
		err      string
	}

	test(t, []testSpec[getUserSpec]{
		{
			name: "Empty",
			spec: getUserSpec{
				user: "myuser",
				err:  "failed with status code 404",
			},
		},
		{
			name: "Found",
			existing: []users.User{
				{Name: "myuser", Password: "Password1?"},
			},
			spec: getUserSpec{
				user: "myuser",
				response: &client.GetUserResponse{
					Username: "myuser",
				},
			},
		},
		{
			name: "Wrong",
			existing: []users.User{
				{Name: "otheruser", Password: "Password1?"},
			},
			spec: getUserSpec{
				user: "myuser",
				err:  "failed with status code 404",
			},
		},
		{
			name: "Picked",
			existing: []users.User{
				{Name: "myuser", Password: "Password1?"},
				{Name: "youruser", Password: "Password2!"},
				{Name: "diffuser", Password: "Password3@"},
			},
			spec: getUserSpec{
				user: "youruser",
				response: &client.GetUserResponse{
					Username: "youruser",
				},
			},
		},
	}, func(u *testing.T, c *client.Client, spec getUserSpec) {
		response, err := c.GetUser(context.TODO(), spec.user)

		if spec.response != nil {
			assert.Equal(u, spec.response.Username, response.Username)
			assert.True(u, IS_UUID.MatchString(response.ID), "response.ID is not a UUID")
		} else {
			assert.Nil(u, response)
		}

		if spec.err != "" {
			assert.EqualError(u, err, spec.err)
		} else {
			assert.NoError(u, err)
		}
	})
}
