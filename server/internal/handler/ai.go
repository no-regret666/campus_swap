package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GenerateDescription AI智能生成物品描述
// 本地模板引擎实现，无需外部API即可运行
// 如需接入真实大模型，替换 generateByTemplate 为 HTTP 调用即可
func GenerateDescription(c *gin.Context) {
	var req struct {
		Title     string `json:"title"`
		Category  string `json:"category"`
		Condition string `json:"condition"`
		Campus    string `json:"campus"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}

	if req.Title == "" && req.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请至少提供标题或分类"})
		return
	}

	desc := generateByTemplate(req.Title, req.Category, req.Condition, req.Campus)

	c.JSON(http.StatusOK, gin.H{
		"message":     "ok",
		"description": desc,
	})
}

// generateByTemplate 基于模板和规则生成描述（模拟大模型输出）
func generateByTemplate(title, category, condition, campus string) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 开场白模板
	openers := []string{
		"这是一件保存完好的闲置物品，",
		"诚意出交换，",
		"校园闲置好物分享，",
		"物品状态良好，适合有需要的同学，",
		"自用闲置转出，",
	}

	// 成色描述
	conditionDesc := map[string]string{
		"全新":     "物品全新未拆封，吊牌/包装完整。",
		"九成新":    "仅使用过几次，几乎看不出使用痕迹，功能完好。",
		"八成新":    "正常使用过一段时间，有轻微使用痕迹，不影响功能。",
		"七成新":    "使用时间较长，外观有一定磨损，但核心功能正常。",
		"六成新及以下": "有明显使用痕迹和磨损，但基本功能仍然可用。",
	}

	// 分类相关描述
	categoryDesc := map[string]string{
		"数码电子": "电子产品性能正常，无维修记录，配件齐全。",
		"图书教材": "书页干净无涂画，内容完整，适合下一届同学使用。",
		"运动户外": "运动装备功能正常，定期保养，适合日常锻炼使用。",
		"生活用品": "生活好物，日常使用方便，品质有保障。",
		"技能服务": "可提供相关技能服务，有实际经验，沟通后可约时间。",
	}

	// 校区信息
	campusInfo := ""
	if campus != "" {
		campusInfo = fmt.Sprintf("目前在%s，支持面交验货。", campus)
	}

	// 结尾
	endings := []string{
		"感兴趣的同学欢迎申请交换，可以私信沟通细节！",
		"期待与你以物换物，循环利用减少浪费！",
		"有意交换的同学可以直接发起申请，看到会及时回复。",
		"欢迎交换，让闲置物品找到新主人！",
	}

	// 组装描述
	var parts []string

	// 开场
	parts = append(parts, openers[rng.Intn(len(openers))])

	// 标题相关
	if title != "" {
		parts = append(parts, fmt.Sprintf("【%s】", title))
	}

	// 成色
	if desc, ok := conditionDesc[condition]; ok {
		parts = append(parts, desc)
	}

	// 分类
	if desc, ok := categoryDesc[category]; ok {
		parts = append(parts, desc)
	}

	// 校区
	if campusInfo != "" {
		parts = append(parts, campusInfo)
	}

	// 结尾
	parts = append(parts, endings[rng.Intn(len(endings))])

	return strings.Join(parts, "")
}
