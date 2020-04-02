package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rh01/oss-storage/meta"
	"github.com/rh01/oss-storage/utils"
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
		// w.Write(bytes)
		io.WriteString(w, string(bytes))
		w.WriteHeader(http.StatusOK)
	} else if r.Method == "POST" {
		// 接受文件流存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Fprintf(w, "Failed to get data, err: %s\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: uploadPath + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 12:02:22"),
		}
		newFile, err := os.Create(uploadPath + fileMeta.FileName)
		if err != nil {
			fmt.Fprintf(w, "Failed to create file, err: %v\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Fprintf(w, "Failed to write file, err: %v\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = utils.FileSha1(newFile)
		// meta.UpdateFileMeta(fileMeta)
		_ = meta.UploadFileMetaDB(&fileMeta)

		var buff = &bytes.Buffer{}
		err = json.NewEncoder(buff).Encode(fileMeta)
		if err != nil {
			fmt.Fprintf(w, "Failed to encode file, err: %v\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, buff.String())

		// http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

// UploadHandlerSuccess 上传成功
func UploadHandlerSuccess(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "success upload")
	w.WriteHeader(http.StatusOK)
}

// GetFileMetaHandler 查询文件的元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	filehash := r.Form.Get("filehash")
	// fMeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Fprintf(w, "Failed to encode fileMeta, err: %v\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var buff bytes.Buffer
	err := json.NewEncoder(&buff).Encode(&fMeta)
	if err != nil {
		fmt.Fprintf(w, "Failed to encode fileMeta, err: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, buff.String())
	w.WriteHeader(http.StatusOK)
}

//DownloadHandler 下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fSha1)

	loc := fMeta.Location
	file, err := os.Open(loc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Description", "attachment;filename=\""+fMeta.FileName+"\"")
	w.Write(data)
}

// FileMetaUpdataHandler 更新元数据信息
func FileMetaUpdataHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName

	err := os.Rename(curFileMeta.Location, uploadPath+curFileMeta.FileName)
	if err != nil {
		w.Write([]byte("rename failed."))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	curFileMeta.Location = uploadPath + curFileMeta.FileName

	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(&curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	w.WriteHeader(200)
}

// FileDeleteHandler 删除文件操作
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fSha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fSha1)
	w.WriteHeader(http.StatusOK)
}
