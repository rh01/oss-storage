package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:heng130509@tcp(127.0.0.1:3306)/oss?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Filed to connect mysql, err: ", err.Error())
		os.Exit(1)
	}
}

// DBConn 返回数据亏连接对象
func DBConn() *sql.DB {
	return db
}
