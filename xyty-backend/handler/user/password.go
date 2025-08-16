package user

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"ini/dao/mysql"
	"ini/handler"
	model "ini/model/user_struct"
	"ini/pkg/auth"
	"ini/pkg/errno"
)

// ChangePasswordRequest 修改密码请求结构
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags user
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} handler.Response
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/user/password [put]
// @Security ApiKeyAuth
func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.BindJSON(&req); err != nil {
		handler.SendBadResponse(c, "请求参数错误", err.Error())
		return
	}

	// 解析JWT获取当前用户
	claims, err := auth.ParseRequest(c)
	if err != nil {
		handler.SendError(c, errno.ErrTokenInvalid, err.Error())
		return
	}

	// 查询用户信息
	var user model.User
	if err := mysql.DB.Where("username = ?", claims.Username).First(&user).Error; err != nil {
		handler.SendError(c, errno.ErrUserNotFound, err.Error())
		return
	}

	// 验证原密码（与数据库中存储的Base64加密密码比对）
	oldPasswordEnc := base64.StdEncoding.EncodeToString([]byte(req.OldPassword))
	if user.Password != oldPasswordEnc {
		handler.SendBadResponse(c, "原密码错误", nil)
		return
	}

	// 加密新密码并更新到数据库
	newPasswordEnc := base64.StdEncoding.EncodeToString([]byte(req.NewPassword))
	if err := mysql.DB.Model(&user).Update("password", newPasswordEnc).Error; err != nil {
		handler.SendError(c, "修改密码失败", err.Error())
		return
	}

	handler.SendGoodResponse(c, "密码修改成功", nil)
}
