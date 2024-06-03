package tests

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go-hexagonal/config"
)

func SetupPostgreSQL(t *testing.T) (mysqlConf *config.PostgresDBConf) {

	t.Log("Setting up an instance of PostGreSQL with testcontainers-go")

	ctx := context.TODO()

	user, password, database := "postgres", "123456", "postgres"

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"MySQL_User":     user,
			"MySQL_Password": password,
			"MySQL_Database": database,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForExec([]string{"pg is ready"}).WithPollInterval(1*time.Second).WithExitCodeMatcher(func(exitCode int) bool {
				return exitCode == 0
			}),
			wait.ForLog("database system is ready to accept connections"),
		),
	}

	pg, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("could not start Docker container, err: %s", err)
	}

	// clean up the container after the test is complete
	t.Cleanup(func() {
		t.Log("Removing pg container from Docker")
		if err := pg.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pg container, err: %s", err)
		}
	})

	host, err := pg.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get host where the container host is exposed, err: %s", err)
	}

	port, err := pg.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("failed to get externally mapped port to pg database, err: %s", err)
	}

	t.Log("Got connection port to PostgreSQL: ", port)

	return &config.PostgresDBConf{
		Host:     host,
		Port:     port.Int(),
		Username: "postgres",
		Password: "123456",
		DbName:   "postgres",
		SSLMode:  "disable",
		TimeZone: "UTC",
	}
}
