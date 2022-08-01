package dao

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func CleanDeletedFile() {
	c := cron.New()
	c.AddFunc("@daily", func() {
		_, err := Db.Exec("delete from private where deleted<=? or (share=1 and expr_time <= ?)", time.Now().Add(time.Hour*24*30), time.Now())
		if err != nil {
			log.Fatal(err)
		}
	})
	c.Start()
	select {}
}
