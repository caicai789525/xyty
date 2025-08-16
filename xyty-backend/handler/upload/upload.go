package upload

import (
	"github.com/gin-gonic/gin"
	"ini/services/qiniu"
)

func UploadImg(c *gin.Context) {
	f, err := c.FormFile("img")
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}
	n, str := qiniu.UploadToQiNiu(f, "")
	if n != 1 {
		c.JSON(400, gin.H{"message": "出错"})
		return
	}
	c.JSON(200, gin.H{"message": str})
}
