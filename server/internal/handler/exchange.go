package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"campus-swap-server/internal/model"
	"campus-swap-server/internal/service"
)

// ExchangeDetail 返回给前端的交换详情（包含关联信息）
type ExchangeDetail struct {
	ID          string `json:"id"`
	ItemID      string `json:"itemId"`
	RequesterId string `json:"requesterId"`
	OwnerId     string `json:"ownerId"`
	OfferDesc   string `json:"offerDesc"`
	Message     string `json:"message"`
	MeetTime    string `json:"meetTime"`
	MeetPlace   string `json:"meetPlace"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	ItemTitle   string `json:"itemTitle"`
	ItemImage   string `json:"itemImage"`
	ItemCategory string `json:"itemCategory"`
	RequesterName string `json:"requesterName"`
	OwnerName string `json:"ownerName"`
	UnreadCount int `json:"unreadCount"` // 未读消息数
}

// GetExchanges 获取当前用户相关的交换（包含物品和用户详情）
func GetExchanges(c *gin.Context) {
	userId := c.GetString("userId")

	// 获取查询参数
	statusFilter := c.Query("status")
	searchKeyword := c.Query("search")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	var results []ExchangeDetail
	service.WithRead(func(d *model.Database) {
		// 建立物品和用户的索引
		itemMap := make(map[string]model.Item)
		for _, item := range d.Items {
			itemMap[item.ID] = item
		}
		userMap := make(map[string]model.User)
		for _, user := range d.Users {
			userMap[user.ID] = user
		}

		// 建立消息未读数索引
		unreadMap := make(map[string]int) // exchangeId -> unreadCount
		for _, msg := range d.Messages {
			// 找到消息对应的交换，然后统计发给当前用户的未读消息
			for _, ex := range d.Exchanges {
				if ex.ID == msg.ExchangeID {
					var otherUserId string
					if ex.OwnerID == userId {
						otherUserId = ex.RequesterID
					} else if ex.RequesterID == userId {
						otherUserId = ex.OwnerID
					}
					if msg.SenderID == otherUserId && !msg.IsRead {
						unreadMap[msg.ExchangeID]++
					}
					break
				}
			}
		}

		for _, ex := range d.Exchanges {
			if ex.RequesterID == userId || ex.OwnerID == userId {
				// 状态筛选
				if statusFilter != "" && ex.Status != statusFilter {
					continue
				}

				// 日期筛选
				if startDate != "" || endDate != "" {
					createdTime, err := time.Parse("2006-01-02 15:04:05", ex.CreatedAt)
					if err == nil {
						if startDate != "" {
							start, err := time.Parse("2006-01-02", startDate)
							if err == nil && createdTime.Before(start) {
								continue
							}
						}
						if endDate != "" {
							end, err := time.Parse("2006-01-02", endDate)
							if err == nil {
								end = end.Add(24*time.Hour - time.Second) // 包含结束日期当天
								if createdTime.After(end) {
									continue
								}
							}
						}
					}
				}

				detail := ExchangeDetail{
					ID:          ex.ID,
					ItemID:      ex.ItemID,
					RequesterId: ex.RequesterID,
					OwnerId:     ex.OwnerID,
					OfferDesc:   ex.OfferDesc,
					Message:     ex.Message,
					MeetTime:    ex.MeetTime,
					MeetPlace:   ex.MeetPlace,
					Status:      ex.Status,
					CreatedAt:   ex.CreatedAt,
					UpdatedAt:   ex.UpdatedAt,
					UnreadCount: unreadMap[ex.ID],
				}
				// 填充物品信息
				if item, ok := itemMap[ex.ItemID]; ok {
					detail.ItemTitle = item.Title
					detail.ItemCategory = item.Category
					if len(item.Images) > 0 {
						detail.ItemImage = item.Images[0]
					}
					// 搜索筛选
					if searchKeyword != "" {
						found := false
						keywords := []string{item.Title, item.Description, ex.OfferDesc, ex.Message}
						for _, kw := range keywords {
							if containsIgnoreCase(kw, searchKeyword) {
								found = true
								break
							}
						}
						if !found {
							continue
						}
					}
				}
				// 填充用户名
				if user, ok := userMap[ex.RequesterID]; ok {
					detail.RequesterName = user.Nickname
				}
				if user, ok := userMap[ex.OwnerID]; ok {
					detail.OwnerName = user.Nickname
				}
				results = append(results, detail)
			}
		}
	})

	if results == nil {
		results = []ExchangeDetail{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "ok",
		"exchanges": results,
	})
}

// containsIgnoreCase 字符串包含检查（不区分大小写）
func containsIgnoreCase(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		result[i] = c
	}
	return string(result)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// CreateExchange 发起交换申请
func CreateExchange(c *gin.Context) {
	userId := c.GetString("userId")

	var req struct {
		ItemID    string `json:"itemId" binding:"required"`
		OfferDesc string `json:"offerDesc"`
		Message   string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请填写交换信息"})
		return
	}

	var ownerID string
	var itemTitle string
	var itemFound bool
	service.WithRead(func(d *model.Database) {
		for _, item := range d.Items {
			if item.ID == req.ItemID {
				itemFound = true
				ownerID = item.OwnerID
				itemTitle = item.Title
				break
			}
		}
	})

	if !itemFound {
		c.JSON(http.StatusNotFound, gin.H{"message": "物品不存在"})
		return
	}

	if ownerID == userId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "不能交换自己的物品"})
		return
	}

	newExchange := model.Exchange{
		ID:          service.GenID("ex"),
		ItemID:      req.ItemID,
		RequesterID: userId,
		OwnerID:     ownerID,
		OfferDesc:   req.OfferDesc,
		Message:     req.Message,
		Status:      "pending",
		CreatedAt:   model.TimeNow(),
		UpdatedAt:   model.TimeNow(),
	}

	var requesterName string
	service.WithRead(func(d *model.Database) {
		for _, user := range d.Users {
			if user.ID == userId {
				requesterName = user.Nickname
				break
			}
		}
	})

	service.WithWrite(func(d *model.Database) {
		d.Exchanges = append(d.Exchanges, newExchange)
	})

	// 创建通知告知物品所有者
	createNotification(ownerID, "exchange_request", "新的交换申请",
		requesterName+" 想用 '"+req.OfferDesc+"' 交换你的 '"+itemTitle+"'",
		newExchange.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":  "交换申请已发送",
		"exchange": newExchange,
	})
}

// UpdateExchange 更新交换状态
func UpdateExchange(c *gin.Context) {
	userId := c.GetString("userId")
	id := c.Param("id")

	var req struct {
		Status    string `json:"status"`
		MeetTime  string `json:"meetTime"`
		MeetPlace string `json:"meetPlace"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请填写信息"})
		return
	}

	// 验证状态值
	validStatuses := map[string]bool{"pending": true, "accepted": true, "rejected": true, "completed": true, "cancelled": true}
	if req.Status != "" && !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的状态值"})
		return
	}

	var found bool
	var currentExchange model.Exchange
	service.WithRead(func(d *model.Database) {
		for _, ex := range d.Exchanges {
			if ex.ID == id {
				currentExchange = ex
				found = true
				break
			}
		}
	})

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"message": "交换记录不存在"})
		return
	}

	// 权限检查：只有参与者可以操作
	if currentExchange.OwnerID != userId && currentExchange.RequesterID != userId {
		c.JSON(http.StatusForbidden, gin.H{"message": "无权操作此交换"})
		return
	}

	// 处理交换取消（发起方可以取消 pending 状态的交换）
	if req.Status == "cancelled" {
		if currentExchange.RequesterID != userId {
			c.JSON(http.StatusForbidden, gin.H{"message": "只有发起方可以取消交换"})
			return
		}
		if currentExchange.Status != "pending" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "只能取消待确认的交换"})
			return
		}
	}

	// 状态流转权限：只有 owner 可以 accepted/rejected/completed
	ownerOnlyStatuses := map[string]bool{"accepted": true, "rejected": true, "completed": true}
	if req.Status != "" && ownerOnlyStatuses[req.Status] && currentExchange.OwnerID != userId {
		c.JSON(http.StatusForbidden, gin.H{"message": "只有物品所有者可以执行此操作"})
		return
	}

	var otherUserId string
	service.WithRead(func(d *model.Database) {
		if currentExchange.OwnerID == userId {
			otherUserId = currentExchange.RequesterID
		} else {
			otherUserId = currentExchange.OwnerID
		}
	})

	service.WithWrite(func(d *model.Database) {
		for i := range d.Exchanges {
			if d.Exchanges[i].ID == id {
				if req.Status != "" {
					d.Exchanges[i].Status = req.Status
				}
				if req.MeetTime != "" {
					d.Exchanges[i].MeetTime = req.MeetTime
				}
				if req.MeetPlace != "" {
					d.Exchanges[i].MeetPlace = req.MeetPlace
				}
				d.Exchanges[i].UpdatedAt = model.TimeNow()

				// 如果交换完成，更新物品状态
				if req.Status == "completed" {
					for j := range d.Items {
						if d.Items[j].ID == d.Exchanges[i].ItemID {
							d.Items[j].Status = "exchanged"
							break
						}
					}
				}

				currentExchange = d.Exchanges[i]
				break
			}
		}
	})

	// 发送通知
	if req.Status != "" && req.Status != "cancelled" {
		var userName string
		service.WithRead(func(d *model.Database) {
			for _, user := range d.Users {
				if user.ID == userId {
					userName = user.Nickname
					break
				}
			}
		})

		statusText := map[string]string{
			"accepted":  "接受了",
			"rejected":  "拒绝了",
			"completed": "确认完成了",
		}

		if text, ok := statusText[req.Status]; ok {
			var itemTitle string
			service.WithRead(func(d *model.Database) {
				for _, item := range d.Items {
					if item.ID == currentExchange.ItemID {
						itemTitle = item.Title
						break
					}
				}
			})
			createNotification(otherUserId, "exchange_response",
				"交换状态更新",
				userName+" "+text+" 交换请求 '"+itemTitle+"'",
				currentExchange.ID)
		}
	}

	// 如果是取消，通知另一方
	if req.Status == "cancelled" {
		var userName string
		service.WithRead(func(d *model.Database) {
			for _, user := range d.Users {
				if user.ID == userId {
					userName = user.Nickname
					break
				}
			}
		})
		createNotification(otherUserId, "exchange_response",
			"交换已取消",
			userName+" 取消了交换请求",
			currentExchange.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "更新成功",
		"exchange": currentExchange,
	})
}

