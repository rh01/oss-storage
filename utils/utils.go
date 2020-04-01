package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"os"
	"path/filepath"
)

var (
	FileNotFound = errors.New("File Not Found") // FileNotFound xxxx
)

// Sha1Stream is sha1 enc
type Sha1Stream struct {
	_sha1 hash.Hash
}

// Update 更新hash
func (o *Sha1Stream) Update(data []byte) {
	if o._sha1 == nil {
		o._sha1 = sha1.New()
	}
	o._sha1.Write(data)
}

// Sum xxxx
func (o *Sha1Stream) Sum() string {
	return hex.EncodeToString(o._sha1.Sum([]byte("")))
}

// Sha1 编码
func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

// FileSha1 文件hash, 内容哈希
func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

// MD5 编码
func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

// FileMD5 文件MD5
func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

// PathExists 检查文件/路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// GetFileSize 获取文件大小
func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}
