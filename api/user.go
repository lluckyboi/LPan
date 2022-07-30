package api

import (
	"LPan/dao"
	"LPan/model"
	"LPan/service"
	"LPan/tool"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

//验证码登录，如果未注册直接注册
func login(c *gin.Context) {
	UserMail := c.PostForm("UserMail")
	AcCode := c.PostForm("AccessCode")
	//checkAcCode
	code := dao.RedisDB.Get("mkey" + AcCode).String()
	log.Println(code)

	//解析
	cidx := 0
	for {
		n := len(code[:])
		if code[cidx] != ':' {
			cidx++
			if cidx > n {
				break
			}
			continue
		}
		cidx += 2
		break
	}

	if code[cidx:] != AcCode {
		c.JSON(200, gin.H{
			"code": 201,
			"err":  "验证码错误或过期",
		})
		return
	}

	//check
	isnok, err := service.IsUserExistByMail(UserMail)
	if err != nil {
		tool.RespInternalError(c)
		log.Println(err)
		return
	}

	//可以创建新账户
	if isnok {
		User := model.User{
			UserMail: UserMail,
		}
		err = service.NewUser(User)
		if err != nil {
			tool.RespInternalError(c)
			log.Println(err)
			return
		}
	}

	User, errr := service.GetUserInfoByMail(UserMail)
	if errr != nil {
		tool.RespInternalError(c)
		log.Println(errr)
		return
	}

	//JWT
	AcUser := model.AcUser{
		UserId: User.UserId,
		AcCode: AcCode,
	}

	AcUser.UserId = User.UserId
	// 生成Token
	tokenString, _ := model.GenToken(AcUser)
	c.JSON(http.StatusOK, gin.H{
		"name":   User.UserName,
		"status": 200,
		"token":  tokenString,
	})
	return
}

//获取邮箱验证码
func getmailac(c *gin.Context) {
	mail := c.PostForm("UserMail")
	var mails []string
	mails = append(mails, mail)
	err, bl := SendMail(mails)
	if err != nil {
		tool.RespInternalError(c)
		log.Println(err)
	} else if bl == true {
		c.JSON(200, gin.H{
			"code": 200,
			"info": "Send Message Success",
		})
	} else {
		c.JSON(200, gin.H{
			"code": 203,
			"info": "Send times limited!",
		})
	}
}

func updateinfo(c *gin.Context) {
	UserName := c.PostForm("UserName")
	User := model.User{
		UserName: UserName,
	}
	_, bl := tool.LengthCheck(UserName)
	if !bl {
		tool.RespErrorWithData(c, "用户名长度不符合要求")
		return
	}

	err := service.UpdateUserName(User.UserName)
	if err != nil {
		tool.RespInternalError(c)
		log.Println(err)
		return
	}
	tool.RespSuccessful(c, "修改用户名")
}
