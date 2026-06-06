package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"campus-swap-server/internal/model"
	"campus-swap-server/internal/service"
)

// GetItems 获取物品列表（支持筛选）
func GetItems(c *gin.Context) {
	keyword := c.Query("keyword")
	category := c.Query("category")
	campus := c.Query("campus")
	owner := c.Query("owner")

	var results []model.Item
	service.WithRead(func(d *model.Database) {
		for _, item := range d.Items {
			// 如果指定了owner，只返回该用户的物品（不限status）
			if owner != "" {
				if item.OwnerID != owner {
					continue
				}
			} else {
				// 广场只展示available的物品
				if item.Status != "available" {
					continue
				}
			}
			if category != "" && item.Category != category {
				continue
			}
			if campus != "" && item.Campus != campus {
				continue
			}
			if keyword != "" {
				kw := strings.ToLower(keyword)
				if !strings.Contains(strings.ToLower(item.Title), kw) &&
					!strings.Contains(strings.ToLower(item.Description), kw) {
					continue
				}
			}
			results = append(results, item)
		}
	})

	if results == nil {
		results = []model.Item{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"items":   results,
	})
}

// GetItem 获取物品详情
func GetItem(c *gin.Context) {
	id := c.Param("id")

	var found *model.Item
	var ownerName string
	service.WithWrite(func(d *model.Database) {
		for i := range d.Items {
			if d.Items[i].ID == id {
				d.Items[i].Views++
				found = &d.Items[i]
				break
			}
		}
		if found != nil {
			for _, u := range d.Users {
				if u.ID == found.OwnerID {
					ownerName = u.Nickname
					if ownerName == "" {
						ownerName = u.Username
					}
					break
				}
			}
		}
	})

	if found == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "物品不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "ok",
		"item":      *found,
		"ownerName": ownerName,
	})
}

// CreateItem 创建物品
func CreateItem(c *gin.Context) {
	userId := c.GetString("userId")

	var req struct {
		Title        string   `json:"title" binding:"required"`
		Description  string   `json:"description"`
		Images       []string `json:"images"`
		Category     string   `json:"category" binding:"required"`
		Campus       string   `json:"campus"`
		Condition    string   `json:"condition"`
		WantExchange string   `json:"wantExchange"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请填写物品信息"})
		return
	}

	images := req.Images
	if images == nil {
		images = []string{}
	}

	newItem := model.Item{
		ID:           service.GenID("item"),
		Title:        req.Title,
		Description:  req.Description,
		Images:       images,
		Category:     req.Category,
		Campus:       req.Campus,
		Condition:    req.Condition,
		OwnerID:      userId,
		WantExchange: req.WantExchange,
		Status:       "available",
		Views:        0,
		CreatedAt:    model.TimeNow(),
	}

	service.WithWrite(func(d *model.Database) {
		d.Items = append(d.Items, newItem)
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "发布成功",
		"item":    newItem,
	})
}
