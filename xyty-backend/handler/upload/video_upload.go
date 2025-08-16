package upload

import (
	"ini/handler"
	"ini/pkg/auth"
	"ini/pkg/errno"
	"ini/services/qiniu"

	"github.com/gin-gonic/gin"
)

// UploadVideo 上传视频
// @Summary 上传视频
// @Description 上传视频文件
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param video formData file true "视频文件"
// @Param scenario_id formData string true "情景ID"
// @Success 200 {object} handler.Response{data=string}
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/upload/video [post]
// @Security ApiKeyAuth
func UploadVideo(c *gin.Context) {
	// 验证用户
	_, err := auth.ParseRequest(c)
	if err != nil {
		handler.SendError(c, errno.ErrTokenInvalid, err.Error())
		return
	}

	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	scenarioID := c.PostForm("scenario_id")
	if scenarioID == "" {
		c.JSON(400, gin.H{"message": "情景ID不能为空"})
		return
	}

	// 上传视频到七牛云，使用scenario_id作为文件夹名
	n, str := qiniu.UploadToQiNiu(file, "videos/"+scenarioID+"/")
	if n != 1 {
		c.JSON(400, gin.H{"message": "视频上传失败", "error": str})
		return
	}

	c.JSON(200, gin.H{"message": "视频上传成功", "video_url": str})
}
