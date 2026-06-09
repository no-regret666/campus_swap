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
			{ID: "u1", Username: "202401001", Password: "123456", Nickname: "林可", Avatar: "", Role: "student", Campus: "主校区", CreditScore: 100},
			{ID: "u2", Username: "202401002", Password: "123456", Nickname: "周然", Avatar: "", Role: "student", Campus: "主校区", CreditScore: 95},
			{ID: "u3", Username: "202401003", Password: "123456", Nickname: "陈舟", Avatar: "", Role: "student", Campus: "南校区", CreditScore: 88},
			{ID: "admin", Username: "admin", Password: "admin123", Nickname: "管理员", Avatar: "", Role: "admin", Campus: "主校区", CreditScore: 100},
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
			// ===== u1 林可 发布的物品 =====
			{ID: "item5", Title: "机械键盘 Cherry红轴", Description: "Cherry MX Board 3.0S，87键，红轴手感好，用了一学期", Images: []string{}, Category: "digital", Campus: "主校区", Condition: "几乎全新", OwnerID: "u1", WantExchange: "想换Kindle或蓝牙耳机", Status: "available", Views: 18, CreatedAt: "2026-05-05 09:20:00"},
			{ID: "item6", Title: "瑜伽垫 加厚款", Description: "NBR材质10mm加厚，183×61cm，防滑纹理，用了一学期", Images: []string{}, Category: "sport", Campus: "主校区", Condition: "轻微使用", OwnerID: "u1", WantExchange: "哑铃或拉力带", Status: "available", Views: 6, CreatedAt: "2026-05-06 11:00:00"},
			{ID: "item7", Title: "水彩颜料套装", Description: "马利24色管装水彩，附带调色盘和画笔", Images: []string{}, Category: "art", Campus: "主校区", Condition: "全新", OwnerID: "u1", WantExchange: "手账素材或彩铅", Status: "available", Views: 9, CreatedAt: "2026-05-07 15:30:00"},
			{ID: "item8", Title: "保温杯 500ml", Description: "316不锈钢内胆，保温12小时，白色", Images: []string{}, Category: "life", Campus: "主校区", Condition: "几乎全新", OwnerID: "u1", WantExchange: "雨伞或收纳盒", Status: "available", Views: 4, CreatedAt: "2026-05-08 10:10:00"},
			{ID: "item9", Title: "Python编程辅导", Description: "Python基础/数据分析/爬虫，可线上或图书馆辅导，每次1.5小时", Images: []string{}, Category: "skill", Campus: "线上", Condition: "", OwnerID: "u1", WantExchange: "英语口语陪练或PS教学", Status: "available", Views: 22, CreatedAt: "2026-05-09 08:45:00"},
			{ID: "item10", Title: "小米蓝牙耳机", Description: "小米Air2 SE，降噪效果不错，充电盒轻微划痕", Images: []string{}, Category: "digital", Campus: "主校区", Condition: "轻微使用", OwnerID: "u1", WantExchange: "充电宝或数据线", Status: "available", Views: 15, CreatedAt: "2026-05-10 14:00:00"},
			// ===== u2 周然 发布的物品 =====
			{ID: "item11", Title: "Kindle Paperwhite 5", Description: "6.8寸300ppi，32G，带官方皮套，屏幕完好", Images: []string{}, Category: "digital", Campus: "主校区", Condition: "几乎全新", OwnerID: "u2", WantExchange: "iPad或手绘板", Status: "available", Views: 20, CreatedAt: "2026-05-02 16:00:00"},
			{ID: "item12", Title: "吉他入门教学", Description: "民谣吉他零基础教学，可上门或视频，送自编教材", Images: []string{}, Category: "skill", Campus: "主校区", Condition: "", OwnerID: "u2", WantExchange: "英语四六级辅导", Status: "available", Views: 11, CreatedAt: "2026-05-03 10:30:00"},
			{ID: "item13", Title: "运动背包 40L", Description: "Nike运动双肩包，多隔层设计，防水面料", Images: []string{}, Category: "sport", Campus: "主校区", Condition: "轻微使用", OwnerID: "u2", WantExchange: "登山杖或运动水壶", Status: "available", Views: 7, CreatedAt: "2026-05-05 13:20:00"},
			{ID: "item14", Title: "手账本套装", Description: "A5活页手账本+彩色笔+贴纸+胶带，未拆封", Images: []string{}, Category: "art", Campus: "主校区", Condition: "全新", OwnerID: "u2", WantExchange: "水彩或彩铅", Status: "available", Views: 10, CreatedAt: "2026-05-06 09:00:00"},
			{ID: "item15", Title: "电热水壶", Description: "美的1.5L不锈钢电水壶，自动断电，宿舍可用", Images: []string{}, Category: "life", Campus: "主校区", Condition: "轻微使用", OwnerID: "u2", WantExchange: "台灯或收纳架", Status: "available", Views: 5, CreatedAt: "2026-05-07 17:40:00"},
			{ID: "item16", Title: "Switch 健身环", Description: "含健身环和腿部绑带，卡带版，玩过两周", Images: []string{}, Category: "digital", Campus: "主校区", Condition: "几乎全新", OwnerID: "u2", WantExchange: "塞尔达卡带或其他Switch游戏", Status: "available", Views: 25, CreatedAt: "2026-05-08 11:30:00"},
			// ===== u3 陈舟 发布的物品 =====
			{ID: "item17", Title: "索尼WH-1000XM4", Description: "头戴式降噪耳机，黑色，续航30小时，耳罩轻微磨损", Images: []string{}, Category: "digital", Campus: "南校区", Condition: "轻微使用", OwnerID: "u3", WantExchange: "AirPods Pro或键盘", Status: "available", Views: 30, CreatedAt: "2026-05-01 20:00:00"},
			{ID: "item18", Title: "画架 折叠款", Description: "铝合金折叠画架，高度可调，适合油画/水彩", Images: []string{}, Category: "art", Campus: "南校区", Condition: "几乎全新", OwnerID: "u3", WantExchange: "颜料或画布", Status: "available", Views: 6, CreatedAt: "2026-05-04 08:30:00"},
			{ID: "item19", Title: "PS修图教学", Description: "Photoshop从入门到进阶，含人像精修/海报设计", Images: []string{}, Category: "skill", Campus: "线上", Condition: "", OwnerID: "u3", WantExchange: "编程辅导或摄影教学", Status: "available", Views: 14, CreatedAt: "2026-05-05 19:00:00"},
			{ID: "item20", Title: "跑鞋 Nike Pegasus 40", Description: "42码，黑色，跑了约200km，大底状态好", Images: []string{}, Category: "sport", Campus: "南校区", Condition: "轻微使用", OwnerID: "u3", WantExchange: "羽毛球拍或网球拍", Status: "available", Views: 9, CreatedAt: "2026-05-06 14:15:00"},
			{ID: "item21", Title: "桌面收纳架", Description: "竹木三层层架，35×25cm，宿舍桌面整理利器", Images: []string{}, Category: "life", Campus: "南校区", Condition: "几乎全新", OwnerID: "u3", WantExchange: "台灯或挂钟", Status: "available", Views: 3, CreatedAt: "2026-05-07 10:00:00"},
			{ID: "item22", Title: "小米充电宝 20000mAh", Description: "双USB输出，支持18W快充，白色", Images: []string{}, Category: "digital", Campus: "南校区", Condition: "几乎全新", OwnerID: "u3", WantExchange: "蓝牙耳机或数据线", Status: "available", Views: 11, CreatedAt: "2026-05-08 16:45:00"},
			// ===== 额外物品：混合分类和校区 =====
			{ID: "item23", Title: "飞利浦电动牙刷", Description: "HX6730声波电动牙刷，含两个刷头", Images: []string{}, Category: "life", Campus: "主校区", Condition: "轻微使用", OwnerID: "u1", WantExchange: "收纳用品或台灯", Status: "available", Views: 7, CreatedAt: "2026-05-11 09:00:00"},
			{ID: "item24", Title: "Wacom手绘板", Description: "CTL-472，适合入门板绘，附带压感笔", Images: []string{}, Category: "digital", Campus: "主校区", Condition: "几乎全新", OwnerID: "u2", WantExchange: "Kindle或iPad配件", Status: "available", Views: 16, CreatedAt: "2026-05-11 10:30:00"},
			{ID: "item25", Title: "羽毛球拍 尤尼克斯", Description: "ARC-5i，4U/G5，含球拍袋，线尚好", Images: []string{}, Category: "sport", Campus: "南校区", Condition: "轻微使用", OwnerID: "u3", WantExchange: "跑鞋或瑜伽垫", Status: "available", Views: 8, CreatedAt: "2026-05-11 11:00:00"},
			{ID: "item26", Title: "彩铅48色套装", Description: "辉柏嘉水溶性彩铅，铁盒装，未使用过", Images: []string{}, Category: "art", Campus: "主校区", Condition: "全新", OwnerID: "u1", WantExchange: "手账素材或水彩", Status: "available", Views: 12, CreatedAt: "2026-05-12 08:20:00"},
			{ID: "item27", Title: "摄影入门教学", Description: "单反/微单相机基础+构图+后期，可周末约拍", Images: []string{}, Category: "skill", Campus: "南校区", Condition: "", OwnerID: "u3", WantExchange: "PS教学或编程辅导", Status: "available", Views: 19, CreatedAt: "2026-05-12 13:00:00"},
			{ID: "item28", Title: "折叠书桌", Description: "80×40cm白色折叠桌，适合宿舍，承重50kg", Images: []string{}, Category: "life", Campus: "主校区", Condition: "几乎全新", OwnerID: "u2", WantExchange: "收纳架或椅子", Status: "available", Views: 5, CreatedAt: "2026-05-12 15:30:00"},
			{ID: "item29", Title: "罗技G304鼠标", Description: "无线游戏鼠标，LIGHTSPEED技术，5号电池供电", Images: []string{}, Category: "digital", Campus: "南校区", Condition: "几乎全新", OwnerID: "u3", WantExchange: "机械键盘或耳机", Status: "available", Views: 13, CreatedAt: "2026-05-13 09:45:00"},
			{ID: "item30", Title: "英语四级辅导", Description: "系统备考+真题讲解+作文批改，线上为主", Images: []string{}, Category: "skill", Campus: "线上", Condition: "", OwnerID: "u2", WantExchange: "吉他教学或绘画教学", Status: "available", Views: 17, CreatedAt: "2026-05-13 14:00:00"},
			{ID: "item31", Title: "跳绳 钢丝绳", Description: "竞速钢丝跳绳，轴承顺滑，适合减脂和体能训练", Images: []string{}, Category: "sport", Campus: "主校区", Condition: "全新", OwnerID: "u1", WantExchange: "瑜伽垫或拉力器", Status: "available", Views: 4, CreatedAt: "2026-05-14 10:00:00"},
			{ID: "item32", Title: "复古胶片相机", Description: "柯尼卡C35，旁轴胶片机，40mm f2.8镜头，测光正常", Images: []string{}, Category: "digital", Campus: "南校区", Condition: "轻微使用", OwnerID: "u3", WantExchange: "拍立得或单反镜头", Status: "available", Views: 21, CreatedAt: "2026-05-14 16:00:00"},
			{ID: "item33", Title: "加湿器 300ml", Description: "超声波雾化加湿器，静音设计，宿舍冬季必备", Images: []string{}, Category: "life", Campus: "南校区", Condition: "几乎全新", OwnerID: "u2", WantExchange: "电热水壶或台灯", Status: "available", Views: 6, CreatedAt: "2026-05-15 09:30:00"},
			{ID: "item34", Title: "口琴 C调十孔", Description: "和来Special 20，C调，适合初学蓝调口琴", Images: []string{}, Category: "art", Campus: "南校区", Condition: "几乎全新", OwnerID: "u3", WantExchange: "吉他配件或乐谱", Status: "available", Views: 5, CreatedAt: "2026-05-15 11:00:00"},
		},
		Exchanges: []model.Exchange{
			{
				ID: "ex1", ItemID: "item1", RequesterID: "u2", OwnerID: "u1",
				OfferDesc: "我可以用尤克里里来换", Message: "你好，我对你的iPad很感兴趣",
				Status: "pending", CreatedAt: "2026-05-10 08:00:00", UpdatedAt: "2026-05-10 08:00:00",
			},
			{
				ID: "ex2", ItemID: "item17", RequesterID: "u1", OwnerID: "u3",
				OfferDesc: "我可以用机械键盘来换", Message: "索尼耳机还在吗？",
				Status: "accepted", MeetTime: "周六 14:00", MeetPlace: "图书馆门口",
				CreatedAt: "2026-05-08 10:00:00", UpdatedAt: "2026-05-09 09:00:00",
			},
			{
				ID: "ex3", ItemID: "item11", RequesterID: "u3", OwnerID: "u2",
				OfferDesc: "我可以用索尼耳机来换", Message: "想换Kindle很久了",
				Status: "completed", CreatedAt: "2026-05-03 10:00:00", UpdatedAt: "2026-05-05 16:00:00",
			},
			{
				ID: "ex4", ItemID: "item9", RequesterID: "u3", OwnerID: "u1",
				OfferDesc: "我可以用摄影教学来换", Message: "我想学Python",
				Status: "completed", CreatedAt: "2026-05-09 10:00:00", UpdatedAt: "2026-05-12 18:00:00",
			},
			{
				ID: "ex5", ItemID: "item16", RequesterID: "u1", OwnerID: "u2",
				OfferDesc: "我可以用蓝牙耳机来换", Message: "健身环看着很有趣",
				Status: "rejected", CreatedAt: "2026-05-08 12:00:00", UpdatedAt: "2026-05-08 15:00:00",
			},
		},
		Messages: []model.Message{
			{ID: "msg1", ExchangeID: "ex1", SenderID: "u2", Content: "你好，iPad还在吗？我想用尤克里里交换", CreatedAt: "2026-05-10 08:00:00"},
			{ID: "msg2", ExchangeID: "ex1", SenderID: "u1", Content: "在的，你的尤克里里是什么牌子？", CreatedAt: "2026-05-10 08:30:00"},
			{ID: "msg3", ExchangeID: "ex2", SenderID: "u1", Content: "索尼耳机还在吗？我想用机械键盘换", CreatedAt: "2026-05-08 10:00:00"},
			{ID: "msg4", ExchangeID: "ex2", SenderID: "u3", Content: "在的，红轴还是青轴？", CreatedAt: "2026-05-08 10:15:00"},
			{ID: "msg5", ExchangeID: "ex2", SenderID: "u1", Content: "红轴，手感很舒服。周六图书馆门口见面？", CreatedAt: "2026-05-08 10:20:00"},
			{ID: "msg6", ExchangeID: "ex2", SenderID: "u3", Content: "好的，周六14点见", CreatedAt: "2026-05-08 10:25:00"},
		},
		Favorites: []model.Favorite{
			{ID: "fav1", UserID: "u1", ItemID: "item11", CreatedAt: "2026-05-02 17:00:00"},
			{ID: "fav2", UserID: "u1", ItemID: "item17", CreatedAt: "2026-05-01 21:00:00"},
			{ID: "fav3", UserID: "u1", ItemID: "item16", CreatedAt: "2026-05-08 12:00:00"},
			{ID: "fav4", UserID: "u2", ItemID: "item5", CreatedAt: "2026-05-05 10:00:00"},
			{ID: "fav5", UserID: "u2", ItemID: "item24", CreatedAt: "2026-05-11 11:00:00"},
			{ID: "fav6", UserID: "u3", ItemID: "item10", CreatedAt: "2026-05-10 15:00:00"},
			{ID: "fav7", UserID: "u3", ItemID: "item26", CreatedAt: "2026-05-12 09:00:00"},
		},
		Reports: []model.Report{},
		Ratings: []model.Rating{
			{ID: "rat1", ExchangeID: "ex3", FromUserID: "u3", ToUserID: "u2", Score: 5, Comment: "Kindle状态很好，交易很顺利", CreatedAt: "2026-05-05 16:30:00"},
			{ID: "rat2", ExchangeID: "ex3", FromUserID: "u2", ToUserID: "u3", Score: 4, Comment: "耳机降噪效果不错，就是有点磨损", CreatedAt: "2026-05-05 16:35:00"},
			{ID: "rat3", ExchangeID: "ex4", FromUserID: "u3", ToUserID: "u1", Score: 5, Comment: "Python讲得很清晰，很有耐心", CreatedAt: "2026-05-12 18:30:00"},
			{ID: "rat4", ExchangeID: "ex4", FromUserID: "u1", ToUserID: "u3", Score: 4, Comment: "摄影教学实用，构图部分讲得好", CreatedAt: "2026-05-12 18:35:00"},
		},
		BrowseRecords: []model.BrowseRecord{
			{ID: "br1", UserID: "u1", ItemID: "item2", Category: "art", CreatedAt: "2026-05-05 10:00:00"},
			{ID: "br2", UserID: "u1", ItemID: "item3", Category: "sport", CreatedAt: "2026-05-06 14:00:00"},
			{ID: "br3", UserID: "u1", ItemID: "item17", Category: "digital", CreatedAt: "2026-05-01 20:30:00"},
			{ID: "br4", UserID: "u1", ItemID: "item11", Category: "digital", CreatedAt: "2026-05-02 17:10:00"},
			{ID: "br5", UserID: "u1", ItemID: "item16", Category: "digital", CreatedAt: "2026-05-08 12:05:00"},
			{ID: "br6", UserID: "u1", ItemID: "item24", Category: "digital", CreatedAt: "2026-05-11 10:45:00"},
			{ID: "br7", UserID: "u1", ItemID: "item12", Category: "skill", CreatedAt: "2026-05-03 11:00:00"},
			{ID: "br8", UserID: "u2", ItemID: "item4", Category: "life", CreatedAt: "2026-05-07 09:00:00"},
			{ID: "br9", UserID: "u2", ItemID: "item1", Category: "digital", CreatedAt: "2026-05-08 11:00:00"},
			{ID: "br10", UserID: "u2", ItemID: "item5", Category: "digital", CreatedAt: "2026-05-05 09:30:00"},
			{ID: "br11", UserID: "u2", ItemID: "item24", Category: "digital", CreatedAt: "2026-05-11 11:00:00"},
			{ID: "br12", UserID: "u2", ItemID: "item9", Category: "skill", CreatedAt: "2026-05-09 09:00:00"},
			{ID: "br13", UserID: "u3", ItemID: "item1", Category: "digital", CreatedAt: "2026-05-09 16:00:00"},
			{ID: "br14", UserID: "u3", ItemID: "item5", Category: "digital", CreatedAt: "2026-05-05 09:30:00"},
			{ID: "br15", UserID: "u3", ItemID: "item10", Category: "digital", CreatedAt: "2026-05-10 14:10:00"},
			{ID: "br16", UserID: "u3", ItemID: "item26", Category: "art", CreatedAt: "2026-05-12 08:30:00"},
			{ID: "br17", UserID: "u3", ItemID: "item19", Category: "skill", CreatedAt: "2026-05-05 19:10:00"},
		},
	}

	return saveDBUnsafe()
}

// DBFileExists 检查数据文件是否存在
func DBFileExists() bool {
	_, err := os.Stat(config.DataFile)
	return err == nil
}
