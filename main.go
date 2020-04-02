package main

import (
	"log"
	"net/http"

	"github.com/rh01/oss-storage/handler"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/success", handler.UploadHandlerSuccess)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdataHandler)
	http.HandleFunc("/file/remove", handler.FileDeleteHandler)

	http.HandleFunc("/user/signup", handler.SignUpHandler)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
