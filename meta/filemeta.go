package meta

import "github.com/rh01/baiduyun/utils"

// FileMeta 文件元数据信息结构体
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta 更新文件的元数据信息
func UpdateFileMeta(fmeta FileMeta) {
	fSha1 := fmeta.FileSha1
	if _, ok := fileMetas[fSha1]; ok {
		delete(fileMetas, fSha1)
	}
	fileMetas[fSha1] = fmeta
}

// GetFileMeta 根据哈希获取文件的元数据信息
func GetFileMeta(fSha1 string) FileMeta {
	return fileMetas[fSha1]
}

// RemoveFileMeta 删除文件的元数据信息
func RemoveFileMeta(fSha1 string) error {
	if _, ok := fileMetas[fSha1]; !ok {
		return utils.FileNotFound
	}
	delete(fileMetas, fSha1)
	return nil
}
