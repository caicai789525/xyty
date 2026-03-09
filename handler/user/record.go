package user

import (
	"ini/dao/mysql"
	"ini/handler"
	model "ini/model/user_struct"
	"ini/pkg/auth"
	"ini/pkg/errno"
	"ini/services/qiniu"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetVideoRecords 获取视频观看记录
// @Summary 获取视频观看记录
// @Description 获取当前登录用户的视频观看记录
// @Tags user
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} handler.Response{data=[]model.VideoRecord}
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/user/video-records [get]
// @Security ApiKeyAuth
func GetVideoRecords(c *gin.Context) {
	claims, err := auth.ParseRequest(c)
	if err != nil {
		handler.SendError(c, errno.ErrTokenInvalid, err.Error())
		return
	}

	var videoRecords []model.VideoRecord
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	// 转换分页参数为整数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // 默认为第1页
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10 // 默认为10条/页，最大100条
	}

	// 计算偏移量，实现分页查询
	offset := (page - 1) * size
	if err := mysql.DB.Where("username = ?", claims.Username).
		Order("created_at DESC"). // 按创建时间倒序
		Limit(size). // 限制每页条数
		Offset(offset). // 偏移量（跳过前面的记录）
		Find(&videoRecords).Error; err != nil {
		handler.SendError(c, errno.ErrDatabase, err.Error())
		return
	}

	handler.SendGoodResponse(c, "获取视频记录成功", videoRecords)
}

// AddVideoRecord 添加视频观看记录
// @Summary 添加视频观看记录
// @Description 添加视频观看记录
// @Tags user
// @Accept json
// @Produce json
// @Param request body model.VideoRecord true "视频记录信息"  // 修正：将 user_struct.VideoRecord 改为 model.VideoRecord
// @Success 200 {object} handler.Response
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/user/video-records [post]
func AddVideoRecord(c *gin.Context) {
	claims, err := auth.ParseRequest(c)
	if err != nil {
		handler.SendError(c, errno.ErrTokenInvalid, err.Error())
		return
	}

	var record model.VideoRecord
	if err := c.BindJSON(&record); err != nil {
		handler.SendBadResponse(c, "请求参数错误", err.Error())
		return
	}

	record.Username = claims.Username
	record.WatchDate = time.Now().Format("2006-01-02")

	if err := mysql.DB.Create(&record).Error; err != nil {
		handler.SendError(c, errno.ErrDatabase, err.Error())
		return
	}

	handler.SendGoodResponse(c, "添加视频记录成功", nil)
}

// GetQiniuVideos 获取七牛云视频列表
// @Summary 获取七牛云视频列表
// @Description 根据情景ID获取七牛云存储的视频文件列表
// @Tags user
// @Accept json
// @Produce json
// @Param scenario_id query string true "情景ID"
// @Success 200 {object} handler.Response{data=[]string}
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/user/qiniu-videos [get]
// @Security ApiKeyAuth
func GetQiniuVideos(c *gin.Context) {
	_, err := auth.ParseRequest(c)
	if err != nil {
		handler.SendError(c, errno.ErrTokenInvalid, err.Error())
		return
	}

	scenarioID := c.Query("scenario_id")
	if scenarioID == "" {
		handler.SendError(c, errno.ErrQuery, "情景ID不能为空")
		return
	}

	prefix := "drifting/" + scenarioID + "/"
	urls, status, err := qiniu.ListFilesByPrefix(prefix)
	if err != nil || status != 1 {
		handler.SendError(c, errno.ErrDatabase, err.Error())
		return
	}

	handler.SendGoodResponse(c, "获取视频列表成功", urls)
}

// GetScenarioRecords 获取情景体验记录
// @Summary 获取情景体验记录
// @Description 获取当前登录用户的情景体验记录
// @Tags user
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(10)
// @Success 200 {object} handler.Response{data=[]model.ScenarioExperience}  // 修正：ScenarioRecord → ScenarioExperience
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/user/scenario-records [get]
// @Security ApiKeyAuth
func GetScenarioRecords(c *gin.Context) {
	claims, err := auth.ParseRequest(c)
	if err != nil {
		handler.SendError(c, errno.ErrTokenInvalid, err.Error())
		return
	}

	var scenarioRecords []model.ScenarioExperience
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	// 转换分页参数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	// 分页查询
	offset := (page - 1) * size
	if err := mysql.DB.Where("username = ?", claims.Username).
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&scenarioRecords).Error; err != nil {
		handler.SendError(c, errno.ErrDatabase, err.Error())
		return
	}

	handler.SendGoodResponse(c, "获取情景体验记录成功", scenarioRecords)
}

// AddScenarioRecord 添加情景体验记录
// @Summary 添加情景体验记录
// @Description 添加当前登录用户的情景体验记录
// @Tags user
// @Accept json
// @Produce json
// @Param request body model.ScenarioExperience true "情景体验记录信息"  // 修正：ScenarioRecord → ScenarioExperience
// @Success 200 {object} handler.Response
// @Failure 400 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/user/scenario-records [post]
// @Security ApiKeyAuth
func AddScenarioRecord(c *gin.Context) {
	claims, err := auth.ParseRequest(c)
	if err != nil {
		handler.SendError(c, errno.ErrTokenInvalid, err.Error())
		return
	}

	var record model.ScenarioExperience
	if err := c.BindJSON(&record); err != nil {
		handler.SendBadResponse(c, "请求参数错误", err.Error())
		return
	}

	// 关联当前用户和体验日期
	record.Username = claims.Username
	record.ExperienceDate = time.Now().Format("2006-01-02")

	if err := mysql.DB.Create(&record).Error; err != nil {
		handler.SendError(c, errno.ErrDatabase, err.Error())
		return
	}

	handler.SendGoodResponse(c, "添加情景体验记录成功", nil)
}

// @Summary 获取所有视频列表
// @Tags user
// @Success 200 {object} handler.Response{data=[]string}
// @Router /user/videos [get]
func GetAllVideos(c *gin.Context) {
	// 无需权限验证，直接查询所有视频
	urls, _, err := qiniu.ListFilesByPrefix("")
	if err != nil {
		handler.SendError(c, errno.ErrDatabase, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    urls,
	})
}
