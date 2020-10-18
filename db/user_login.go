// Package db 主要实现了数据库层的逻辑 -- user token op
package db

import (
	"fmt"

	mydb "github.com/rh01/oss-storage/db/mysql"
)

// UpdateToken : 更新用户token，用来登陆使用
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`,`user_token`) values (?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
