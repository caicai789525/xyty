package router

import (
	"github.com/gin-gonic/gin"
	"ini/handler/login_and_signup"
	"ini/handler/upload"
	"ini/handler/user"
	"ini/router/middleware"
)

func RouterInit() *gin.Engine {
	e := gin.Default()

	e.Use(middleware.Cors())

	LoginGroup := e.Group("api/v1/")
	{
		LoginGroup.POST("/signup", login_and_signup.Signup)
		LoginGroup.POST("/pwd_login", login_and_signup.LoginWithPwd)
		LoginGroup.POST("/send_mail", login_and_signup.SendCode)
		LoginGroup.POST("/code_login", login_and_signup.LoginWithCode)
	}

	// 添加个人中心路由组
	userGroup := e.Group("api/v1/user")
	userGroup.Use(middleware.Auth()) // 需要认证
	{
		userGroup.GET("/profile", user.GetProfile)                  // 获取个人资料
		userGroup.POST("/avatar", user.UpdateAvatar)                // 更新头像
		userGroup.PUT("/password", user.ChangePassword)             // 修改密码
		userGroup.GET("/video-records", user.GetVideoRecords)       // 获取视频记录
		userGroup.POST("/video-records", user.AddVideoRecord)       // 添加视频记录
		userGroup.GET("/scenario-records", user.GetScenarioRecords) // 情景记录查询
		userGroup.POST("/scenario-records", user.AddScenarioRecord) // 情景记录添加
	}

	e.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	e.POST("/upload", upload.UploadImg)

	return e
}
