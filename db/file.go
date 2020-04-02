package db

import (
	"fmt"

	mydb "github.com/rh01/oss-storage/db/mysql"
)

// OnFileUploadFinished 文件上传完成
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file (`file_sha1`, `file_name`,`file_size`,`file_addr`,`status`) values(?,?,?,?,1)",
	)
	if err != nil {
		fmt.Println("Failed to prepare statement, err: ", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Printf("File with hash: %s has been uploaded before", filehash)
			return false
		}
	}
	return true
}
