package service

import (
	"LPan/dao"
	"LPan/model"
	"database/sql"
)

func NewFile(FileName string, UserId int, FatherPath, hash string, size int64) error {
	return dao.AddFile(FileName, UserId, FatherPath, hash, size)
}

func CheckAuthorityToDownload(FileID, UserID int) (bool, error, model.Private) {
	private, err := dao.SelectPrivateByUserIdAndFileId(FileID, UserID)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err, private
		} else {
			return false, nil, private
		}
	}
	return true, nil, private
}

func FindTrueNameInPubilcByFileId(FileId int) (string, error) {
	return dao.SelectFileNameByFileId(FileId)
}

// CheckHash 检查要上传的文件是否在公共存储中心存在
func CheckHash(hash string) (bool, error, int) {
	return dao.CheckHash(hash)
}

// AddHashedFile 添加已经存在于公共存储中心的文件
func AddHashedFile(FileName string, UserId int, FatherPath string, FileId int) error {
	return dao.AddPrivateFile(FileName, UserId, FatherPath, FileId)
}

func DeleteFileByUserIdAndFileId(FileId, UserId int) error {
	return dao.DeletePrivateByUserIdAndFileId(UserId, FileId)
}

func RecoverPrivateByUserIdAndFileId(FileId, UserId int) error {
	return dao.RecoverPrivateFileByUserIdAndFleId(UserId, FileId)
}
