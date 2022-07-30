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
	err := Db.QueryRow("select user_id,user_name from user where user_mail=?", UserMail).
		Scan(&User.UserId, &User.UserName)
	return User, err
}

func UpadteUserName(UserName string) error {
	_, err := Db.Exec("update user set user_name=?", UserName)
	return err
}

func SelectUserByUserId(UserId int) (model.User, error) {
	User := model.User{}
	err := Db.QueryRow("select user_id,user_name,user_mail,vip from user where user_id=?", UserId).
		Scan(&User.UserId, &User.UserName, &User.UserMail, &User.Vip)
	return User, err
}
