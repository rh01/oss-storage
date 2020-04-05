package main

import (
	"log"
	"net/http"

	"github.com/rh01/oss-storage/handler"
)

func main() {
	// 静态资源处理
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(assets.AssetFS())))
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// 文件存取接口
	http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/upload/success", handler.HTTPInterceptor(handler.UploadHandlerSuccess))
	http.HandleFunc("/file/meta", handler.HTTPInterceptor(handler.GetFileMetaHandler))
	http.HandleFunc("/file/query", handler.HTTPInterceptor(handler.FileQueryHandler))
	http.HandleFunc("/file/download", handler.HTTPInterceptor(handler.DownloadHandler))
	http.HandleFunc("/file/update", handler.HTTPInterceptor(handler.FileMetaUpdataHandler))
	http.HandleFunc("/file/remove", handler.HTTPInterceptor(handler.FileDeleteHandler))
	// 秒传接口
	http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(handler.TryFastUploadHandler))

	// 用户登录和注册控制器
	http.HandleFunc("/user/signup", handler.SignUpHandler)
	http.HandleFunc("/user/signin", handler.SiginHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	// 首页展示
	http.HandleFunc("/home", handler.HomeHandler)

	// 分块上传接口
	// 初始化分块信息
	http.HandleFunc("/file/mpupload/init", handler.HTTPInterceptor(handler.InitiateMultipartUploadHandler))
	// 上传分块
	http.HandleFunc("/file/mpupload/uppart", handler.HTTPInterceptor(handler.UploadPartHandler))
	// 通知分块上传完成
	http.HandleFunc("/file/mpupload/complete", handler.HTTPInterceptor(handler.CompleteUploadPartHandler))
	// 取消上传分块
	// http.HandleFunc("/file/mpupload/init", handler.HTTPInterceptor(handler.CancelUploadHandler))
	// 查看分块上传的整体状态
	// http.HandleFunc("/file/mpupload/init", handler.HTTPInterceptor(handler.MultipartUploadStatusHandler))

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
