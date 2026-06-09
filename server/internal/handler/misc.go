package handler

import (
	"net/http"
	"time"

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

		// 热门分类统计
		categoryCount := make(map[string]int)
		for _, item := range d.Items {
			if item.Status == "available" {
				categoryCount[item.Category]++
			}
		}
		type hotCategory struct {
			Category string `json:"category"`
			Name     string `json:"name"`
			Count    int    `json:"count"`
		}
		var hotCategories []hotCategory
		for catID, count := range categoryCount {
			name := catID
			for _, cat := range d.Categories {
				if cat.ID == catID {
					name = cat.Name
					break
				}
			}
			hotCategories = append(hotCategories, hotCategory{
				Category: catID,
				Name:     name,
				Count:    count,
			})
		}
		// 按数量降序排序
		for i := 0; i < len(hotCategories); i++ {
			for j := i + 1; j < len(hotCategories); j++ {
				if hotCategories[j].Count > hotCategories[i].Count {
					hotCategories[i], hotCategories[j] = hotCategories[j], hotCategories[i]
				}
			}
		}

		stats = gin.H{
			"userCount":      len(d.Users),
			"itemCount":      len(d.Items),
			"exchangeCount":  len(d.Exchanges),
			"reportCount":    len(d.Reports),
			"pendingReports": pendingReports,
			"messageCount":   len(d.Messages),
			"hotCategories":  hotCategories,
		}
	})

	c.JSON(http.StatusOK, gin.H{
		"message":   "ok",
		"stats":     stats,
		"dashboard": stats,
	})
}

