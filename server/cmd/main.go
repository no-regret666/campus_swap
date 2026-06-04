package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"campus-swap-server/internal/config"
	"campus-swap-server/internal/handler"
	"campus-swap-server/internal/middleware"
	"campus-swap-server/internal/service"
)

func main() {
	// 初始化数据库
	if service.DBFileExists() {
		if err := service.LoadDB(); err != nil {
			log.Printf("加载数据文件失败，重新初始化: %v", err)
			if err := service.SeedDB(); err != nil {
				log.Fatalf("初始化数据失败: %v", err)
			}
		}
	} else {
		log.Println("数据文件不存在，创建初始数据...")
		if err := service.SeedDB(); err != nil {
			log.Fatalf("初始化数据失败: %v", err)
		}
	}

	// 设置 Gin
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.CorsMiddleware())

	// 路由分组
	api := r.Group("/api")

	// 认证路由
	auth := api.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)
		auth.GET("/me", middleware.AuthMiddleware(), handler.GetMe)
	}

	// 物品路由
	items := api.Group("/items")
	{
		items.GET("", handler.GetItems)
		items.GET("/:id", handler.GetItem)
		items.POST("", middleware.AuthMiddleware(), handler.CreateItem)
	}

	// 交换路由
	exchanges := api.Group("/exchanges")
	exchanges.Use(middleware.AuthMiddleware())
	{
		exchanges.GET("", handler.GetExchanges)
		exchanges.POST("", handler.CreateExchange)
		exchanges.PATCH("/:id", handler.UpdateExchange)
	}

	// 消息路由
	messages := api.Group("/messages")
	messages.Use(middleware.AuthMiddleware())
	{
		messages.GET("/:exchangeId", handler.GetMessages)
		messages.POST("", handler.SendMessage)
	}

	// 收藏路由
	favorites := api.Group("/favorites")
	favorites.Use(middleware.AuthMiddleware())
	{
		favorites.GET("", handler.GetFavorites)
		favorites.POST("", handler.AddFavorite)
		favorites.DELETE("/:itemId", handler.RemoveFavorite)
	}

	// 分类路由
	api.GET("/categories", handler.GetCategories)

	// 后台统计
	api.GET("/dashboard", handler.GetDashboard)

	// 推荐
	api.GET("/recommendations", middleware.AuthMiddleware(), handler.GetRecommendations)

	// 举报
	api.POST("/reports", middleware.AuthMiddleware(), handler.CreateReport)

	// 评价
	api.POST("/ratings", middleware.AuthMiddleware(), handler.CreateRating)

	// 通知
	api.GET("/notifications", middleware.AuthMiddleware(), handler.GetNotifications)

	// 启动服务
	fmt.Printf("服务启动在 %s\n", config.ServerPort)
	if err := r.Run(config.ServerPort); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
