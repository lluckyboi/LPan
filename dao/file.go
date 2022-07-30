package dao

import (
	"LPan/model"
	"database/sql"
	"log"
	"time"
)

func AddFile(FileName string, UserId int, FatherPath, hash string, size int64) error {
	FileId := 0
	//开启事务
	tx, err := Db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("insert into public_file(file_name,hash,size)values(?,?,?)", FileName, hash, size)
	if err != nil {
		errr := tx.Rollback()
		if errr != nil {
			log.Println("rollback err", errr)
		}
		return err
	}
	err = tx.QueryRow("select file_id from public_file where file_name=?", FileName).Scan(&FileId)
	if err != nil {
		errr := tx.Rollback()
		if errr != nil {
			log.Println("rollback err", errr)
		}
		return err
	}
	_, err = tx.Exec("insert into private(user_id, file_name,file_id,father_path)values(?,?,?,?)", UserId, FileName, FileId, FatherPath)
	//提交事务
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func AddPrivateFile(FileName string, UserId int, FatherPath string, FileId int) error {
	_, err := Db.Exec("insert into private (user_id, file_name,file_id,father_path) values (?,?,?,?)", UserId, FileName, FileId, FatherPath)
	return err
}

func SelectPrivateByUserIdAndFileId(FileID, UserID int) (model.Private, error) {
	Private := model.Private{}
	err := Db.QueryRow("select user_id,file_name,file_id from private where file_id=? and user_id=? and deleted is null", FileID, UserID).Scan(&Private.UserId, &Private.FileName, &Private.FileId)
	return Private, err
}

func SelectFileNameByFileId(FileId int) (string, error) {
	FileName := ""
	err := Db.QueryRow("select file_name from public_file where file_id=?", FileId).Scan(&FileName)
	return FileName, err
}

func CheckHash(hash string) (bool, error, int) {
	fileId := 0
	err := Db.QueryRow("select file_id from public_file where hash=?", hash).Scan(&fileId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil, 0
		}
		return false, err, 0
	}
	return true, nil, fileId
}

func DeletePrivateByUserIdAndFileId(UserId, FileId int) error {
	_, err := Db.Exec("update private set deleted=? where user_id=? and file_id=?", time.Now(), UserId, FileId)
	return err
}

func RecoverPrivateFileByUserIdAndFleId(UserId, FileId int) error {
	_, err := Db.Exec("update private set deleted=null where user_id=? and file_id=?", UserId, FileId)
	return err
}

func RenamePrivateFileByUserIdAndFileId(UserId, FileId int, NewName string) error {
	_, err := Db.Exec("update private set file_name = ? where user_id=? and file_id=?", NewName, UserId, FileId)
	return err
}

func UpdatePathByUserIdAndFileId(UserId, FileId int, NewPath string) error {
	_, err := Db.Exec("update private set father_path = ? where user_id=? and file_id=?", NewPath, UserId, FileId)
	return err
}

func SetShareByUserIdAndFileId(UserId, FileId int, Expr time.Time) error {
	_, err := Db.Exec("update private set share=1 ,expr_time=? where user_id=? and file_id=?", Expr, UserId, FileId)
	return err
}
