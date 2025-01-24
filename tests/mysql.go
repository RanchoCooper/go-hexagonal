package tests

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-hexagonal/config"
	"go-hexagonal/tests/migrations/migrate"
)

const (
	MysqlStartTimeout = 2 * time.Minute
)

func SetupMySQL(t *testing.T) *config.MySQLConfig {
	t.Log("Setting up an instance of MySQL with testcontainers-go")
	ctx := context.Background()

	user, password, dbName := "user", "123456", "test"

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_USER":          user,
			"MYSQL_ROOT_PASSWORD": password,
			"MYSQL_PASSWORD":      password,
			"MYSQL_DATABASE":      dbName,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("3306/tcp").WithStartupTimeout(MysqlStartTimeout),
			wait.ForLog("ready for connections"),
		),
	}

	db, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("could not start Docker container, err: %s", err)
	}

	t.Cleanup(func() {
		t.Log("Removing MySQL container from Docker")
		if err := db.Terminate(ctx); err != nil {
			t.Errorf("failed to terminate MySQL container, err: %s", err) // 改为 t.Errorf
		}
	})

	host, err := db.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get host where the container is exposed, err: %s", err)
	}

	port, err := db.MappedPort(ctx, "3306/tcp")
	if err != nil {
		t.Fatalf("failed to get externally mapped port to MySQL database, err: %s", err) // 修改描述
	}

	t.Log("Got connection port to MySQL: ", port)

	return &config.MySQLConfig{
		User:      user,
		Password:  password,
		Host:      host,
		Port:      port.Int(),
		Database:  dbName,
		CharSet:   "utf8mb4",
		ParseTime: false,
		TimeZone:  "UTC",
	}
}

func MockMySQLData(t *testing.T, conf *config.Config, sqls []string) *gorm.DB {
	err := migrate.MySQLMigrateDrop(conf)
	if err != nil {
		t.Fatalf("MySQLMigrateDrop fail, err: %+v\n", err)
	}

	err = migrate.MySQLMigrateUp(conf)
	if err != nil {
		t.Fatalf("MySQLMigrateUp fail %+v\n", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		config.GlobalConfig.MySQL.User,
		config.GlobalConfig.MySQL.Password,
		config.GlobalConfig.MySQL.Host,
		config.GlobalConfig.MySQL.Port,
		config.GlobalConfig.MySQL.Database,
		config.GlobalConfig.MySQL.CharSet,
		config.GlobalConfig.MySQL.ParseTime,
		config.GlobalConfig.MySQL.TimeZone,
	)
	dialect := driver.New(driver.Config{
		DSN:                       dsn,
		DriverName:                "mysql",
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	for _, sql := range sqls {
		tx := db.Exec(sql)
		if tx.Error != nil {
			t.Fatalf("Unable to insert data %+v\n", err)
		}
	}

	return db
}
