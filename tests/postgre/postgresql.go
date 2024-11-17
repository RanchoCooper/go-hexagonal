package postgre

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go-hexagonal/config"
	"go-hexagonal/internal/adapter/repository/postgre"
	"go-hexagonal/tests/postgre/migrations/migrate"
)

func SetupPostgreSQL(t *testing.T) (postgreSQLConfig *config.PostgreSQLConfig) {

	t.Log("Setting up an instance of PostgreSQL with testcontainers-go")
	ctx := context.TODO()

	user, dbName, password := "postgres", "postgres", "123456"

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_DATABASE": dbName,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForExec([]string{"pg_isready"}).
				WithPollInterval(1*time.Second).
				WithExitCodeMatcher(func(exitCode int) bool {
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
		t.Fatalf("Could not start Docker container: %s", err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		t.Log("Removing pg container from Docker")
		if err := pg.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pg container: %s", err)
		}
	})

	host, err := pg.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get host where the container host is exposed: %s", err)
	}

	port, err := pg.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("Failed to get externally mapped port to pg database: %s", err)
	}
	t.Log("Got connection port to PostgreSQL: ", port)

	return &config.PostgreSQLConfig{
		Host:     host,
		Port:     port.Int(),
		Username: user,
		Password: password,
		DbName:   dbName,
		SSLMode:  "disable",
		TimeZone: "UTC",
	}
}

func MockPgSQLData(t *testing.T, conf *config.GlobalConfig, sqls []string) {

	err := migrate.GolangMigrateDrop(conf)
	if err != nil {
		t.Fatalf("GolangMigrateDrop fail %+v\n", err)
	}

	err = migrate.GolangMigrateUp(conf)
	if err != nil {
		t.Fatalf("GolangMigrateUp fail %+v\n", err)
	}

	pgPool, err := postgre.NewConnPool(conf.Postgre)
	if err != nil {
		t.Fatalf("NewConnPool fail %+v\n", err)
	}

	err = pgPool.Ping(context.Background())
	if err != nil {
		t.Fatalf("pgPool Ping fail %+v\n", err)
	}

	for _, sql := range sqls {
		_, err = pgPool.Exec(context.Background(), sql)
		if err != nil {
			t.Fatalf("Unable to insert data %+v\n", err)
		}
	}
}
