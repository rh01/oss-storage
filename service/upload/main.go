package main

import (
	"log"
	"net/http"

	"github.com/rh01/oss-storage/handler"
)

func main() {
	// 静态资源处理
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	//离职申请与管理接口
	http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/query", handler.HTTPInterceptor(handler.ApplyQueryHandlerByUserName))
	http.HandleFunc("/file/queryAll", handler.HTTPInterceptor(handler.ApplyQueryHandler))
	http.HandleFunc("/file/modify", handler.HTTPInterceptor(handler.ModifyHandler))
	http.HandleFunc("/file/delete", handler.HTTPInterceptor(handler.DeleteHandler))
	http.HandleFunc("/file/manage", handler.HTTPInterceptor(handler.ManageHandler))
	http.HandleFunc("/file/agree", handler.HTTPInterceptor(handler.AgreeUserApply))

	// 用户登录和注册控制器
	http.HandleFunc("/user/signup", handler.SignUpHandler)
	http.HandleFunc("/user/signin", handler.SiginHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	// 首页展示
	http.HandleFunc("/home", handler.HTTPInterceptor(handler.HomeHandler))

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
