package model

type User struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
	UserMail string `json:"user_mail"`
	Vip      int    `json:"vip"`
}

type AcUser struct {
	UserId int    `json:"user_id"`
	AcCode string `json:"ac_code"`
}
