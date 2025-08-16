package method

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/jordan-wright/email"
	"ini/dao/mysql"
	model2 "ini/model"
	model "ini/model/user_struct"
	"ini/pkg/errno"
	"ini/router/middleware"
	"ini/services/parseyaml"
	"log"
	"math/big"
	"net/smtp"
)

// 账号密码登录
func UserLogin1(LoginUser model.User) (string, error, int) {
	enc_str := base64.StdEncoding.EncodeToString([]byte(LoginUser.Password))
	err := mysql.DB.Where("username = ? AND password = ?", LoginUser.Username, enc_str).First(&LoginUser).Error
	if err != nil {
		return "", err, 1
	}
	signedToken, errr := middleware.GenerateToken(LoginUser.Username)
	if err != nil {
		log.Println(errr)
		return "", errr, 2
	}
	return signedToken, nil, 0
}

// 随机生成验证码
func generateRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	randomString := make([]rune, length)

	for i := range randomString {
		// 生成随机索引
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		randomString[i] = letters[n.Int64()]
	}

	return string(randomString)
}

// 发送验证邮箱
func SendLoginCode(addr string) error {
	err := mysql.DB.Where("mail_addr = ?", addr).First(&model.User{}).Error
	if err != nil {
		return errno.ErrMailNotExist
	}
	code := generateRandomString(6)
	// 简单设置 log 参数
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	v := parseyaml.GetYaml()

	em := email.NewEmail()
	// 设置 sender 发送方 的邮箱
	em.From = v.GetString("mail.addr")

	// 设置 receiver 接收方 的邮箱
	em.To = []string{addr}

	// 设置主题
	em.Subject = "偶遇华夏登录验证码"

	// 简单设置文件发送的内容
	em.Text = []byte(code)

	//设置服务器相关的配置
	err = em.Send("smtp.163.com:25", smtp.PlainAuth("", v.GetString("mail.addr"), v.GetString("mail.auth_code"), v.GetString("mail.smtp_addr")))
	if err != nil {
		return err
	}

	//更新数据库中信息
	err = mysql.DB.Update("code", code).Where("mail_addr = ?", addr).Error
	if err != nil {
		return err
	}
	return nil
}

func UserLogin2(LoginU model2.LoginU2) (string, error, int) {
	var Login model.User
	err := mysql.DB.Where("mail_addr = ?", LoginU.MailAddr).First(&Login).Error
	if err != nil {
		return "", err, 1
	}
	err = mysql.DB.Where("mail_addr = ? AND code = ?", LoginU.MailAddr, LoginU.Code).First(&model.UserCode{}).Error
	if err != nil {
		return "", err, 2
	}
	signedToken, errr := middleware.GenerateToken(Login.Username)
	if err != nil {
		log.Println(errr)
		return "", errr, 3
	}
	return signedToken, nil, 0
}
