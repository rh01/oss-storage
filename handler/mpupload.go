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
	chunkPrefix      = "chkidx_"
)

// MultipartUploadInfo : 分块上传的信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
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
		w.Write(utils.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}

	//获取redis的一个连接
	rConn := cache.RedisPool().Get()
	defer rConn.Close()

	//生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  chunkSize,
		ChunkCount: int(math.Ceil(float64(filesize) / chunkSize)),
	}

	//将初始化信息写入到redis缓存中
	rConn.Do("HSET", mpPrefix+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", mpPrefix+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", mpPrefix+upInfo.UploadID, "filesize", upInfo.FileSize)
	//将响应初始化数据返回客户端
	w.Write(utils.NewRespMsg(0, "OK", upInfo).JSONBytes())
}

// UploadPartHandler : 上传文件分块接口
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//1 解析用户请求参数
	r.ParseForm()

	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	//2 获得redis连接池中的一个连接
	rConn := cache.RedisPool().Get()
	defer rConn.Close()

	//3 获取文件句柄，并用于存储分块内容
	fpath := mpLoadPathPrefix + uploadID + "/" + chunkIndex
	// 如果该文件前缀的文件不存在，则创建文件夹
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(utils.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
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
	rConn.Do("HSET", mpPrefix+uploadID, chunkPrefix+chunkIndex, 1)
	//5 返回结果给客户端
	w.Write(utils.NewRespMsg(0, "OK", nil).JSONBytes())
}

// CompleteUploadPartHandler : 通知上传合并
func CompleteUploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//1 解析请求参数
	r.ParseForm()

	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	//2 获取redis连接吃的一个连接
	rConn := cache.RedisPool().Get()
	defer rConn.Close()

	//3 通过uploadid查询redis并判断是否所有分块上传完成
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+upid))
	if err != nil {
		w.Write(utils.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
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
		} else if strings.HasPrefix(k, chunkPrefix) && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		w.Write(utils.NewRespMsg(-2, "invaild request", nil).JSONBytes())
		return
	}
	//4 TODO 合并分块

	//5 更新唯一文件表以及用户文件表
	fsize, _ := strconv.Atoi(filesize)
	db.OnFileUploadFinished(filehash, filename, int64(fsize), "")
	db.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))

	//6 响处理结果
	w.Write(utils.NewRespMsg(0, "OK", nil).JSONBytes())
}
