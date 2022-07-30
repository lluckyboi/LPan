package model

import "time"

type Private struct {
	UserId   int       `json:"user_id"`
	FileName string    `json:"file_name"`
	FileId   int       `json:"file_id"`
	Deleted  time.Time `json:"-"`
}

type Public struct {
	FileName string
	FileId   int
}
