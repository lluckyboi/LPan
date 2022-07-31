package api

import (
	"LPan/service"
	"LPan/tool"
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

	fileGroup := r.Group("file", JWTAuthMiddleware(), RateLimitBuck(time.Millisecond*100, 2048))
	{
		fileGroup.POST("/upload", uploadfile)
		fileGroup.POST("/uploadbysilce", uploadfilebysile)
		fileGroup.GET("/download/:file_id", downloadfile)
		fileGroup.GET("/downloadbysilce/:file_id", downloadbyslice)
		fileGroup.DELETE("/delete/:file_id", deletefile)
		fileGroup.GET("/recover/:file_id", recoverfile)
		fileGroup.POST("/rename/:file_id", renamefile)
		fileGroup.POST("/modify-path/:file_id", modifypath)
		fileGroup.POST("/share/:file_id", sharefile)
		fileGroup.POST("/share/secret/:file_id", sharefilewithsecret)
	}

	//解密链接重定向
	r.GET("/secret/:val", JWTAuthMiddleware(), RateLimitBuck(time.Millisecond*100, 2048), func(c *gin.Context) {
		val := c.Param("val")
		path, err := service.GetOriginBySec(val)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{
					"info": "no such link",
				})
				return
			}
			log.Println(err)
			tool.RespInternalError(c)
			return
		}
		c.Redirect(http.StatusMovedPermanently, "http://39.106.81.229:9925/"+path)
	})

	r.Run(":9925")
}
