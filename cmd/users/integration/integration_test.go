package integration

import (
	"flag"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/benjamin-wright/auth-server/cmd/users/internal/server"
	"github.com/benjamin-wright/auth-server/cmd/users/pkg/client"
	"github.com/benjamin-wright/auth-server/internal/users"
	"github.com/benjamin-wright/db-operator/pkg/test/cockroach"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		m.Run()
		return
	}

	if testing.Verbose() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	close := cockroach.Run("cockroach", 26257)
	defer close()

	cockroach.Migrate("../../../deploy/chart/resources/users-1.sql")

	m.Run()
}

type testSpec[T any] struct {
	name     string
	existing []users.User
	spec     T
}

func test[T any](t *testing.T, specs []testSpec[T], f func(t *testing.T, c *client.Client, spec T)) {
	if testing.Short() {
		t.SkipNow()
	}

	cli, err := users.New()
	if !assert.NoError(t, err) {
		return
	}

	r := server.Router(cli)
	srv := httptest.NewServer(r.Handler())
	c := client.New(srv.URL)

	for _, spec := range specs {
		t.Run(spec.name, func(u *testing.T) {
			if !assert.NoError(u, cli.DeleteAllUsers()) {
				return
			}

			for _, user := range spec.existing {
				if _, err := cli.AddUser(user); !assert.NoError(t, err) {
					return
				}
			}

			f(u, c, spec.spec)
		})
	}
}
