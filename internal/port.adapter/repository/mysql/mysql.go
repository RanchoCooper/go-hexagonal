package mysql

import (
    "fmt"
    "log"
    "os"
    "sync"
    "time"

    driver "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/gorm/schema"

    "go-hexagonal/config"
)

/**
 * @author Rancho
 * @date 2021/12/21
 */

var (
    once  sync.Once
    MySQL *MySQLRepository
)

type MySQLRepository struct {
    db *gorm.DB
    // add other repository
}

func init() {
    once.Do(func() {
        var err error
        MySQL, err = NewMySQLRepository()
        if err != nil {
            panic("init MySQL fail, err: " + err.Error())
        }
    })
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

    db, err := gorm.Open(driver.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true,
        },
        Logger: logger.New(
            log.New(os.Stdout, "\r\n", log.LstdFlags),
            logger.Config{
                SlowThreshold:             200 * time.Millisecond,
                LogLevel:                  logger.Info,
                IgnoreRecordNotFoundError: false,
                Colorful:                  true,
            }),
    })
    if err != nil {
        return nil, err
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    sqlDB.SetMaxIdleConns(config.Config.MySQL.MaxIdleConns)
    sqlDB.SetMaxOpenConns(config.Config.MySQL.MaxOpenConns)

    return db, nil
}

func NewMySQLRepository() (*MySQLRepository, error) {
    db, err := NewGormDB()
    if err != nil {
        return nil, err
    }
    MySQL = &MySQLRepository{
        db: db,
    }

    return MySQL, nil
}
