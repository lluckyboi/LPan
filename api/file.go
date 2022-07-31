package api

import (
	"LPan/service"
	"LPan/tool"
	sha12 "crypto/sha1"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"github.com/skip2/go-qrcode"
	"io"
	"log"
	"net/http"
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

	//检查HASH
	hash, err := tool.GetHash(formFile)
	if err != nil {
		log.Printf("GetHash error : %v", err)
		tool.RespInternalError(c)
		return
	}

	bl, err, fileId := service.CheckHash(hash)
	if err != nil {
		log.Printf("check hash error : %v", err)
		tool.RespInternalError(c)
		return
	}
	//如果哈希值存在 秒传
	if bl {
		err := service.AddHashedFile(FileName, UserId, FatherPath, fileId)
		if err != nil {
			log.Printf("add hashed file error : %v", err)
			tool.RespInternalError(c)
			return
		}
		tool.RespSuccessful(c, "秒传")
		return
	}

	//否则继续上传
	_, err = io.Copy(file, formFile)
	if err != nil && err != io.EOF {
		log.Printf("copy file error : %v", err)
		tool.RespInternalError(c)
	}
	log.Printf("user %v upload file success", UserId)

	//size
	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		tool.RespInternalError(c)
		return
	}
	size := info.Size()

	//入库
	err = service.NewFile(FileName, UserId, FatherPath, hash, size)
	if err != nil {
		log.Println(err)
		tool.RespInternalError(c)
		return
	}

	tool.RespSuccessful(c, "上传文件")
}

//断点续传上传文件
func uploadfilebysile(c *gin.Context) {
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

	//检查HASH
	hash, err := tool.GetHash(formFile)
	if err != nil {
		log.Printf("GetHash error : %v", err)
		tool.RespInternalError(c)
		return
	}
	bl, err, fileId := service.CheckHash(hash)
	if err != nil {
		log.Printf("check hash error : %v", err)
		tool.RespInternalError(c)
		return
	}
	//如果哈希值存在 秒传
	if bl {
		err := service.AddHashedFile(FileName, UserId, FatherPath, fileId)
		if err != nil {
			log.Printf("add hashed file error : %v", err)
			tool.RespInternalError(c)
			return
		}
		tool.RespSuccessful(c, "秒传")
		return
	}

	//否则继续断点上传
	//先生成临时文件
	tempFile := FilePath + "temp.txt"
	tempf, _ := os.OpenFile(tempFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer tempf.Close()
	_, err = tempf.Seek(0, io.SeekStart)
	if err != nil {
		log.Printf("Seek error : %v", err)
		tool.RespInternalError(c)
		return
	}

	buffer := make([]byte, 100, 100)
	len1, err := tempf.Read(buffer)
	ctstr := string(buffer[:len1])
	ct, _ := strconv.ParseInt(ctstr, 10, 64)
	_, err = formFile.Seek(ct, 0)
	if err != nil {
		log.Printf("Seek error : %v", err)
		tool.RespInternalError(c)
		return
	}

	_, err = file.Seek(ct, 0)
	if err != nil {
		log.Printf("Seek error : %v", err)
		tool.RespInternalError(c)
		return
	}

	data := make([]byte, 2048, 2048)
	n2 := -1         // 读取的数据量
	n3 := -1         //写出的数据量
	total := int(ct) //读取的总量
	for {
		n2, err = formFile.Read(data)
		if err == io.EOF {
			log.Println("copy finished")
			os.Remove(tempFile)
			break
		}
	}
	//将数据写入到目标文件
	n3, _ = file.Write(data[:n2])
	total += n3
	//将复制总量，存储到临时文件中
	tempf.Seek(0, io.SeekStart)

	_, err = tempf.WriteString(strconv.Itoa(total))
	if err != nil {
		log.Printf("WriteString error : %v", err)
		tool.RespInternalError(c)
		return
	}

	log.Printf("%v upload file success", UserId)

	//size
	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		tool.RespInternalError(c)
		return
	}
	size := info.Size()
	//入库
	err = service.NewFile(FileName, UserId, FatherPath, hash, size)
	if err != nil {
		log.Println(err)
		tool.RespInternalError(c)
		return
	}

	tool.RespSuccessful(c, "断点上传文件")
}

//通过id 下载文件
func downloadfile(c *gin.Context) {
	FileId := tool.StringTOInt(c.Param("file_id"))
	UserID := c.MustGet("UserId").(int)

	//是否为分享的文件
	share := c.Query("share")
	if share == "true" {
		expr := c.Query("expr")
		if expr == "" {
			log.Println("expr get err ")
			tool.RespInternalError(c)
			return
		}
		//校验
		isok, err, _ := service.CheckAuthorityToDownload(FileId, UserID)
		if err != nil {
			log.Println("CheckAuthorityToDownload err ", err)
			tool.RespInternalError(c)
			return
		}
		if !isok {
			inexpr := tool.StringTOInt(expr)
			ExprTime := time.Now().Add(time.Duration(inexpr) * time.Hour * 24)
			err := service.AddHashedFile(strconv.Itoa(FileId), UserID, "/", FileId)
			if err != nil {
				log.Println("AddHashedFile err ", err)
				tool.RespInternalError(c)
				return
			}
			err = service.SetShareByUserIdAndFileId(UserID, FileId, ExprTime)
			if err != nil {
				log.Println("SetShareByUserIdAndFileId err ", err)
				tool.RespInternalError(c)
				return
			}
		}
	}

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
		info, err := file.Stat()
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
		http.ServeContent(w, req, "./file/"+FileName, info.ModTime(), lr)
	}
}

