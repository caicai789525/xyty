package send_email

import (
	"fmt"
	"ini/services/parseyaml"
	"log"
	"math/rand"
	"net/smtp"
	"time"
)

type Email struct{}

// Send 发送验证码
func (e *Email) Send(mail string) int {
	// 解析yaml配置文件
	v := parseyaml.GetYaml()
	addr := v.GetString("addr")
	pwd := v.GetString("password")
	// 配置SMTP服务器
	auth := smtp.PlainAuth("", addr, pwd, "smtp.163.com")
	to := []string{mail}
	// 生成6位随机验证码
	rand.Seed(time.Now().Unix())
	num := rand.Intn(900000) + 100000
	str := fmt.Sprintf("From:%s\r\nTo:%s\r\nSubject:您的验证码是\r\n\r\n%d\r\n", addr, mail, num) //邮件格式
	msg := []byte(str)
	// 发送邮件
	err := smtp.SendMail("smtp.163.com:25", auth, "SJMbaiyang@163.com", to, msg)
	if err != nil {
		log.Fatal(err)
	}
	return num
}
