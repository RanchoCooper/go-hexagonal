package repository

import (
	"context"
	"fmt"
	builtinLog "log"
	"os"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/spf13/cast"
	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

type MySQL struct {
	db *gorm.DB
}

func NewMySQLClient() *MySQL {
	db, err := openGormDB()
	if err != nil {
		panic(err)
	}
	return &MySQL{db: db}
}

func (c *MySQL) GetDB(ctx context.Context) *gorm.DB {
	return c.db.WithContext(ctx)
}

func (c *MySQL) SetDB(db *gorm.DB) {
	c.db = db
}

func (c *MySQL) Close(ctx context.Context) error {
	sqlDB, err := c.GetDB(ctx).DB()
	if err != nil {
		log.SugaredLogger.Errorf("get MySQL DB fail. err: %v", err)
		return err
	}
	if sqlDB != nil {
		if err := sqlDB.Close(); err != nil {
			log.SugaredLogger.Errorf("close MySQL fail. err: %v", err)
			return err
		}
	}

	log.Logger.Info("MySQL closed")
	return nil
}

func (c *MySQL) MockClient() (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic("mock MySQL fail, err: " + err.Error())
	}

	dialect := driver.New(
		driver.Config{
			Conn:                      sqlDB,
			DriverName:                "mysql-mock",
			SkipInitializeWithVersion: true,
		},
	)

	c.db, err = gorm.Open(dialect, buildGormConfig())

	return c.db, mock
}

func openGormDB() (*gorm.DB, error) {
	var (
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
			config.GlobalConfig.MySQL.User,
			config.GlobalConfig.MySQL.Password,
			config.GlobalConfig.MySQL.Host,
			config.GlobalConfig.MySQL.Database,
			config.GlobalConfig.MySQL.CharSet,
			config.GlobalConfig.MySQL.ParseTime,
			config.GlobalConfig.MySQL.TimeZone,
		)
		dialect = driver.New(driver.Config{
			DSN:                       dsn,
			DriverName:                "mysql",
			SkipInitializeWithVersion: true,
		})
	)

	db, err := gorm.Open(dialect, buildGormConfig())
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(config.GlobalConfig.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.GlobalConfig.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cast.ToDuration(config.GlobalConfig.MySQL.MaxLifeTime))
	sqlDB.SetConnMaxIdleTime(cast.ToDuration(config.GlobalConfig.MySQL.MaxIdleTime))

	return db, nil
}

func buildGormConfig() *gorm.Config {
	logger := gormLogger.New(
		builtinLog.New(os.Stdout, "\r\n", builtinLog.LstdFlags),
		gormLogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormLogger.Info, // Log level
			IgnoreRecordNotFoundError: false,           // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,            // Disable color
		},
	)

	return &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         logger,
	}
}
