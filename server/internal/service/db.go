package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"campus-swap-server/internal/config"
	"campus-swap-server/internal/model"
)

var (
	db   *model.Database
	mu   sync.RWMutex
	rng  *rand.Rand
)

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// GenID 生成唯一ID
func GenID(prefix string) string {
	return fmt.Sprintf("%s_%d_%04d", prefix, time.Now().UnixMilli(), rng.Intn(10000))
}

// GetDB 获取数据库实例（读锁）
func GetDB() *model.Database {
	mu.RLock()
	defer mu.RUnlock()
	return db
}

// WithWrite 在写锁保护下操作数据库并保存
func WithWrite(fn func(d *model.Database)) {
	mu.Lock()
	defer mu.Unlock()
	fn(db)
	saveDBUnsafe()
}

// WithRead 在读锁保护下操作数据库
func WithRead(fn func(d *model.Database)) {
	mu.RLock()
	defer mu.RUnlock()
	fn(db)
}

// LoadDB 从 JSON 文件读取 Database
func LoadDB() error {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(config.DataFile)
	if err != nil {
		return err
	}

	var database model.Database
	if err := json.Unmarshal(data, &database); err != nil {
		return err
	}
	db = &database
	return nil
}

// SaveDB 保存到 JSON 文件
func SaveDB() error {
	mu.Lock()
	defer mu.Unlock()
	return saveDBUnsafe()
}

func saveDBUnsafe() error {
	dir := filepath.Dir(config.DataFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(config.DataFile, data, 0644)
}

// SeedDB 创建初始数据
func SeedDB() error {
	mu.Lock()
	defer mu.Unlock()

	db = &model.Database{
		Users: []model.User{
			{ID: "u1", Username: "202401001", Password: "123456", Nickname: "林可", Avatar: "", Role: "student", Campus: "主校区"},
			{ID: "u2", Username: "202401002", Password: "123456", Nickname: "周然", Avatar: "", Role: "student", Campus: "主校区"},
			{ID: "u3", Username: "202401003", Password: "123456", Nickname: "陈舟", Avatar: "", Role: "student", Campus: "南校区"},
			{ID: "admin", Username: "admin", Password: "admin123", Nickname: "管理员", Avatar: "", Role: "admin", Campus: "主校区"},
		},
		Categories: []model.Category{
			{ID: "digital", Name: "数码设备", Icon: "💻"},
			{ID: "sport", Name: "运动户外", Icon: "⚽"},
			{ID: "life", Name: "生活用品", Icon: "🏠"},
			{ID: "art", Name: "文创乐器", Icon: "🎨"},
			{ID: "skill", Name: "技能服务", Icon: "🛠"},
		},
		Items: []model.Item{
			{
				ID: "item1", Title: "iPad Air 5", Description: "2022款 64G WiFi版，带笔和键盘",
				Images: []string{}, Category: "digital", Campus: "主校区",
				Condition: "几乎全新", OwnerID: "u1", WantExchange: "想换一台Kindle或机械键盘",
				Status: "available", Views: 12, CreatedAt: "2026-05-01 10:00:00",
			},
			{
				ID: "item2", Title: "尤克里里", Description: "23寸桃花心木面单，附带琴包和调音器",
				Images: []string{}, Category: "art", Campus: "主校区",
				Condition: "轻微使用", OwnerID: "u2", WantExchange: "绘画工具或小音箱",
				Status: "available", Views: 8, CreatedAt: "2026-05-02 14:30:00",
			},
			{
				ID: "item3", Title: "哑铃套装 20kg", Description: "包含杆和多片重量片",
				Images: []string{}, Category: "sport", Campus: "南校区",
				Condition: "明显使用", OwnerID: "u3", WantExchange: "瑜伽垫或跳绳等运动装备",
				Status: "available", Views: 5, CreatedAt: "2026-05-03 09:15:00",
			},
			{
				ID: "item4", Title: "台灯 护眼款", Description: "LED护眼台灯，三档色温可调",
				Images: []string{}, Category: "life", Campus: "主校区",
				Condition: "全新", OwnerID: "u1", WantExchange: "任何生活用品都可以",
				Status: "available", Views: 3, CreatedAt: "2026-05-04 16:00:00",
			},
		},
		Exchanges: []model.Exchange{
			{
				ID: "ex1", ItemID: "item1", RequesterID: "u2", OwnerID: "u1",
				OfferDesc: "我可以用尤克里里来换", Message: "你好，我对你的iPad很感兴趣",
				Status: "pending", CreatedAt: "2026-05-10 08:00:00", UpdatedAt: "2026-05-10 08:00:00",
			},
		},
		Messages: []model.Message{
			{ID: "msg1", ExchangeID: "ex1", SenderID: "u2", Content: "你好，iPad还在吗？我想用尤克里里交换", CreatedAt: "2026-05-10 08:00:00"},
			{ID: "msg2", ExchangeID: "ex1", SenderID: "u1", Content: "在的，你的尤克里里是什么牌子？", CreatedAt: "2026-05-10 08:30:00"},
		},
		Favorites: []model.Favorite{},
		Reports:   []model.Report{},
		Ratings:   []model.Rating{},
	}

	return saveDBUnsafe()
}

// DBFileExists 检查数据文件是否存在
func DBFileExists() bool {
	_, err := os.Stat(config.DataFile)
	return err == nil
}
