package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"campus-swap-server/internal/model"
	"campus-swap-server/internal/service"
)

// GetCategories 获取分类列表
func GetCategories(c *gin.Context) {
	var categories []model.Category
	service.WithRead(func(d *model.Database) {
		categories = d.Categories
	})

	c.JSON(http.StatusOK, gin.H{
		"message":    "ok",
		"categories": categories,
	})
}

// GetDashboard 后台统计面板
func GetDashboard(c *gin.Context) {
	var stats gin.H
	service.WithRead(func(d *model.Database) {
		pendingReports := 0
		for _, r := range d.Reports {
			if r.Status == "pending" {
				pendingReports++
			}
		}
		stats = gin.H{
			"userCount":      len(d.Users),
			"itemCount":      len(d.Items),
			"exchangeCount":  len(d.Exchanges),
			"reportCount":    len(d.Reports),
			"pendingReports": pendingReports,
			"messageCount":   len(d.Messages),
		}
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"stats":   stats,
	})
}

// GetRecommendations 推荐物品（基于收藏偏好）
func GetRecommendations(c *gin.Context) {
	userId := c.GetString("userId")

	var recommendations []model.Item
	service.WithRead(func(d *model.Database) {
		// 获取用户收藏的分类偏好
		categoryPrefs := make(map[string]bool)
		for _, fav := range d.Favorites {
			if fav.UserID == userId {
				for _, item := range d.Items {
					if item.ID == fav.ItemID {
						categoryPrefs[item.Category] = true
						break
					}
				}
			}
		}

		// 根据偏好推荐物品
		for _, item := range d.Items {
			if item.Status != "available" || item.OwnerID == userId {
				continue
			}
			if len(categoryPrefs) == 0 || categoryPrefs[item.Category] {
				recommendations = append(recommendations, item)
			}
		}

		// 如果没有偏好，返回所有可用物品（最多10个）
		if len(recommendations) > 10 {
			recommendations = recommendations[:10]
		}
	})

	if recommendations == nil {
		recommendations = []model.Item{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"items":   recommendations,
	})
}

// CreateReport 创建举报
func CreateReport(c *gin.Context) {
	userId := c.GetString("userId")

	var req struct {
		TargetID string `json:"targetId" binding:"required"`
		Type     string `json:"type" binding:"required"`
		Reason   string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请填写举报信息"})
		return
	}

	newReport := model.Report{
		ID:        service.GenID("rpt"),
		UserID:    userId,
		TargetID:  req.TargetID,
		Type:      req.Type,
		Reason:    req.Reason,
		Status:    "pending",
		CreatedAt: model.TimeNow(),
	}

	service.WithWrite(func(d *model.Database) {
		d.Reports = append(d.Reports, newReport)
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "举报已提交",
		"report":  newReport,
	})
}

// CreateRating 创建评价
func CreateRating(c *gin.Context) {
	userId := c.GetString("userId")

	var req struct {
		ExchangeID string `json:"exchangeId" binding:"required"`
		ToUserID   string `json:"toUserId" binding:"required"`
		Score      int    `json:"score" binding:"required"`
		Comment    string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请填写评价信息"})
		return
	}

	if req.Score < 1 || req.Score > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "评分需在1-5之间"})
		return
	}

	newRating := model.Rating{
		ID:         service.GenID("rat"),
		ExchangeID: req.ExchangeID,
		FromUserID: userId,
		ToUserID:   req.ToUserID,
		Score:      req.Score,
		Comment:    req.Comment,
		CreatedAt:  model.TimeNow(),
	}

	service.WithWrite(func(d *model.Database) {
		d.Ratings = append(d.Ratings, newRating)
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "评价成功",
		"rating":  newRating,
	})
}

// GetNotifications 获取通知列表（当前返回空列表）
func GetNotifications(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":       "ok",
		"notifications": []interface{}{},
	})
}
