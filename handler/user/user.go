package user

import (
	"github.com/gin-gonic/gin"
	"ini/dao/mysql"
	"ini/handler"
	model "ini/model/user_struct"
	"ini/pkg/auth"
	"ini/pkg/errno"
	"ini/services/qiniu"
)

// GetProfile 获取用户个人资料
// @Summary 获取用户个人资料
// @Description 获取当前登录用户的个人资料
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} handler.Response{data=model.User}
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/user/profile [get]
// @Security ApiKeyAuth
func GetProfile(c *gin.Context) {
	claims, err := auth.ParseRequest(c)
	if err != nil {
		handler.SendError(c, errno.ErrTokenInvalid, err.Error())
		return
	}

	var user model.User
	if err := mysql.DB.Where("username = ?", claims.Username).First(&user).Error; err != nil {
		handler.SendError(c, errno.ErrUserNotFound, err.Error())
		return
	}

	// 隐藏密码信息
	user.Password = ""
	handler.SendGoodResponse(c, "获取个人资料成功", user)
}

// UpdateAvatar 更新用户头像
// @Summary 更新用户头像
// @Description 更新当前登录用户的头像
// @Tags user
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "用户头像"
// @Success 200 {object} handler.Response{data=string}
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/user/avatar [post]
// @Security ApiKeyAuth
func UpdateAvatar(c *gin.Context) {
	claims, err := auth.ParseRequest(c)
	if err != nil {
		handler.SendError(c, errno.ErrTokenInvalid, err.Error())
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		handler.SendBadResponse(c, "获取头像文件失败", err.Error())
		return
	}

	// 调用七牛云上传接口，上传到avatar文件夹
	status, url := qiniu.UploadToQiNiu(file, "avatar/"+claims.Username+"/")
	if status != 1 {
		handler.SendError(c, "上传头像失败", url)
		return
	}

	// 更新用户头像URL
	if err := mysql.DB.Model(&model.User{}).Where("username = ?", claims.Username).Update("avatar", url).Error; err != nil {
		handler.SendError(c, "更新头像失败", err.Error())
		return
	}

	handler.SendGoodResponse(c, "更新头像成功", url)
}
