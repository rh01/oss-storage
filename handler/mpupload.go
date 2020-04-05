// Package handler 主要实现controller功能
package handler

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	cache "github.com/rh01/oss-storage/cache/redis"
	"github.com/rh01/oss-storage/db"
	"github.com/rh01/oss-storage/utils"
)

const (
	chunkSize        = 5 * 1024 * 1024 // 5MB
	buffSize         = 1024 * 1024     // 1MB
	mpLoadPathPrefix = "/data/"        // 分片上传的前缀
	mpPrefix         = "MP_"
)

// MultipartUploadInfo : 分块上传的信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int64
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

// InitiateMultipartUploadHandler : 初始化分块信息上传
func InitiateMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	//解析用户请求信息
	r.ParseForm()

	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		resp := utils.RespMsg{
			Code: -1,
			Msg:  "Params invaild",
		}
		w.Write(resp.JSONBytes())
		return
	}

	//获取redis的一个连接
	rConn := cache.RedisPoll().Get()
	defer rConn.Close()

	//生成分块上传的初始化信息
	mpInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   int64(filesize),
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  chunkSize,
		ChunkCount: int(math.Ceil(float64(filesize) / chunkSize)),
	}

	//将初始化信息写入到redis缓存中
	rConn.Do("HSET", mpPrefix+mpInfo.UploadID, "chunkcount", mpInfo.ChunkCount)
	rConn.Do("HSET", mpPrefix+mpInfo.UploadID, "filehash", mpInfo.FileHash)
	rConn.Do("HSET", mpPrefix+mpInfo.UploadID, "filesize", mpInfo.FileSize)
	//将响应初始化数据返回客户端
	resp := utils.RespMsg{
		Code: 0,
		Msg:  "OK",
	}
	w.Write(resp.JSONBytes())
}

// UploadPartHandler : 上传文件分块接口
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//1 解析用户请求参数
	r.ParseForm()

	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	//2 获得redis连接池中的一个连接
	rConn := cache.RedisPoll().Get()
	defer rConn.Close()

	//3 获取文件句柄，并用于存储分块内容
	fpath := mpLoadPathPrefix + uploadID + "/" + chunkIndex
	// 如果该文件前缀的文件不存在，则创建该文件夹
	os.Mkdir(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		resp := utils.RespMsg{
			Code: -1,
			Msg:  "Upload part failed",
		}
		w.Write(resp.JSONBytes())
		return
	}
	defer fd.Close()

	buf := make([]byte, buffSize)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	//4 更新redis缓存状态
	rConn.Do("HSET", mpPrefix+uploadID, 1)
	//5 返回结果给客户端
	resp := utils.RespMsg{
		Code: 0,
		Msg:  "OK",
	}
	w.Write(resp.JSONBytes())
}

// CompleteUploadPartHandler : 通知上传合并
func CompleteUploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//1 解析请求参数
	r.ParseForm()

	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		resp := utils.RespMsg{
			Code: -1,
			Msg:  "invailed params.",
		}
		w.Write(resp.JSONBytes())
		return
	}

	filename := r.Form.Get("filename")

	//2 获取redis连接吃的一个连接
	rConn := cache.RedisPoll().Get()
	defer rConn.Close()

	//3 通过uploadid查询redis并判断是否所有分块上传完成
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+upid))
	if err != nil {
		resp := utils.RespMsg{
			Code: -1,
			Msg:  "complete upload failed",
		}
		w.Write(resp.JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	// 由于data返回的结果kv是连续的，占两个位置，因此这里+2
	for i := 0; i < len(data); i += 2 {
		// 类型转换，可以用断言
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		resp := utils.RespMsg{
			Code: -2,
			Msg:  "invaold request",
		}
		w.Write(resp.JSONBytes())
		return
	}
	//4 TODO 合并分块

	//5 更新唯一文件表以及用户文件表
	suc := db.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))
	if !suc {
		resp := utils.RespMsg{
			Code: -3,
			Msg:  "update userfile failed",
		}
		w.Write(resp.JSONBytes())
		return
	}

	//6 响处理结果
	resp := utils.RespMsg{
		Code: 0,
		Msg:  "OK",
	}
	w.Write(resp.JSONBytes())
}
