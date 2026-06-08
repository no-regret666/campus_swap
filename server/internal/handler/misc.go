package handler

import (
	"math/rand"
	"net/http"
	"sort"

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

// GetDashboard 后台统计面板（管理员）
func GetDashboard(c *gin.Context) {
	var result gin.H
	service.WithRead(func(d *model.Database) {
		pendingReports := 0
		for _, r := range d.Reports {
			if r.Status == "pending" {
				pendingReports++
			}
		}

		// 交换完成率
		completedExchanges := 0
		for _, ex := range d.Exchanges {
			if ex.Status == "completed" {
				completedExchanges++
			}
		}
		exchangeRate := 0.0
		if len(d.Exchanges) > 0 {
			exchangeRate = float64(completedExchanges) / float64(len(d.Exchanges)) * 100
		}

		// 热门分类排行
		categoryCount := make(map[string]int)
		for _, item := range d.Items {
			categoryCount[item.Category]++
		}
		type catRank struct {
			Category string `json:"category"`
			Name     string `json:"name"`
			Count    int    `json:"count"`
		}
		var hotCategories []catRank
		for _, cat := range d.Categories {
			if cnt, ok := categoryCount[cat.ID]; ok && cnt > 0 {
				hotCategories = append(hotCategories, catRank{
					Category: cat.ID,
					Name:     cat.Name,
					Count:    cnt,
				})
			}
		}
		sort.Slice(hotCategories, func(i, j int) bool {
			return hotCategories[i].Count > hotCategories[j].Count
		})

		// 最近审计日志（最多20条）
		recentLogs := d.AuditLogs
		if len(recentLogs) > 20 {
			recentLogs = recentLogs[len(recentLogs)-20:]
		}

		result = gin.H{
			"userCount":      len(d.Users),
			"itemCount":      len(d.Items),
			"exchangeCount":  len(d.Exchanges),
			"reportCount":    len(d.Reports),
			"pendingReports": pendingReports,
			"messageCount":   len(d.Messages),
			"exchangeRate":   exchangeRate,
			"hotCategories":  hotCategories,
			"auditLogs":      recentLogs,
		}
	})

	c.JSON(http.StatusOK, gin.H{
		"message":   "ok",
		"dashboard": result,
	})
}

// recommendItem 推荐物品（带推荐理由）
type recommendItem struct {
	model.Item
	Reason string `json:"reason"`
}

// GetRecommendations 推荐物品（基于收藏偏好 + 浏览量加权 + 随机打乱）
func GetRecommendations(c *gin.Context) {
	userId := c.GetString("userId")

	var recommendations []recommendItem
	service.WithRead(func(d *model.Database) {
		// 获取用户收藏的分类偏好
		categoryPrefs := make(map[string]bool)
		favoritedItems := make(map[string]bool)
		for _, fav := range d.Favorites {
			if fav.UserID == userId {
				favoritedItems[fav.ItemID] = true
				for _, item := range d.Items {
					if item.ID == fav.ItemID {
						categoryPrefs[item.Category] = true
						break
					}
				}
			}
		}

		// 收集可推荐物品
		for _, item := range d.Items {
			if item.Status != "available" || item.OwnerID == userId {
				continue
			}
			// 排除已收藏物品
			if favoritedItems[item.ID] {
				continue
			}

			reason := ""
			if len(categoryPrefs) > 0 && categoryPrefs[item.Category] {
				reason = "根据你的收藏偏好推荐"
			} else if item.Views >= 10 {
				reason = "热门物品，浏览量高"
			} else {
				reason = "为你发现的好物"
			}

			// 有偏好时只推荐偏好分类，无偏好时推荐全部
			if len(categoryPrefs) > 0 && !categoryPrefs[item.Category] {
				continue
			}

			recommendations = append(recommendations, recommendItem{
				Item:   item,
				Reason: reason,
			})
		}

		// 按浏览量降序排序
		sort.Slice(recommendations, func(i, j int) bool {
			return recommendations[i].Views > recommendations[j].Views
		})

		// 随机打乱前半部分，增加多样性
		if len(recommendations) > 4 {
			topN := len(recommendations) / 2
			if topN > 10 {
				topN = 10
			}
			rand.Shuffle(topN, func(i, j int) {
				recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
			})
		}

		// 最多返回10个
		if len(recommendations) > 10 {
			recommendations = recommendations[:10]
		}
	})

	if recommendations == nil {
		recommendations = []recommendItem{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "ok",
		"recommendations": recommendations,
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

	// 验证 type 合法性
	if req.Type != "item" && req.Type != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "举报类型只能为 item 或 user"})
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

// GetReports 获取举报列表（管理员）
func GetReports(c *gin.Context) {
	userId := c.GetString("userId")

	// 验证管理员权限
	isAdmin := false
	service.WithRead(func(d *model.Database) {
		for _, u := range d.Users {
			if u.ID == userId && u.Role == "admin" {
				isAdmin = true
				break
			}
		}
	})
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "仅管理员可访问"})
		return
	}

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

