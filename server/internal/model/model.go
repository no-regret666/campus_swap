package model

import "time"

// User 用户模型
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"` // student / admin
	Campus   string `json:"campus"`
}

// Category 分类
type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

// Item 物品
type Item struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Category    string   `json:"category"`
	Campus      string   `json:"campus"`
	Condition   string   `json:"condition"` // 全新 / 几乎全新 / 轻微使用 / 明显使用
	OwnerID     string   `json:"ownerId"`
	WantExchange string  `json:"wantExchange"` // 期望交换的物品描述
	Status      string   `json:"status"`       // available / exchanged / removed
	Views       int      `json:"views"`
	CreatedAt   string   `json:"createdAt"`
}

// Exchange 交换申请
type Exchange struct {
	ID          string `json:"id"`
	ItemID      string `json:"itemId"`
	RequesterID string `json:"requesterId"`
	OwnerID     string `json:"ownerId"`
	OfferDesc   string `json:"offerDesc"` // 申请人提供的交换物描述
	Message     string `json:"message"`
	Status      string `json:"status"` // pending / accepted / rejected / completed / cancelled
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// Message 站内消息
type Message struct {
	ID         string `json:"id"`
	ExchangeID string `json:"exchangeId"`
	SenderID   string `json:"senderId"`
	Content    string `json:"content"`
	CreatedAt  string `json:"createdAt"`
}

// Favorite 收藏
type Favorite struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	ItemID    string `json:"itemId"`
	CreatedAt string `json:"createdAt"`
}

// Report 举报
type Report struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	TargetID  string `json:"targetId"`
	Type      string `json:"type"` // item / user
	Reason    string `json:"reason"`
	Status    string `json:"status"` // pending / resolved
	CreatedAt string `json:"createdAt"`
}

// Rating 评价
type Rating struct {
	ID         string `json:"id"`
	ExchangeID string `json:"exchangeId"`
	FromUserID string `json:"fromUserId"`
	ToUserID   string `json:"toUserId"`
	Score      int    `json:"score"` // 1-5
	Comment    string `json:"comment"`
	CreatedAt  string `json:"createdAt"`
}

// Database JSON 文件存储的完整结构
type Database struct {
	Users      []User     `json:"users"`
	Categories []Category `json:"categories"`
	Items      []Item     `json:"items"`
	Exchanges  []Exchange `json:"exchanges"`
	Messages   []Message  `json:"messages"`
	Favorites  []Favorite `json:"favorites"`
	Reports    []Report   `json:"reports"`
	Ratings    []Rating   `json:"ratings"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// TimeNow 获取当前时间字符串
func TimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
