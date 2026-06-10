package handler

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"

	"campus-swap-server/internal/model"
	"campus-swap-server/internal/service"
)

// GetItems 获取物品列表（支持筛选、排序）
func GetItems(c *gin.Context) {
	keyword := c.Query("keyword")
	category := c.Query("category")
	campus := c.Query("campus")
	owner := c.Query("owner")
	sortBy := c.Query("sortBy") // views / newest，默认 newest

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

		// 排序
		sort.Slice(results, func(i, j int) bool {
			if sortBy == "views" {
				return results[j].Views < results[i].Views
			}
			// 默认按时间 newest 倒序
			return results[j].CreatedAt < results[i].CreatedAt
		})
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

	// 尝试获取当前用户ID（未登录则为空）
	userId, _ := c.Get("userId")
	userIdStr, _ := userId.(string)

	var found *model.Item
	var ownerName string
	var isFavorited bool
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
			// 登录用户记录浏览历史
			if userIdStr != "" {
				d.BrowseRecords = append(d.BrowseRecords, model.BrowseRecord{
					ID:        service.GenID("br"),
					UserID:    userIdStr,
					ItemID:    id,
					Category:  found.Category,
					CreatedAt: model.TimeNow(),
				})
				// 检查是否已收藏
				for _, fav := range d.Favorites {
					if fav.UserID == userIdStr && fav.ItemID == id {
						isFavorited = true
						break
					}
				}
			}
		}
	})

	if found == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "物品不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "ok",
		"item":        *found,
		"ownerName":   ownerName,
		"isFavorited": isFavorited,
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

// UpdateItem 修改自己发布的物品
func UpdateItem(c *gin.Context) {
	userId := c.GetString("userId")
	id := c.Param("id")

	var req struct {
		Title        string   `json:"title"`
		Description  string   `json:"description"`
		Images       []string `json:"images"`
		Category     string   `json:"category"`
		Campus       string   `json:"campus"`
		Condition    string   `json:"condition"`
		WantExchange string   `json:"wantExchange"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	var found bool
	var updatedItem model.Item
	service.WithWrite(func(d *model.Database) {
		for i := range d.Items {
			if d.Items[i].ID == id {
				if d.Items[i].OwnerID != userId {
					return
				}
				// 更新非空字段
				if req.Title != "" {
					d.Items[i].Title = req.Title
				}
				if req.Description != "" {
					d.Items[i].Description = req.Description
				}
				if req.Images != nil {
					d.Items[i].Images = req.Images
				}
				if req.Category != "" {
					d.Items[i].Category = req.Category
				}
				if req.Campus != "" {
					d.Items[i].Campus = req.Campus
				}
				if req.Condition != "" {
					d.Items[i].Condition = req.Condition
				}
				if req.WantExchange != "" {
					d.Items[i].WantExchange = req.WantExchange
				}
				updatedItem = d.Items[i]
				found = true
				break
			}
		}
	})

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"message": "物品不存在或无权修改"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "修改成功",
		"item":    updatedItem,
	})
}

// DeleteItem 删除自己发布的物品
func DeleteItem(c *gin.Context) {
	userId := c.GetString("userId")
	id := c.Param("id")

	var found bool
	service.WithWrite(func(d *model.Database) {
		for i := range d.Items {
			if d.Items[i].ID == id {
				if d.Items[i].OwnerID != userId {
					return
				}
				// 从列表中删除
				d.Items = append(d.Items[:i], d.Items[i+1:]...)
				found = true
				break
			}
		}
		if found {
			// 同时删除关联的收藏记录
			var newFavorites []model.Favorite
			for _, f := range d.Favorites {
				if f.ItemID != id {
					newFavorites = append(newFavorites, f)
				}
			}
			d.Favorites = newFavorites
		}
	})

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"message": "物品不存在或无权删除"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
