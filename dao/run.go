package dao

import (
	"database/sql"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// Db 全局变量
var Db *sql.DB
var RedisDB *redis.Client

func RUNDB() {
	//启用mysql数据库
	db, err := sql.Open("mysql", "lpan:pzSyrjstpNzyp8wc@/lpan")
	if err != nil {
		log.Fatal(err)
	}
	Db = db

	//初始化redis
	initClient()
}

func initClient() (err error) {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err = RedisDB.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
