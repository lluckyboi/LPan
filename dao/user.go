package dao

import (
	"LPan/model"
)

func NewUser(User model.User) error {
	_, err := Db.Exec("insert into user(user_mail)values(?)", User.UserMail)
	return err
}

func SelectUserByUserMail(UserMail string) (model.User, error) {
	User := model.User{}
	err := Db.QueryRow("select user_id,user_name from user where user_mail=?", UserMail).Scan(&User)
	return User, err
}

func UpadteUserName(UserName string) error {
	_, err := Db.Exec("update user set user_name=?", UserName)
	return err
}
