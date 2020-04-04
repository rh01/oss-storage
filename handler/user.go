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
	passWord := r.Form.Get("password")
	encpasswd := utils.Sha1([]byte(passWord + pwdSalt))

	// 校验用户名和密码
	pwdChecked := userdb.UserSignin(username, encpasswd)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}
	// 生成访问凭证 -- token
	token := genToken(username)
	updateResult := userdb.UpdateToken(username, token)
	if !updateResult {
		w.Write([]byte("FAILED"))
		return
	}
	// 登陆成功之后重定向首页--上传页面
	// w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := utils.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

// UserInfoHandler  获取用户信息的接口
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	r.ParseForm()

	username := r.Form.Get("username")
	// TODO：使用拦截器处理这一部分的逻辑
	// token := r.Form.Get("token")
	// 验证token是否有效
	// isValidToken := isTokenVaild(token)
	// if !isValidToken {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	// 查询用户信息
	user, err := userdb.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 组装并且响应用户数据
	resp := &utils.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

func genToken(username string) string {
	// 32 bits + 8 bits
	// 40bits: md5(username + timestamp + token_salt) +timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := utils.MD5([]byte(username + ts + tokenSalt))
	return tokenPrefix + ts[:8]
}

// isTokenVaild 判断当前的token是否有效
func isTokenVaild(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}

// HomeHandler 首页控制器
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/home.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// TODO: 比较 io.WriteString() 和 w.Write 以及 fmt.Fprint 性能问题
		w.Write(data)
		return
	}
}
