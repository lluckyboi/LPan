package api

import (
	"LPan/dao"
	"LPan/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	QQMailSMTPCode = "vxngdjqoeaoajajd"
	QQMailSender   = "1598273095@qq.com"
	QQMailTitle    = "验证"
	SMTPAdr        = "smtp.qq.com"
	SMTPPort       = 587
	MailListSize   = 2048

	SecretId    = "AKIDSjlgKhpgFXYIkvpRmqc8MK1ScU5PSGo4"
	SecreKey    = "1lDOJxKQFn44nFPItDH4ZnkGWUCdmeFZ"
	SmsSdkAppId = "1400696970"
	TemplateId  = "1452825"
)

type MailboxConf struct {
	// 邮件标题
	Title string
	// 邮件内容
	Body string
	// 收件人列表
	RecipientList []string
	// 发件人账号
	Sender string
	// 发件人密码，QQ邮箱这里配置授权码
	SPassword string
	// SMTP 服务器地址， QQ邮箱是smtp.qq.com
	SMTPAddr string
	// SMTP端口 QQ邮箱是25
	SMTPPort int
}

// SendMail QQ邮箱验证码
func SendMail(mails []string) (error, bool) {
	var mailConf MailboxConf
	mailConf.Title = QQMailTitle
	mailConf.RecipientList = mails
	mailConf.Sender = QQMailSender
	mailConf.SPassword = QQMailSMTPCode
	mailConf.SMTPAddr = SMTPAdr
	mailConf.SMTPPort = SMTPPort

	//产生六位数验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	//发送的内容
	html := fmt.Sprintf(`<div>
        <div>
            尊敬的用户，您好！
        </div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>你本次的验证码为%s,为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
        </div>
        <div>
            <p>此邮箱为系统邮箱，请勿回复。</p>
        </div>
    </div>`, vcode)

	m := gomail.NewMessage()

	m.SetHeader(`From`, mailConf.Sender, "LPan")
	m.SetHeader(`To`, mailConf.RecipientList...)
	m.SetHeader(`Subject`, mailConf.Title)
	m.SetBody(`text/html`, html)

	err := gomail.NewDialer(mailConf.SMTPAddr, mailConf.SMTPPort, mailConf.Sender, mailConf.SPassword).DialAndSend(m)
	if err != nil {
		log.Fatalf("Send Email Fail, %s", err.Error())
		return err, false
	}
	for _, j := range mails {
		//	1.定义验证码和手机号的key的形式
		codeKey := "mkey" + vcode
		mailKey := "mkey" + j
		//	2.查看用户发送验证码的数量
		ctn := dao.RedisDB.Get(mailKey).String()
		log.Println(ctn)
		//  ctn为零值的时候，说明该用户还没有发送过验证码
		if ctn == "" {
			dao.RedisDB.Set(mailKey, "1", time.Minute*60*24) //设置手请求验证码的时间为一天。
		}
		cidx := 0
		for {
			n := len(ctn[:])
			if ctn[cidx] != ':' {
				cidx++
				if cidx > n {
					break
				}
				continue
			}
			cidx += 2
			break
		}
		ctn1, _ := strconv.ParseInt(ctn[cidx:], 10, 32)
		//log.Println(ctn1)
		//  发送次数没有超过100次
		if ctn1 <= 99 {
			dao.RedisDB.Incr(mailKey)
		} else {
			return nil, false
		}
		//	将验证码存入redis,设置过期时间为2分钟
		dao.RedisDB.Set(codeKey, vcode, time.Minute*2)
		//MalilList[j] = vcode
	}
	log.Printf("Send Email Success")
	return nil, true
}

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		mc, err := model.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"info": "无效的Token",
			})
			c.Abort()
			return
		}
		c.Set("UserId", mc.UserId)
		c.Next()
	}
}

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
