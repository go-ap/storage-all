package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	pg "github.com/go-ap/storage-pg"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupContainer(t testing.TB) pg.Config {
	ctx := context.Background()

	if dockerHost := os.Getenv("DOCKER_HOST"); dockerHost == "" {
		t.Skipf("no DOCKER_HOST environment variable set to use for go-containers setup")
		return pg.Config{}
	}
	pgContainer, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithInitScripts(filepath.Join("images", "init-db.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("err getting connection string: %s", err)
	}

	pconf, err := pgx.ParseConfig(connStr)
	if err != nil {
		t.Fatalf("err getting config: %s", err)
	}

	return pg.Config{
		Host:     pconf.Host,
		Port:     pconf.Port,
		Database: pconf.Database,
		User:     pconf.User,
		Password: pconf.Password,
		ErrFn:    t.Logf,
	}
}
