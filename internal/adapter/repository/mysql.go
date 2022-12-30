package repository

import (
	"context"
	"fmt"
	buitinLog "log"
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

/**
 * @author Rancho
 * @date 2021/12/21
 */

func buildGormConfig() *gorm.Config {
	logger := gormLogger.New(
		buitinLog.New(os.Stdout, "\r\n", buitinLog.LstdFlags),
		gormLogger.Config{
			SlowThreshold:             time.Second,     // Slow SQL threshold
			LogLevel:                  gormLogger.Info, // Log level
			IgnoreRecordNotFoundError: false,           // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,            // Disable color
		},
	)
	// logger := zapgorm2.New(log.Logger)
	// logger.SetAsDefault()
	// logger.LogMode(gormLogger.Info)

	return &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         logger,
	}
}

type MySQL struct {
	db *gorm.DB
}

func (c *MySQL) GetDB(ctx context.Context) *gorm.DB {
	return c.db.WithContext(ctx)
}

func (c *MySQL) SetDB(db *gorm.DB) {
	c.db = db
}

func (c *MySQL) Close(ctx context.Context) {
	sqlDB, _ := c.GetDB(ctx).DB()
	if sqlDB != nil {
		err := sqlDB.Close()
		if err != nil {
			log.SugaredLogger.Errorf("close MySQL fail. err: %v", err)
		}
	}
	log.Logger.Info("MySQL closed")
}

func (c *MySQL) MockClient() (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic("mock MySQL fail, err: " + err.Error())
	}
	dialect := driver.New(driver.Config{
		Conn:                      sqlDB,
		DriverName:                "mysql-mock",
		SkipInitializeWithVersion: true,
	})

	c.db, err = gorm.Open(dialect, buildGormConfig())

	return c.db, mock
}

func openGormDB() (*gorm.DB, error) {
	var (
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
			config.Config.MySQL.User,
			config.Config.MySQL.Password,
			config.Config.MySQL.Host,
			config.Config.MySQL.Database,
			config.Config.MySQL.CharSet,
			config.Config.MySQL.ParseTime,
			config.Config.MySQL.TimeZone,
		)
		dialect = driver.New(driver.Config{
			DSN:                       dsn,
			DriverName:                "mysql",
			DefaultStringSize:         255,
			SkipInitializeWithVersion: true,
			// ServerVersion:                 "",
			// DSNConfig:                     nil,
			// Conn:                          nil,
			// DefaultDatetimePrecision:      nil,
			// DisableWithReturning:          false,
			// DisableDatetimePrecision:      false,
			// DontSupportRenameIndex:        false,
			// DontSupportRenameColumn:       false,
			// DontSupportForShareClause:     false,
			// DontSupportNullAsDefaultValue: false,
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
	sqlDB.SetMaxIdleConns(config.Config.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.Config.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cast.ToDuration(config.Config.MySQL.MaxLifeTime))
	sqlDB.SetConnMaxIdleTime(cast.ToDuration(config.Config.MySQL.MaxIdleTime))

	return db, nil
}

func NewMySQLClient() *MySQL {
	db, err := openGormDB()
	if err != nil {
		panic(err)
	}
	return &MySQL{db: db}
}
