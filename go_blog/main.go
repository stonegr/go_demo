package main

import (
	"fmt"
	"go_blog/models"
	"go_blog/routes"
	"go_blog/utils"
	"log"
)

func main() {
	// 加载配置
	if err := utils.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	if err := models.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 设置路由
	r := routes.SetupRouter()

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%s",
		utils.AppConfig.Server.Host,
		utils.AppConfig.Server.Port)
	r.Run(serverAddr)
}