// GetRecommendations 推荐物品（基于收藏+浏览的加权分类偏好）
func GetRecommendations(c *gin.Context) {
	userId := c.GetString("userId")

	type ScoredItem struct {
		Item  model.Item
		Score float64
	}

	var results []ScoredItem
	service.WithRead(func(d *model.Database) {
		// 1. 计算分类偏好权重
		//    收藏同一分类的物品：每个 +5 分
		//    浏览同一分类的物品：每个 +1 分
		categoryPrefs := make(map[string]float64)

		for _, fav := range d.Favorites {
			if fav.UserID == userId {
				for _, item := range d.Items {
					if item.ID == fav.ItemID {
						categoryPrefs[item.Category] += 5
						break
					}
				}
			}
		}

		for _, br := range d.BrowseRecords {
			if br.UserID == userId {
				categoryPrefs[br.Category] += 1
			}
		}

		// 2. 对每个可用物品计算推荐分数
		for _, item := range d.Items {
			if item.Status != "available" || item.OwnerID == userId {
				continue
			}

			score := categoryPrefs[item.Category] // 分类偏好分

			// 时间衰减：越新越好，7天内加分
			if item.CreatedAt != "" {
				createdTime, err := time.Parse("2006-01-02 15:04:05", item.CreatedAt)
				if err == nil {
					hoursAgo := time.Since(createdTime).Hours()
					if hoursAgo < 24 {
						score += 3 // 1天内
					} else if hoursAgo < 72 {
						score += 2 // 3天内
					} else if hoursAgo < 168 {
						score += 1 // 7天内
					}
				}
			}

			// 没有任何偏好时，浏览量高的排前面
			if len(categoryPrefs) == 0 {
				score = float64(item.Views)*0.1 + 1
			}

			results = append(results, ScoredItem{Item: item, Score: score})
		}

		// 3. 按分数降序排序
		for i := 0; i < len(results); i++ {
			for j := i + 1; j < len(results); j++ {
				if results[j].Score > results[i].Score {
					results[i], results[j] = results[j], results[i]
				}
			}
		}

		// 4. 最多返回 10 个
		if len(results) > 10 {
			results = results[:10]
		}
	})

	// 提取物品列表
	var recommendations []model.Item
	for _, r := range results {
		recommendations = append(recommendations, r.Item)
	}

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
		TargetID   string `json:"targetId" binding:"required"`
		Type       string `json:"type"`
		TargetType string `json:"targetType"`
		Reason     string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请填写举报信息"})
		return
	}

	// 兼容 type 和 targetType 两种字段名
	reportType := req.Type
	if reportType == "" {
		reportType = req.TargetType
	}
	if reportType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请指定举报类型"})
		return
	}

	newReport := model.Report{
		ID:        service.GenID("rpt"),
		UserID:    userId,
		TargetID:  req.TargetID,
		Type:      reportType,
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

// GetReports 获取举报列表（管理员）
func GetReports(c *gin.Context) {
	var reports []model.Report
	service.WithRead(func(d *model.Database) {
		reports = d.Reports
	})

	if reports == nil {
		reports = []model.Report{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"reports": reports,
	})
}

// UpdateReport 更新举报状态（管理员）
func UpdateReport(c *gin.Context) {
	reportId := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请指定状态"})
		return
	}

	if req.Status != "pending" && req.Status != "resolved" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "状态必须是 pending 或 resolved"})
		return
	}

	var updated bool
	service.WithWrite(func(d *model.Database) {
		for i := range d.Reports {
			if d.Reports[i].ID == reportId {
				d.Reports[i].Status = req.Status
				updated = true
				break
			}
		}
	})

	if !updated {
		c.JSON(http.StatusNotFound, gin.H{"message": "未找到举报记录"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// GetRatings 获取评价列表
func GetRatings(c *gin.Context) {
	toUserId := c.Query("toUserId")
	userId := c.Query("userId") // 兼容前端传 userId 参数

	// toUserId 优先，否则用 userId
	targetUserId := toUserId
	if targetUserId == "" {
		targetUserId = userId
	}

	var ratings []model.Rating
	service.WithRead(func(d *model.Database) {
		if targetUserId != "" {
			for _, r := range d.Ratings {
				if r.ToUserID == targetUserId {
					ratings = append(ratings, r)
				}
			}
		} else {
			ratings = d.Ratings
		}
	})

	if ratings == nil {
		ratings = []model.Rating{}
	}

	// 计算平均分和评价数
	var avgScore float64
	var ratingCount int
	if len(ratings) > 0 {
		var total int
		for _, r := range ratings {
			total += r.Score
		}
		avgScore = float64(total) / float64(len(ratings))
		ratingCount = len(ratings)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "ok",
		"ratings":     ratings,
		"avgScore":    avgScore,
		"ratingCount": ratingCount,
	})
}

// GetNotifications 获取通知列表
func GetNotifications(c *gin.Context) {
	userId := c.GetString("userId")

	var notifications []model.Notification
	service.WithRead(func(d *model.Database) {
		for _, n := range d.Notifications {
			if n.UserID == userId {
				notifications = append(notifications, n)
			}
		}
	})

	// 按时间倒序
	for i := 0; i < len(notifications)/2; i++ {
		notifications[i], notifications[len(notifications)-1-i] = notifications[len(notifications)-1-i], notifications[i]
	}

	if notifications == nil {
		notifications = []model.Notification{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "ok",
		"notifications": notifications,
	})
}

// MarkNotificationRead 标记通知为已读
func MarkNotificationRead(c *gin.Context) {
	userId := c.GetString("userId")
	notificationId := c.Param("id")

	service.WithWrite(func(d *model.Database) {
		for i := range d.Notifications {
			if d.Notifications[i].ID == notificationId && d.Notifications[i].UserID == userId {
				d.Notifications[i].IsRead = true
				break
			}
		}
	})

	c.JSON(http.StatusOK, gin.H{"message": "已标记已读"})
}

// MarkAllNotificationsRead 标记所有通知为已读
func MarkAllNotificationsRead(c *gin.Context) {
	userId := c.GetString("userId")

	service.WithWrite(func(d *model.Database) {
		for i := range d.Notifications {
			if d.Notifications[i].UserID == userId {
				d.Notifications[i].IsRead = true
			}
		}
	})

	c.JSON(http.StatusOK, gin.H{"message": "已全部标记已读"})
}

// GetUnreadNotificationCount 获取未读通知数
func GetUnreadNotificationCount(c *gin.Context) {
	userId := c.GetString("userId")

	var count int
	service.WithRead(func(d *model.Database) {
		for _, n := range d.Notifications {
			if n.UserID == userId && !n.IsRead {
				count++
			}
		}
	})

	c.JSON(http.StatusOK, gin.H{"unreadCount": count})
}

// CreateNotification 创建通知（内部方法）
func createNotification(userId, notifType, title, content, relatedId string) {
	newNotif := model.Notification{
		ID:        service.GenID("notif"),
		UserID:    userId,
		Type:      notifType,
		Title:     title,
		Content:   content,
		RelatedID: relatedId,
		IsRead:    false,
		CreatedAt: model.TimeNow(),
	}

	service.WithWrite(func(d *model.Database) {
		d.Notifications = append(d.Notifications, newNotif)
	})
}
