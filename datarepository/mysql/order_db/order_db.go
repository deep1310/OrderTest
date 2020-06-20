package order_db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

var (
	sqlDB *gorm.DB
)

func Close(sqlDB *gorm.DB) {
	err := sqlDB.Close()
	if err != nil {
		panic(fmt.Errorf("Fatal error while close DB: %s", err.Error()))
	}
}

func GetSqlConn() *gorm.DB {
	return sqlDB
}

func init() {

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		"root", "1234", "127.0.0.1:3306", "order_db")
	db, err := gorm.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	db.DB().SetConnMaxLifetime(5 * time.Minute)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	if db.Error != nil {
		panic(err)
	}
	sqlDB = db
}

