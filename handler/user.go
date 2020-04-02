// Package handler 处理用户的逻辑
package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	userdb "github.com/rh01/oss-storage/db"
	"github.com/rh01/oss-storage/utils"
)

const (
	// 用于加密的盐值(自定义)
	pwdSalt   = "^%#890"
	tokenSalt = "_tokensalt"
)

// SignUpHandler 处理用户注册请求
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// TODO: 比较 io.WriteString() 和 w.Write 以及 fmt.Fprint 性能问题
		w.Write(data)
		return
	}

	r.ParseForm()
	userName := r.Form.Get("username")
	passWord := r.Form.Get("password")
	// pre vailation
	if len(userName) < 3 || len(passWord) < 5 {
		w.Write([]byte("Invaild Parameter"))
		return
	}
	// encroption password
	encPassword := utils.Sha1([]byte(passWord + pwdSalt))
	// TODO: 需要处理
	succ := userdb.UserSignup(userName, encPassword)
	if succ {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

// SiginHandler : 用户登陆接口
func SiginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	encpasswd := r.Form.Get("password")
	// 校验用户名和密码
	pwdChecked := userdb.UserSignin(username, encpasswd)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}
	// 生成访问凭证 -- token
	_ = genToken(username)

	// 登陆成功之后重定向首页--上传页面
}

func genToken(username string) string {
	// 32 bits + 8 bits
	// 40bits: md5(username + timestamp + token_salt) +timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := utils.MD5([]byte(username + ts + tokenSalt))
	return tokenPrefix + ts[:8]
}