// UpdateReport 处理举报（管理员）
func UpdateReport(c *gin.Context) {
	userId := c.GetString("userId")
	id := c.Param("id")

	// 验证管理员权限
	isAdmin := false
	service.WithRead(func(d *model.Database) {
		for _, u := range d.Users {
			if u.ID == userId && u.Role == "admin" {
				isAdmin = true
				break
			}
		}
	})
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "仅管理员可操作"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请填写处理状态"})
		return
	}

	if req.Status != "resolved" && req.Status != "dismissed" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "状态只能为 resolved 或 dismissed"})
		return
	}

	var found bool
	var report model.Report
	service.WithWrite(func(d *model.Database) {
		for i := range d.Reports {
			if d.Reports[i].ID == id {
				d.Reports[i].Status = req.Status
				report = d.Reports[i]
				found = true

				// 如果举报成立且是物品举报，下架物品
				if req.Status == "resolved" && d.Reports[i].Type == "item" {
					for j := range d.Items {
						if d.Items[j].ID == d.Reports[i].TargetID {
							d.Items[j].Status = "removed"
							break
						}
					}
				}

				// 记录审计日志
				actionDesc := "举报已驳回"
				if req.Status == "resolved" {
					actionDesc = "举报已处理，目标已下架"
				}
				d.AuditLogs = append(d.AuditLogs, model.AuditLog{
					ID:          service.GenID("log"),
					Action:      "report_" + req.Status,
					Description: actionDesc,
					OperatorID:  userId,
					TargetID:    id,
					CreatedAt:   model.TimeNow(),
				})
				break
			}
		}
	})

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"message": "举报记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "处理成功",
		"report":  report,
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

	// 验证交换是否完成、是否重复评价
	var exchangeFound, exchangeCompleted, duplicateRating bool
	service.WithRead(func(d *model.Database) {
		for _, ex := range d.Exchanges {
			if ex.ID == req.ExchangeID {
				exchangeFound = true
				if ex.Status == "completed" {
					exchangeCompleted = true
				}
				break
			}
		}
		for _, r := range d.Ratings {
			if r.ExchangeID == req.ExchangeID && r.FromUserID == userId {
				duplicateRating = true
				break
			}
		}
	})

	if !exchangeFound {
		c.JSON(http.StatusNotFound, gin.H{"message": "交换记录不存在"})
		return
	}
	if !exchangeCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"message": "只能评价已完成的交换"})
		return
	}
	if duplicateRating {
		c.JSON(http.StatusBadRequest, gin.H{"message": "你已经评价过此交换"})
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

// GetRatings 获取评价列表（按用户或按交换）
func GetRatings(c *gin.Context) {
	userId := c.Query("userId")
	exchangeId := c.Query("exchangeId")

	var ratings []model.Rating
	service.WithRead(func(d *model.Database) {
		for _, r := range d.Ratings {
			if userId != "" && r.ToUserID == userId {
				ratings = append(ratings, r)
			} else if exchangeId != "" && r.ExchangeID == exchangeId {
				ratings = append(ratings, r)
			} else if userId == "" && exchangeId == "" {
				ratings = append(ratings, r)
			}
		}
	})

	if ratings == nil {
		ratings = []model.Rating{}
	}

	// 如果按用户查询，额外返回平均分
	avgScore := 0.0
	if userId != "" && len(ratings) > 0 {
		total := 0
		for _, r := range ratings {
			total += r.Score
		}
		avgScore = float64(total) / float64(len(ratings))
	}

	result := gin.H{
		"message": "ok",
		"ratings": ratings,
	}
	if userId != "" {
		result["avgScore"] = avgScore
		result["ratingCount"] = len(ratings)
	}

	c.JSON(http.StatusOK, result)
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

	if notifications == nil {
		notifications = []model.Notification{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "ok",
		"notifications": notifications,
	})
}

// MarkNotificationRead 标记通知为已读
func MarkNotificationRead(c *gin.Context) {
	userId := c.GetString("userId")
	id := c.Param("id")

	var found bool
	service.WithWrite(func(d *model.Database) {
		for i := range d.Notifications {
			if d.Notifications[i].ID == id && d.Notifications[i].UserID == userId {
				d.Notifications[i].IsRead = true
				found = true
				break
			}
		}
	})

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"message": "通知不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "已标记为已读",
	})
}
