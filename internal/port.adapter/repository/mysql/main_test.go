package mysql

import (
    "context"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    driver "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

/**
 * @author Rancho
 * @date 2021/12/21
 */

var ctx = context.Background()

func mockMySQL() (*gorm.DB, sqlmock.Sqlmock) {
    sqlDB, mock, err := sqlmock.New()
    if err != nil {
        panic("mock MySQL fail, err: " + err.Error())
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
    MySQL = &MySQLRepository{
        db: db,
    }

    return db, mock
}

func TestMain(m *testing.M) {
    m.Run()
}
