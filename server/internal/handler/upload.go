package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 上传目录（模拟对象存储）
var uploadDir string

func init() {
	// 优先使用环境变量指定的上传目录
	if envDir := os.Getenv("UPLOAD_DIR"); envDir != "" {
		uploadDir = envDir
	} else {
		// 使用相对于可执行文件或当前工作目录的上传目录
		// 尝试获取当前可执行文件所在目录
		exePath, err := os.Executable()
		if err == nil {
			// 检查是否在 go run 的临时目录中
			if !strings.Contains(exePath, "go-build") {
				uploadDir = filepath.Join(filepath.Dir(exePath), "uploads")
			} else {
				// go run 模式，使用 cwd 作为基准
				cwd, err := os.Getwd()
				if err == nil {
					uploadDir = filepath.Join(cwd, "uploads")
				} else {
					uploadDir = "uploads"
				}
			}
		} else {
			uploadDir = "uploads"
		}
	}
	os.MkdirAll(uploadDir, 0755)
}

// UploadImage 处理图片上传，模拟对象存储
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请选择要上传的图片"})
		return
	}

	// 限制文件大小 5MB
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "图片大小不能超过5MB"})
		return
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"message": "仅支持 jpg/png/gif/webp 格式"})
		return
	}

	// 生成唯一文件名（时间戳+原名）模拟对象存储的key
	timestamp := time.Now().UnixNano()
	objectKey := fmt.Sprintf("%d_%s", timestamp, sanitizeFilename(file.Filename))
	savePath := filepath.Join(uploadDir, objectKey)

	// 保存文件到本地存储（模拟对象存储的 PutObject）
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "上传失败"})
		return
	}

	// 返回访问URL（模拟对象存储的 CDN 地址）
	url := fmt.Sprintf("/uploads/%s", objectKey)
	c.JSON(http.StatusOK, gin.H{
		"url":       url,
		"objectKey": objectKey,
		"size":      file.Size,
		"filename":  file.Filename,
	})
}

// sanitizeFilename 清理文件名中的特殊字符
func sanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	return name
}
