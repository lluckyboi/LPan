package main

import (
	"LPan/api"
	"LPan/dao"
)

func main() {
	//先连接数据库
	dao.RUNDB()
	//启动引擎
	api.RUNENGINE()
	//每日清理过期文件
	dao.CleanDeletedFile()
}
