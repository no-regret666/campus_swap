package handler

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"

	"campus-swap-server/internal/model"
	"campus-swap-server/internal/service"
)

// GetMessages 获取某交换的消息列表
func GetMessages(c *gin.Context) {
	exchangeId := c.Param("exchangeId")
	userId := c.GetString("userId")

	// 验证用户是否有权访问此交换的消息
	var isParticipant bool
	service.WithRead(func(d *model.Database) {
		for _, ex := range d.Exchanges {
			if ex.ID == exchangeId && (ex.OwnerID == userId || ex.RequesterID == userId) {
				isParticipant = true
				break
			}
		}
	})

	if !isParticipant {
		c.JSON(http.StatusForbidden, gin.H{"message": "无权查看此交换的消息"})
		return
	}

	var results []model.Message
	service.WithRead(func(d *model.Database) {
		for _, msg := range d.Messages {
			if msg.ExchangeID == exchangeId {
				results = append(results, msg)
			}
		}
	})

	if results == nil {
		results = []model.Message{}
	}

	// 按创建时间升序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt < results[j].CreatedAt
	})

	c.JSON(http.StatusOK, gin.H{
		"message":  "ok",
		"messages": results,
	})
}

// SendMessage 发送消息
func SendMessage(c *gin.Context) {
	userId := c.GetString("userId")

	var req struct {
		ExchangeID string   `json:"exchangeId" binding:"required"`
		Content    string   `json:"content"`
		Images     []string `json:"images"` // 支持多图片
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	// 验证：至少要有文字内容或图片
	if req.Content == "" && (req.Images == nil || len(req.Images) == 0) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请填写消息内容或添加图片"})
		return
	}

	// 验证用户是否有权向此交换发送消息
	var isParticipant bool
	service.WithRead(func(d *model.Database) {
		for _, ex := range d.Exchanges {
			if ex.ID == req.ExchangeID && (ex.OwnerID == userId || ex.RequesterID == userId) {
				isParticipant = true
				break
			}
		}
	})

	if !isParticipant {
		c.JSON(http.StatusForbidden, gin.H{"message": "无权向此交换发送消息"})
		return
	}

	newMsg := model.Message{
		ID:         service.GenID("msg"),
		ExchangeID: req.ExchangeID,
		SenderID:   userId,
		Content:    req.Content,
		Images:     req.Images,
		IsRead:     false,
		IsDeleted:  false,
		CreatedAt:  model.TimeNow(),
	}

	// 查找交换的对方用户ID，发送通知
	var recipientID string
	service.WithRead(func(d *model.Database) {
		for _, ex := range d.Exchanges {
			if ex.ID == req.ExchangeID {
				if ex.OwnerID == userId {
					recipientID = ex.RequesterID
				} else {
					recipientID = ex.OwnerID
				}
				break
			}
		}
	})

	service.WithWrite(func(d *model.Database) {
		d.Messages = append(d.Messages, newMsg)
		// 通知对方：收到新消息
		if recipientID != "" {
			d.Notifications = append(d.Notifications, model.Notification{
				ID:        service.GenID("ntf"),
				UserID:    recipientID,
				Type:      "new_message",
				Title:     "收到新消息",
				Content:   "你有一条新的交换沟通消息",
				RelatedID: req.ExchangeID,
				CreatedAt: model.TimeNow(),
			})
		}
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "发送成功",
		"data":    newMsg,
	})
}

