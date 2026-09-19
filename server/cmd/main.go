package main

import (
	"fmt"
	"log"
	"wk-go/internal/config"
	"wk-go/internal/cron"
	"wk-go/internal/database"
	"wk-go/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.Init("config/config.yaml"); err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	if err := database.Init(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 启动定时任务
	cron.Start()

	gin.SetMode(config.Global.Server.Mode)
	r := gin.Default()
	r.Static("/uploads", "./uploads")
	router.Setup(r)

	addr := fmt.Sprintf(":%d", config.Global.Server.Port)
	log.Printf("服务启动于 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
