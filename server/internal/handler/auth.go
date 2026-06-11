package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"campus-swap-server/internal/config"
	"campus-swap-server/internal/model"
	"campus-swap-server/internal/service"
)

// SafeUser 返回不含密码的用户信息（含信用等级）
func SafeUser(u model.User) gin.H {
	level, badge := creditLevel(u.CreditScore)
	return gin.H{
		"id":          u.ID,
		"username":    u.Username,
		"nickname":    u.Nickname,
		"avatar":      u.Avatar,
		"role":        u.Role,
		"campus":      u.Campus,
		"creditScore": u.CreditScore,
		"creditLevel": level,
		"creditBadge": badge,
	}
}

// creditLevel 根据信用分计算等级和徽章
func creditLevel(score int) (string, string) {
	switch {
	case score >= 150:
		return "信赖", "🏆"
	case score >= 120:
		return "优秀", "⭐"
	case score >= 90:
		return "活跃", "🌟"
	case score >= 60:
		return "普通", "👤"
	default:
		return "警示", "⚠️"
	}
}

// Login 用户登录
func Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请输入用户名和密码"})
		return
	}

	var found *model.User
	service.WithRead(func(d *model.Database) {
		for i := range d.Users {
			if d.Users[i].Username == req.Username && d.Users[i].Password == req.Password {
				found = &d.Users[i]
				break
			}
		}
	})

	if found == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
		return
	}

	// 签发 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": found.ID,
		"exp":    time.Now().Add(time.Duration(config.JWTExpireHour) * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   tokenStr,
		"user":    SafeUser(*found),
	})
}

// Register 用户注册
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname" binding:"required"`
		Campus   string `json:"campus"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请填写完整注册信息"})
		return
	}

	var exists bool
	service.WithRead(func(d *model.Database) {
		for _, u := range d.Users {
			if u.Username == req.Username {
				exists = true
				break
			}
		}
	})

	if exists {
		c.JSON(http.StatusConflict, gin.H{"message": "用户名已存在"})
		return
	}

	newUser := model.User{
		ID:          service.GenID("u"),
		Username:    req.Username,
		Password:    req.Password,
		Nickname:    req.Nickname,
		Avatar:      "",
		Role:        "student",
		Campus:      req.Campus,
		CreditScore: 100,
	}

	service.WithWrite(func(d *model.Database) {
		d.Users = append(d.Users, newUser)
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"user":    SafeUser(newUser),
	})
}

// GetMe 获取当前登录用户信息
func GetMe(c *gin.Context) {
	userId := c.GetString("userId")

	var found *model.User
	service.WithRead(func(d *model.Database) {
		for i := range d.Users {
			if d.Users[i].ID == userId {
				found = &d.Users[i]
				break
			}
		}
	})

	if found == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"user":    SafeUser(*found),
	})
}
