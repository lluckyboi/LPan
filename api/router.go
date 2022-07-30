package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RUNENGINE() {
	r := gin.Default()

	//cors解决跨域问题
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true                                                                                                 //允许所有域名
	config.AllowMethods = []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"}                                                     //允许请求的方法
	config.AllowHeaders = []string{"tus-resumable", "upload-length", "upload-metadata", "cache-control", "x-requested-with", "*"} //允许的Header
	r.Use(cors.New(config))

	//发送邮箱验证码
	r.POST("/getmailac", getmailac)

	//登录注册
	r.POST("/login", login)

	//修改信息
	r.POST("/updateinfo", updateinfo)

	r.Run(":9925")
}
