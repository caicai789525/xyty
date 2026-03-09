package model

// 注册用户
type SignupUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	MailAddr string `json:"mail_addr"`
}

// 账号密码登录
type LoginU1 struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 邮箱验证码登录
type LoginU2 struct {
	MailAddr string `json:"mail_addr"`
	Code     string `json:"code"`
}

// 传输邮箱
type SendMailAddr struct {
	MailAddr string `json:"mail_addr"`
}
