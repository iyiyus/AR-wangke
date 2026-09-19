package handler

import (
	"net/http"

	"wk-go/internal/database"

	"github.com/gin-gonic/gin"
)

// GET /api/site 公开站点配置（favicon/logo/sitename，无需认证）
func SiteInfo(c *gin.Context) {
	cfg := map[string]string{}
	rows := database.DB.Model(nil).Table("qingka_wangke_config").Select("v,k").Find(&[]map[string]string{})
	_ = rows
	// 直接查
	type kv struct {
		V string `gorm:"column:v"`
		K string `gorm:"column:k"`
	}
	var list []kv
	database.DB.Raw("SELECT v,k FROM qingka_wangke_config").Scan(&list)
	for _, r := range list {
		cfg[r.V] = r.K
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"favicon":  cfg["favicon"],
			"logo":     cfg["logo"],
			"sitename": cfg["sitename"],
		},
	})
}
