package handler

import (
	"time"
	"wk-go/internal/database"
	"wk-go/internal/model"

	"github.com/gin-gonic/gin"
)

// GET /api/myprice/list
func MyPriceList(c *gin.Context) {
	uid, _ := c.Get("uid")
	var list []model.MiJia
	database.DB.Where("uid = ?", uid).Find(&list)
	OK(c, list)
}

// POST /api/myprice/set  管理员设置密价
func MyPriceSet(c *gin.Context) {
	var req model.MiJia
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	req.AddTime = time.Now().Format("2006-01-02 15:04:05")

	var existing model.MiJia
	if err := database.DB.Where("uid = ? AND cid = ?", req.UID, req.CID).First(&existing).Error; err == nil {
		database.DB.Model(&existing).Updates(map[string]interface{}{
			"mode": req.Mode, "price": req.Price,
		})
	} else {
		database.DB.Create(&req)
	}
	OK(c, nil)
}

// DELETE /api/myprice/:id  管理员删除密价
func MyPriceDelete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&model.MiJia{}, id)
	OK(c, nil)
}