// MarkMessagesRead 标记消息为已读
func MarkMessagesRead(c *gin.Context) {
	exchangeId := c.Param("exchangeId")
	userId := c.GetString("userId")

	// 验证用户是否有权访问此交换
	var isParticipant bool
	var otherUserId string
	service.WithRead(func(d *model.Database) {
		for _, ex := range d.Exchanges {
			if ex.ID == exchangeId {
				if ex.OwnerID == userId {
					isParticipant = true
					otherUserId = ex.RequesterID
				} else if ex.RequesterID == userId {
					isParticipant = true
					otherUserId = ex.OwnerID
				}
				break
			}
		}
	})

	if !isParticipant {
		c.JSON(http.StatusForbidden, gin.H{"message": "无权访问此交换的消息"})
		return
	}

	service.WithWrite(func(d *model.Database) {
		for i := range d.Messages {
			if d.Messages[i].ExchangeID == exchangeId && d.Messages[i].SenderID == otherUserId && !d.Messages[i].IsRead {
				d.Messages[i].IsRead = true
			}
		}
	})

	c.JSON(http.StatusOK, gin.H{"message": "已标记已读"})
}

// DeleteMessage 删除/撤回消息（只能撤回自己发送的，且在5分钟内）
func DeleteMessage(c *gin.Context) {
	userId := c.GetString("userId")
	messageId := c.Param("id")

	var found bool
	var message model.Message
	service.WithRead(func(d *model.Database) {
		for _, msg := range d.Messages {
			if msg.ID == messageId {
				message = msg
				found = true
				break
			}
		}
	})

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"message": "消息不存在"})
		return
	}

	// 只能删除自己发送的消息
	if message.SenderID != userId {
		c.JSON(http.StatusForbidden, gin.H{"message": "只能撤回自己的消息"})
		return
	}

	// 检查是否在5分钟内（300秒）
	createdTime, _ := time.Parse("2006-01-02 15:04:05", message.CreatedAt)
	if time.Since(createdTime).Seconds() > 300 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "消息已超过5分钟，无法撤回"})
		return
	}

	service.WithWrite(func(d *model.Database) {
		for i := range d.Messages {
			if d.Messages[i].ID == messageId {
				d.Messages[i].IsDeleted = true
				d.Messages[i].Content = "此消息已撤回"
				break
			}
		}
	})

	c.JSON(http.StatusOK, gin.H{"message": "消息已撤回"})
}

// GetUnreadCount 获取未读消息数量
func GetUnreadCount(c *gin.Context) {
	userId := c.GetString("userId")

	var unreadCount int
	service.WithRead(func(d *model.Database) {
		// 建立交换ID到用户ID的映射
		exchangeMap := make(map[string]string) // exchangeId -> otherUserId
		for _, ex := range d.Exchanges {
			if ex.OwnerID == userId {
				exchangeMap[ex.ID] = ex.RequesterID
			} else if ex.RequesterID == userId {
				exchangeMap[ex.ID] = ex.OwnerID
			}
		}

		// 统计来自其他用户且未读的消息
		for _, msg := range d.Messages {
			if otherId, ok := exchangeMap[msg.ExchangeID]; ok {
				if msg.SenderID == otherId && !msg.IsRead {
					unreadCount++
				}
			}
		}
	})

	c.JSON(http.StatusOK, gin.H{
		"unreadCount": unreadCount,
	})
}

// GetAllUnreadByExchange 获取每个交换的未读数
func GetAllUnreadByExchange(c *gin.Context) {
	userId := c.GetString("userId")

	var results map[string]int
	service.WithRead(func(d *model.Database) {
		results = make(map[string]int)

		// 建立交换ID到用户ID的映射
		exchangeMap := make(map[string]string) // exchangeId -> otherUserId
		for _, ex := range d.Exchanges {
			if ex.OwnerID == userId {
				exchangeMap[ex.ID] = ex.RequesterID
			} else if ex.RequesterID == userId {
				exchangeMap[ex.ID] = ex.OwnerID
			}
		}

		// 统计每个交换的未读消息数
		for _, msg := range d.Messages {
			if otherId, ok := exchangeMap[msg.ExchangeID]; ok {
				if msg.SenderID == otherId && !msg.IsRead {
					results[msg.ExchangeID]++
				}
			}
		}
	})

	c.JSON(http.StatusOK, gin.H{"unreadByExchange": results})
}
