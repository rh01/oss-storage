package meta

import (
	mydb "github.com/rh01/oss-storage/db"
	"github.com/rh01/oss-storage/utils"
)

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

// UploadFileMetaDB 上传文件元数据信息刀片数据库
func UploadFileMetaDB(fMeta *FileMeta) bool {
	return mydb.OnFileUploadFinished(fMeta.FileSha1, fMeta.FileName, fMeta.FileSize, fMeta.Location)
}

// GetFileMetaDB 获取元数据信息从数据库
func GetFileMetaDB(fSha1 string) (*FileMeta, error) {
	tFile, err := mydb.GetFileMeta(fSha1)
	if err != nil {
		return nil, err
	}
	fMeta := FileMeta{
		FileSha1: tFile.FileHash,
		FileName: tFile.FileName.String,
		FileSize: tFile.FileSize.Int64,
		Location: tFile.FileAddr.String,
	}
	return &fMeta, nil
}
