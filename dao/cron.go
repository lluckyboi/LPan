package dao

import (
	"github.com/robfig/cron/v3"
	"time"
)

func CleanDeletedFile() {
	c := cron.New()
	c.AddFunc("@daily", func() {
		Db.Exec("delete from private where deleted<?", time.Now().Add(time.Hour*24*30))
	})
	c.Start()
	select {}
}
