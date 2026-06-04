package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"campus-swap-server/internal/model"
	"campus-swap-server/internal/service"
)

// GetMessages 获取某交换的消息列表
func GetMessages(c *gin.Context) {
	exchangeId := c.Param("exchangeId")

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

	c.JSON(http.StatusOK, gin.H{
		"message":  "ok",
		"messages": results,
	})
}

// SendMessage 发送消息
func SendMessage(c *gin.Context) {
	userId := c.GetString("userId")

	var req struct {
		ExchangeID string `json:"exchangeId" binding:"required"`
		Content    string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请填写消息内容"})
		return
	}

	newMsg := model.Message{
		ID:         service.GenID("msg"),
		ExchangeID: req.ExchangeID,
		SenderID:   userId,
		Content:    req.Content,
		CreatedAt:  model.TimeNow(),
	}

	service.WithWrite(func(d *model.Database) {
		d.Messages = append(d.Messages, newMsg)
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "发送成功",
		"data":    newMsg,
	})
}
