package repository

import (
	"context"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/spf13/cast"
	driver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"

	"go-hexagonal/config"
	"go-hexagonal/util/log"
)

/**
 * @author Rancho
 * @date 2021/12/21
 */

var (
	logger     = zapgorm2.New(log.Logger)
	gormConfig = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         logger,
	}
)

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
			log.SugaredLogger.Errorf("close mysql client fail. err: %v", err)
		}
	}
	log.Logger.Info("mysql client closed")
}

func (c *MySQL) MockClient() (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic("mock MySQL fail, err: " + err.Error())
	}
	dialector := driver.New(driver.Config{
		Conn:                      sqlDB,
		DriverName:                "mysql-mock",
		SkipInitializeWithVersion: true,
	})

	c.db, err = gorm.Open(dialector, gormConfig)

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
		dialector = driver.New(driver.Config{
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

	db, err := gorm.Open(dialector, gormConfig)

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
