package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"campus-swap-server/internal/model"
	"campus-swap-server/internal/service"
)

// GetFavorites 获取当前用户的收藏列表（含物品详情）
func GetFavorites(c *gin.Context) {
	userId := c.GetString("userId")

	type FavoriteWithItem struct {
		model.Favorite
		Item *model.Item `json:"item"`
	}

	var results []FavoriteWithItem
	service.WithRead(func(d *model.Database) {
		// 建立物品ID索引
		itemMap := make(map[string]*model.Item)
		for i := range d.Items {
			itemMap[d.Items[i].ID] = &d.Items[i]
		}
		for _, fav := range d.Favorites {
			if fav.UserID == userId {
				fi := FavoriteWithItem{Favorite: fav}
				if item, ok := itemMap[fav.ItemID]; ok {
					fi.Item = item
				}
				results = append(results, fi)
			}
		}
	})

	if results == nil {
		results = []FavoriteWithItem{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "ok",
		"favorites": results,
	})
}

// AddFavorite 添加收藏
func AddFavorite(c *gin.Context) {
	userId := c.GetString("userId")

	var req struct {
		ItemID string `json:"itemId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请指定物品ID"})
		return
	}

	// 检查是否已收藏
	var already bool
	service.WithRead(func(d *model.Database) {
		for _, fav := range d.Favorites {
			if fav.UserID == userId && fav.ItemID == req.ItemID {
				already = true
				break
			}
		}
	})

	if already {
		c.JSON(http.StatusConflict, gin.H{"message": "已收藏该物品"})
		return
	}

	newFav := model.Favorite{
		ID:        service.GenID("fav"),
		UserID:    userId,
		ItemID:    req.ItemID,
		CreatedAt: model.TimeNow(),
	}

	service.WithWrite(func(d *model.Database) {
		d.Favorites = append(d.Favorites, newFav)
	})

	c.JSON(http.StatusOK, gin.H{
		"message":  "收藏成功",
		"favorite": newFav,
	})
}

// RemoveFavorite 取消收藏
func RemoveFavorite(c *gin.Context) {
	userId := c.GetString("userId")
	itemId := c.Param("itemId")

	var removed bool
	service.WithWrite(func(d *model.Database) {
		for i, fav := range d.Favorites {
			if fav.UserID == userId && fav.ItemID == itemId {
				d.Favorites = append(d.Favorites[:i], d.Favorites[i+1:]...)
				removed = true
				break
			}
		}
	})

	if !removed {
		c.JSON(http.StatusNotFound, gin.H{"message": "未找到收藏记录"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消收藏成功"})
}
