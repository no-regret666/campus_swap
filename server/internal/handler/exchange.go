package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"campus-swap-server/internal/model"
	"campus-swap-server/internal/service"
)

// GetExchanges 获取当前用户相关的交换
func GetExchanges(c *gin.Context) {
	userId := c.GetString("userId")

	var results []model.Exchange
	service.WithRead(func(d *model.Database) {
		for _, ex := range d.Exchanges {
			if ex.RequesterID == userId || ex.OwnerID == userId {
				results = append(results, ex)
			}
		}
	})

	if results == nil {
		results = []model.Exchange{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "ok",
		"exchanges": results,
	})
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
	var itemFound bool
	service.WithRead(func(d *model.Database) {
		for _, item := range d.Items {
			if item.ID == req.ItemID {
				itemFound = true
				ownerID = item.OwnerID
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

	service.WithWrite(func(d *model.Database) {
		d.Exchanges = append(d.Exchanges, newExchange)
	})

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
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请填写状态"})
		return
	}

	// 验证状态值
	validStatuses := map[string]bool{"accepted": true, "rejected": true, "completed": true, "cancelled": true}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的状态值"})
		return
	}

	var found bool
	var exchange model.Exchange
	service.WithWrite(func(d *model.Database) {
		for i := range d.Exchanges {
			if d.Exchanges[i].ID == id {
				// 只有物品所有者或申请者可以更新
				if d.Exchanges[i].OwnerID != userId && d.Exchanges[i].RequesterID != userId {
					return
				}
				d.Exchanges[i].Status = req.Status
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

				exchange = d.Exchanges[i]
				found = true
				break
			}
		}
	})

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"message": "交换记录不存在或无权操作"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "状态更新成功",
		"exchange": exchange,
	})
}
