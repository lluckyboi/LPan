package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func RUNENGINE() {
	r := gin.Default()
	//限制上传为文件最大512MB
	r.MaxMultipartMemory = 1 << 29

	//cors解决跨域问题
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true                                                                                                 //允许所有域名
	config.AllowMethods = []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"}                                                     //允许请求的方法
	config.AllowHeaders = []string{"tus-resumable", "upload-length", "upload-metadata", "cache-control", "x-requested-with", "*"} //允许的Header
	r.Use(cors.New(config))

	//发送邮箱验证码
	r.POST("/get-mail-ac", getmailac)

	//登录注册
	r.POST("/login", login)

	//修改信息
	r.POST("/update-info", updateinfo)

	fileGroup := r.Group("file", JWTAuthMiddleware(), RateLimitMiddleware(time.Millisecond*100, 2048))
	{
		fileGroup.POST("/upload", uploadfile)
		fileGroup.GET("/download/:file_id", downloadfile)
		fileGroup.DELETE("/delete/:file_id", deletefile)
		fileGroup.GET("/recover/:file_id", recoverfile)
		fileGroup.POST("/rename/:file_id", renamefile)
		fileGroup.POST("/modify-path/:file_id", modifypath)
	}

	r.Run(":9925")
}
