// Package handler 主要实现controller功能
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/rh01/oss-storage/db"
)

const (
	uploadPath = "/tmp/"
)

// UploadHandler 上传文件
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		bytes, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			fmt.Fprint(w, "static file not found.")
			w.WriteHeader(http.StatusNotFound)
		}
		io.WriteString(w, string(bytes))
		w.WriteHeader(http.StatusOK)
	} else if r.Method == "POST" {
		// TODO 更新用户文件表记录
		username := r.Form.Get("username")
		rdEmail := r.Form.Get("rd_email")
		fmt.Println("username", username, "rd_email", rdEmail)
		suc := db.OnUserApplyFinished(username, rdEmail)
		fmt.Println("suc", suc)
		if suc {
			http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
		} else {
			w.Write([]byte("Apply Failed"))
		}
	}
}

// ManageHandler : 上传文件
func ManageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		bytes, err := ioutil.ReadFile("./static/view/manage.html")
		if err != nil {
			fmt.Fprint(w, "static file not found.")
			w.WriteHeader(http.StatusNotFound)
		}
		io.WriteString(w, string(bytes))
		w.WriteHeader(http.StatusOK)
	} else if r.Method == "POST" {
		// TODO 更新用户文件表记录
		username := r.Form.Get("username")
		rdEmail := r.Form.Get("rd_email")
		fmt.Println("username", username, "rd_email", rdEmail)
		suc := db.OnUserApplyFinished(username, rdEmail)
		fmt.Println("suc", suc)
		if suc {
			http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
		} else {
			w.Write([]byte("Apply Failed"))
		}
	}
}

// AgreeUserApplyInfo 同意用户的申请信息
func AgreeUserApply(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("user_name")
	approvalUser := r.Form.Get("username")
	rdEmail := r.Form.Get("rd_email")
	u, _ := db.GetUserInfo(approvalUser)
	approvalUserRole := u.Role

	fmt.Println("user_name", username, "approvalUserRole", approvalUserRole)
	suc := db.AgreeUserApplyInfo(rdEmail, approvalUserRole)
	fmt.Println("suc", suc)
	if suc {
		http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
	} else {
		w.Write([]byte("Agree Failed"))
	}
	// fmt.Fprint(w, "success upload")
	// w.WriteHeader(http.StatusOK)
}

// UploadHandlerSuccess 上传成功
func UploadHandlerSuccess(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "success upload")
	w.WriteHeader(http.StatusOK)
}

// ApplyQueryHandler : 批量查询的申请信息
func ApplyQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")

	u, err := db.GetUserInfo(username)
	if err != nil {
		return
	}
	role := u.Role
	fmt.Println("role", role)
	uFiles, err := db.QueryUserApplyInfoByRole(role, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(uFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// ApplyQueryHandlerByUserName : 查询自己的申请信息
func ApplyQueryHandlerByUserName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	uFiles, err := db.QueryUserApplyInfoByUserName(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(uFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// ModifyHandler : 更改自己的申请信息
func ModifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		bytes, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			fmt.Fprint(w, "static file not found.")
			w.WriteHeader(http.StatusNotFound)
		}
		io.WriteString(w, string(bytes))
		w.WriteHeader(http.StatusOK)
	} else if r.Method == "POST" {
		r.ParseForm()
		// TODO 更新用户文件表记录
		username := r.Form.Get("username")
		apply := r.Form.Get("apply")
		rdEmail := r.Form.Get("rd_email")
		fmt.Println("username", username, "apply", apply, "rd_email", rdEmail)
		suc := db.UpdateUserApplyInfo(username,
			apply, rdEmail)
		fmt.Println("suc", suc)
		if suc {
			http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
		} else {
			w.Write([]byte("Apply Failed"))
		}
	}
}

// DeleteHandler : 删除自己的申请信息
func DeleteHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	// TODO 更新用户文件表记录
	username := r.Form.Get("username")
	fmt.Println("username", username)
	suc := db.DeleteUserApplyInfo(username)
	fmt.Println("suc", suc)
	if suc {
		http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
	} else {
		w.Write([]byte("Delete Failed"))
	}

}
