package dao

import (
	"LPan/model"
	"log"
)

func AddFile(FileName string, UserId int) error {
	FileId := 0
	//开启事务
	tx, err := Db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("insert into public_file(file_name)values(?)", FileName)
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
	_, err = tx.Exec("insert into private(user_id, file_name,file_id)values(?,?,?)", UserId, FileName, FileId)
	//提交事务
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func SelectPrivateByUserIdAndFileId(FileID, UserID int) (model.Private, error) {
	Private := model.Private{}
	err := Db.QueryRow("select *from private where file_id=?and user_id=? and deleted is null", FileID, UserID).Scan(&Private.UserId, &Private.FileName, &Private.FileId, &Private.Deleted)
	return Private, err
}
