package service

import (
	"LPan/dao"
	"LPan/model"
	"database/sql"
)

func NewFile(FileName string, UserId int, FatherPath string) error {
	return dao.AddFile(FileName, UserId, FatherPath)
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
