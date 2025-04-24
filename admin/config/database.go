package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	var err error
	// Thay đổi thông tin kết nối theo cấu hình của bạn
	dsn := "root:01680865@tcp(127.0.0.1:3306)/go_web?charset=utf8mb4&parseTime=true&loc=Local"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	fmt.Println("Kết nối database thành công!")
	return nil
}
