package service

import (
	"LPan/dao"
	"LPan/model"
	"database/sql"
)

func IsUserExistByMail(UserMail string) (bool, error) {
	_, err := dao.SelectUserByUserMail(UserMail)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	} else if err == sql.ErrNoRows {
		return true, nil
	}
	return false, nil
}

func NewUser(User model.User) error {
	return dao.NewUser(User)
}

func GetUserInfoByMail(UserMail string) (model.User, error) {
	return dao.SelectUserByUserMail(UserMail)
}

func UpdateUserName(UserName string) error {
	return dao.UpadteUserName(UserName)
}
