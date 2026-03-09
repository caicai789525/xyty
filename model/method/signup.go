package method

import (
	"encoding/base64"
	"ini/dao/mysql"
	model2 "ini/model"
	model "ini/model/user_struct"
	"ini/pkg/errno"
)

func Signup(SignupU model2.SignupUser) error {
	var Sign model.User
	enc_str := base64.StdEncoding.EncodeToString([]byte(SignupU.Password))
	Sign = model.User{
		Username: SignupU.Username,
		Password: enc_str,
		MailAddr: SignupU.MailAddr,
	}
	err := mysql.DB.Where("username = ?", SignupU.Username).First(&model.User{}).Error
	if err == nil {
		return errno.ErrUsernameExist
	}
	err = mysql.DB.Where("mail_addr = ?", SignupU.MailAddr).First(&model.User{}).Error
	if err == nil {
		return errno.ErrMailUsed
	}
	err = mysql.DB.Create(&Sign).Error
	if err != nil {
		return errno.ErrDatabase
	}
	err = mysql.DB.Create(&model.UserCode{
		Username: SignupU.Username,
		MailAddr: SignupU.MailAddr,
		Code:     "",
	}).Error
	if err != nil {
		return errno.ErrDatabase
	}
	return nil
}