// GetOverdueExchanges 获取逾期未处理的交换（超过24小时的pending状态）
func GetOverdueExchanges(c *gin.Context) {
	userId := c.GetString("userId")

	type OverdueExchange struct {
		ExchangeDetail
		OverdueHours int `json:"overdueHours"`
	}

	var results []OverdueExchange
	service.WithRead(func(d *model.Database) {
		for _, ex := range d.Exchanges {
			if ex.Status == "pending" && ex.OwnerID == userId {
				// 计算逾期时间
				createdTime, err := time.Parse("2006-01-02 15:04:05", ex.CreatedAt)
				if err == nil {
					hours := int(time.Since(createdTime).Hours())
					if hours >= 24 {
						itemTitle := ""
						for _, item := range d.Items {
							if item.ID == ex.ItemID {
								itemTitle = item.Title
								break
							}
						}
						requesterName := ""
						for _, user := range d.Users {
							if user.ID == ex.RequesterID {
								requesterName = user.Nickname
								break
							}
						}
						results = append(results, OverdueExchange{
							ExchangeDetail: ExchangeDetail{
								ID:           ex.ID,
								ItemID:       ex.ItemID,
								RequesterId:  ex.RequesterID,
								OwnerId:      ex.OwnerID,
								OfferDesc:    ex.OfferDesc,
								Message:      ex.Message,
								Status:       ex.Status,
								CreatedAt:    ex.CreatedAt,
								UpdatedAt:    ex.UpdatedAt,
								ItemTitle:    itemTitle,
								RequesterName: requesterName,
							},
							OverdueHours: hours,
						})
					}
				}
			}
		}
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"overdue": results,
	})
}
