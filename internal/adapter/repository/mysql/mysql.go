package mysql

import (
    "context"
    "fmt"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/pkg/errors"
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

type IMySQL interface {
    GetDB(ctx context.Context) *gorm.DB
    SetDB(DB *gorm.DB)
    Close(ctx context.Context)
    MockClient() (*gorm.DB, sqlmock.Sqlmock)
}

type client struct {
    db *gorm.DB
}

func (c *client) GetDB(ctx context.Context) *gorm.DB {
    return c.db.WithContext(ctx)
}

func (c *client) SetDB(DB *gorm.DB) {
    c.db = DB
}

func (c *client) Close(ctx context.Context) {
    sqlDB, _ := c.GetDB(ctx).DB()
    if sqlDB != nil {
        err := sqlDB.Close()
        if err != nil {
            log.SugaredLogger.Errorf("close mysql client fail. err: %v", err)
        }
    }
    log.Logger.Info("mysql client closed")
}

func (c *client) MockClient() (*gorm.DB, sqlmock.Sqlmock) {
    sqlDB, mock, err := sqlmock.New()
    if err != nil {
        panic("mock MySQLClient fail, err: " + err.Error())
    }
    dialector := driver.New(driver.Config{
        Conn:       sqlDB,
        DriverName: "mysql",
    })

    // a SELECT VERSION() query will be run when gorm opens the database, so we need to expect that here
    columns := []string{"version"}
    mock.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
        mock.NewRows(columns).FromCSVString("1"),
    )
    db, err := gorm.Open(dialector, &gorm.Config{})

    return db, mock
}

func finishTransaction(err error, tx *gorm.DB) error {
    if err != nil {
        if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
            return errors.Wrap(err, rollbackErr.Error())
        }

        return err
    }

    if commitErr := tx.Commit().Error; commitErr != nil {
        return errors.Wrap(err, fmt.Sprintf("failed to commit tx, err: %v", commitErr.Error()))
    }

    return nil
}

func NewGormDB() (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
        config.Config.MySQL.User,
        config.Config.MySQL.Password,
        config.Config.MySQL.Host,
        config.Config.MySQL.Database,
        config.Config.MySQL.CharSet,
        config.Config.MySQL.ParseTime,
        config.Config.MySQL.TimeZone,
    )

    logger := zapgorm2.New(log.Logger)
    logger.SetAsDefault()
    db, err := gorm.Open(
        driver.Open(dsn),
        &gorm.Config{
            NamingStrategy: schema.NamingStrategy{SingularTable: true},
            Logger:         logger,
        },
    )
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

func NewMySQLClient() IMySQL {
    db, err := NewGormDB()
    if err != nil {
        panic(err)
    }
    return &client{db: db}
}