//删除私人仓库中文件
func deletefile(c *gin.Context) {
	FileId := tool.StringTOInt(c.Param("file_id"))
	UserId := c.MustGet("UserId").(int)
	//权限检查
	isok, err, _ := service.CheckAuthorityToDownload(FileId, UserId)
	if err != nil {
		log.Println("CheckAuthorityToDownload err ", err)
		tool.RespInternalError(c)
		return
	}
	if !isok {
		tool.RespErrorWithData(c, "没有删除权限")
		return
	}

	err = service.DeleteFileByUserIdAndFileId(FileId, UserId)
	if err != nil {
		log.Println("CheckAuthorityToDownload err ", err)
		tool.RespInternalError(c)
		return
	}

	tool.RespSuccessful(c, "删除文件")
}

//从回收站找回文件
func recoverfile(c *gin.Context) {
	FileId := tool.StringTOInt(c.Param("file_id"))
	UserId := c.MustGet("UserId").(int)
	//权限检查
	isok, err, _ := service.CheckAuthorityToDownload(FileId, UserId)
	if err != nil {
		log.Println("CheckAuthorityToDownload err ", err)
		tool.RespInternalError(c)
		return
	}
	if !isok {
		tool.RespErrorWithData(c, "没有删除权限")
		return
	}
	err = service.RecoverPrivateByUserIdAndFileId(FileId, UserId)
	if err != nil {
		log.Println("CheckAuthorityToDownload err ", err)
		tool.RespInternalError(c)
		return
	}

	tool.RespSuccessful(c, "找回文件")
}

//文件改名
func renamefile(c *gin.Context) {
	FileId := tool.StringTOInt(c.Param("file_id"))
	UserId := c.MustGet("UserId").(int)
	NewName := c.PostForm("NewName")

	err := service.RenameFileInPrivateByUserIdAndFileId(FileId, UserId, NewName)
	if err != nil {
		log.Println("RenameFileInPrivateByUserIdAndFileId err ", err)
		tool.RespInternalError(c)
		return
	}
	tool.RespSuccessful(c, "改名")
}

//路径修改
func modifypath(c *gin.Context) {
	FileId := tool.StringTOInt(c.Param("file_id"))
	UserId := c.MustGet("UserId").(int)
	NewPath := c.PostForm("NewPath")

	err := service.ModifyPathByUserIdAndFileId(UserId, FileId, NewPath)
	if err != nil {
		log.Println("ModifyPathByUserIdAndFileId err ", err)
		tool.RespInternalError(c)
		return
	}
	tool.RespSuccessful(c, "修改路径 ")
}

//生成分享二维码 不加密 拿到链接的可以下载
func sharefile(c *gin.Context) {
	UserId := c.MustGet("UserId").(int)
	FileId := c.Param("file_id")
	Expr := c.PostForm("expr_time")

	//校验权限
	isok, err, _ := service.CheckAuthorityToDownload(tool.StringTOInt(FileId), UserId)
	if err != nil {
		log.Println("CheckAuthorityToDownload err ", err)
		tool.RespInternalError(c)
		return
	}
	if !isok {
		tool.RespErrorWithData(c, "没有分享权限")
		return
	}

	qrcode.WriteFile("http://39.106.81.229:9925/file/download/"+FileId+"?share=true&expr="+Expr, qrcode.Medium, 256, "./qrcode/"+FileId+".png")
	c.File("./qrcode/" + FileId + ".png")

}

//生成加密分享二维码
func sharefilewithsecret(c *gin.Context) {
	UserId := c.MustGet("UserId").(int)
	FileId := c.Param("file_id")
	Expr := c.PostForm("expr_time")

	//校验权限
	isok, err, _ := service.CheckAuthorityToDownload(tool.StringTOInt(FileId), UserId)
	if err != nil {
		log.Println("CheckAuthorityToDownload err ", err)
		tool.RespInternalError(c)
		return
	}
	if !isok {
		tool.RespErrorWithData(c, "没有分享权限")
		return
	}
	link := "file/download/" + FileId + "?share=true&expr=" + Expr
	sha1 := sha12.New()
	io.WriteString(sha1, link)
	res := base64.StdEncoding.EncodeToString(sha1.Sum(nil))

	//入库
	err = service.AddSha1AndLinkMap(res, link)
	if err != nil {
		log.Println("AddSha1AndLinkMap err ", err)
		tool.RespInternalError(c)
		return
	}

	qrcode.WriteFile("http://39.106.81.229:9925/secert/"+res, qrcode.Medium, 256, "./qrcode/"+FileId+".png")
	c.File("./qrcode/" + FileId + ".png")

}
