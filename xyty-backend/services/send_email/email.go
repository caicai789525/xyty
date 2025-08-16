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

func (e *Email) Send(mail string) int {
	v := parseyaml.GetYaml()
	addr := v.GetString("addr")
	pwd := v.GetString("password")
	auth := smtp.PlainAuth("", addr, pwd, "smtp.163.com")
	to := []string{mail}
	rand.Seed(time.Now().Unix())
	num := rand.Intn(900000) + 100000
	str := fmt.Sprintf("From:%s\r\nTo:%s\r\nSubject:您的验证码是\r\n\r\n%d\r\n", addr, mail, num) //邮件格式
	msg := []byte(str)
	err := smtp.SendMail("smtp.163.com:25", auth, "SJMbaiyang@163.com", to, msg)
	if err != nil {
		log.Fatal(err)
	}
	return num
}
