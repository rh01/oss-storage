// Package db 主要实现了数据库层的逻辑 --- file op
package db

import (
	"fmt"
	"log"
	"strconv"
	"time"

	mydb "github.com/rh01/oss-storage/db/mysql"
)

// UserApply : 用户离职申请表
type UserApply struct {
	UserName  string
	Apply     string //离职申请
	RdEmail   string //rd邮箱
	ApplyAt   string //申请的时间
	FinishTag string //用来标记进行到哪一步骤
}

// OnUserApplyFinished : 插入一条离职申请
// 當用戶發起申請後的數據插入邏輯
func OnUserApplyFinished(username, rdEmail string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert  into tbl_user_apply  (`user_name`, `rd_email`,`apply_at`, `finish_tag`) values(?,?,?,?) ")
	if err != nil {
		fmt.Println("insert failed, err: ", err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, rdEmail, time.Now(), 1)
	if err != nil {
		fmt.Println("err: ", err.Error())
		return false
	}
	return true
}

// QueryUserApplyInfoByUserName : 获取/检索离职用户的离职申请，根据用户名查找
func QueryUserApplyInfoByUserName(username string) (*[]UserApply, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select user_name, rd_email,apply_at,finish_tag from tbl_user_apply")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("err: ", err.Error())
	}

	var uApplies []UserApply
	for rows.Next() {
		var uApply UserApply
		err = rows.Scan(&uApply.UserName, &uApply.RdEmail, &uApply.ApplyAt, &uApply.FinishTag)
		if err != nil {
			fmt.Println("err: ", err.Error())
			break
		}
		uApplies = append(uApplies, uApply)
	}

	return &uApplies, nil
}

// UpdateUserApplyInfo : 更新离职用户的离职申请
func UpdateUserApplyInfo(username, apply, rdEmail string) bool {
	sql := fmt.Sprintf("update tbl_user_apply set apply='%s', rd_email='%s' where user_name='%s'", apply, rdEmail, username)
	_, err := mydb.DBConn().Exec(sql)
	if err != nil {
		log.Println("exec failed:", err, ", sql:", sql)
		return false
	}
	return true
}

// DeleteUserApplyInfo : 删除离职用户的离职申请
func DeleteUserApplyInfo(username string) bool {
	sql := fmt.Sprintf("delete from tbl_user_apply  where user_name='%s'", username)
	_, err := mydb.DBConn().Exec(sql)
	if err != nil {
		log.Println("exec failed:", err, ", sql:", sql)
		return false
	}
	return true
}

// AgreeUserApplyInfo : 同意该用户的离职申请
func AgreeUserApplyInfo(rd_email string, approverRole string) bool {
	iRole, _ := strconv.Atoi(approverRole)
	targetRole := iRole + 1
	sql := fmt.Sprintf("update tbl_user_apply set finish_tag='%s' where rd_email='%s'", strconv.Itoa(targetRole), rd_email)

	_, err := mydb.DBConn().Exec(sql)
	if err != nil {
		log.Println("exec failed:", err, ", sql:", sql)
		return false
	}
	return true
}

// QueryUserApplyInfoByRole : 批量获取/检索离职用户的离职申请，根据用户的权限查找
func QueryUserApplyInfoByRole(role string, limit int) ([]UserApply, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select user_name, rd_email, apply_at, finish_tag from tbl_user_apply where finish_tag=? limit ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(1, limit)
	if err != nil {
		return nil, err
	}

	var uApplies []UserApply
	for rows.Next() {
		var uApply UserApply
		err = rows.Scan(&uApply.UserName, &uApply.RdEmail, &uApply.ApplyAt, &uApply.FinishTag)
		if err != nil {
			fmt.Println("err: ", err.Error())
			break
		}
		uApplies = append(uApplies, uApply)
	}

	return uApplies, nil
}
