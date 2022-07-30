package api

import (
	"LPan/service"
	"LPan/tool"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"

	//"github.com/robfig/cron/v3"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//令牌桶速率和容量 500KB
const (
	rate     = 500 << 10
	capacity = 500 << 10
)

type LimitReaders struct {
	io.ReadSeeker
	r io.Reader
}

//上传文件
func uploadfile(c *gin.Context) {
	//获取用户Id
	UserId := c.MustGet("UserId").(int)
	//目标虚拟路径
	DesPath := c.PostForm("DesPath")

	//获取父级虚拟路径
	Plen := len(DesPath)
	num := 0
	var pd int
	for i := Plen - 1; i >= 0; i-- {
		if DesPath[i] == '/' {
			num++
		}
		if num == 1 {
			pd = i
		} else if num == 2 {
			break
		}
	}
	FatherPath := DesPath[:pd]

	formFile, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("receive formfile error : %v", err)
		tool.RespErrorWithData(c, "文件为空")
		return
	}
	defer formFile.Close()

	//获取文件后缀
	arr := strings.Split(header.Filename, ".")
	extent := arr[len(arr)-1]

	FilePath := "/gopro/src/lpan/file/"
	FileName := strconv.FormatInt(time.Now().Unix(), 10) + "." + extent
	file, err := os.Create(FilePath + FileName)
	if err != nil {
		log.Printf("create file error : %v", err)
		tool.RespInternalError(c)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, formFile)
	if err != nil && err != io.EOF {
		log.Printf("copy file error : %v", err)
		tool.RespInternalError(c)
	}
	log.Printf("%v upload file success", UserId)

	err = service.NewFile(FileName, UserId, FatherPath)
	if err != nil {
		log.Println(err)
		tool.RespInternalError(c)
		return
	}

	tool.RespSuccessful(c, "上传文件")
}

//通过id 下载文件
func downloadfile(c *gin.Context) {
	FileId := tool.StringTOInt(c.Param("file_id"))
	UserID := c.MustGet("UserId").(int)
	//校验下载权限
	isok, err, _ := service.CheckAuthorityToDownload(FileId, UserID)
	if err != nil {
		log.Println("CheckAuthorityToDownload err ", err)
		tool.RespInternalError(c)
		return
	}
	if !isok {
		tool.RespErrorWithData(c, "没有下载权限")
		return
	}
	//在公共存储中心找到真名
	FileName, err := service.FindTrueNameInPubilcByFileId(FileId)
	if err != nil {
		log.Println("FindTrueNameInPubilcByFileId err ", err)
		tool.RespInternalError(c)
		return
	}

	//检查是否VIP决定是否限速
	user, err := service.GetUserInfoByUserId(UserID)
	if err != nil {
		log.Println(err)
		tool.RespInternalError(c)
		return
	}
	//VIP不限速
	if user.Vip == 1 {
		c.File("./file/" + FileName)
	} else {
		//非会员 令牌算法限流安排上
		req := c.Request
		w := c.Writer
		file, err := os.Open("./file/" + FileName)
		if err != nil {
			log.Println(err)
			tool.RespInternalError(c)
			return
		}
		defer file.Close()
		fi, err := file.Stat()
		if err != nil {
			log.Println(err)
			tool.RespInternalError(c)
			return
		}
		bucket := ratelimit.NewBucketWithRate(rate, capacity)
		lr := &LimitReaders{
			ReadSeeker: file,
			r:          ratelimit.Reader(file, bucket),
		}
		http.ServeContent(w, req, "./file/"+FileName, fi.ModTime(), lr)
	}
}

func deletefile() {

}
