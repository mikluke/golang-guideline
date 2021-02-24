//+build integration

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	db  *sql.DB
	ctx = context.Background()
)

func TestMain(m *testing.M) {
	shutdown := setup()

	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() func() {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:12",
		Env:          map[string]string{"POSTGRES_PASSWORD": "root"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
	})
	if err != nil {
		logrus.WithError(err).Fatalf("failed to create container")
	}

	if err := c.Start(ctx); err != nil {
		logrus.WithError(err).Fatal("failed to start container")
	}

	host, err := c.Host(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("failed to get host")
	}

	port, err := c.MappedPort(ctx, "5432")
	if err != nil {
		logrus.WithError(err).Fatal("failed to map port")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=postgres password=root sslmode=disable", host, port.Int())

	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		logrus.WithError(err).Fatal("failed to open connection")
	}

	if err := db.Ping(context.Background()); err != nil {
		logrus.WithError(err).Fatal("failed to ping postgres")
	}

	shutdownFn := func() {
		if c != nil {
			c.Terminate(ctx)
		}
	}

	migrate()

	return shutdownFn
}

func migrate() {
	// run your migrations there
}

func cleanup(t *testing.T) {
	// refresh db state there
}

func TestPg_DoSomething(t *testing.T) {
	defer cleanup(t)

	// test your DoSomething func
}
