// Package db 主要实现了数据库层的逻辑 -- user op
package db

import (
	"database/sql"
	"fmt"

	mydb "github.com/rh01/oss-storage/db/mysql"
)

// User 用户信息，返回的数据结构体
type User struct {
	Username     string
	Email        string
	Phone        string
	Role         string //权限，什么样的role就可以看到什么样子的离职申请
	SignupAt     string
	LastActiveAt string
	Status       int
}

// GetUserInfo : 获取用户信息
func GetUserInfo(username string) (User, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select user_name, signup_at, role from tbl_user where user_name=? limit 1",
	)

	u := User{}
	if err != nil {
		fmt.Printf("User [%s] not found, err: %v\n", username, err)
		return u, err
	}
	defer stmt.Close()

	// 执行查询操作哦
	err = stmt.QueryRow(username).Scan(&u.Username, &u.SignupAt, &u.Role)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return u, err
	}

	return u, nil
}

// UserSignup : 通过用户名及密码完成user表的注册操作
func UserSignup(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	// 关闭资源
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	// 判断是否重复注册
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

// UserSignup : 通过用户名及密码完成user表的注册操作
func UserSignout(username string) bool {
	sql := fmt.Sprintf(
		"delete from tbl_user_token where username = '%s'", username)

	_, err := mydb.DBConn().Exec(sql)
	if err != nil {
		fmt.Println("Failed to delete token, err:" + err.Error())
		return false
	}
	return true
}

type userInfo struct {
	userName string
	passwd   sql.NullString
}

// UserSignin : 检查用户名是否存在并且密码是否正确
func UserSignin(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"select user_name, user_pwd from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Printf("Failed to get username [%s], err: %v\n", username, err.Error())
		return false
	}
	defer stmt.Close()

	u := userInfo{}
	err = stmt.QueryRow(username).Scan(&u.userName, &u.passwd)
	if err != nil {
		fmt.Printf("Failed to query row username [%s], err: %v\n", username, err.Error())
		return false
	}

	if u.passwd.String != encpwd {
		fmt.Printf("password not equal encpwd  [%s]\n", encpwd)
		return false
	}
	return true
}
