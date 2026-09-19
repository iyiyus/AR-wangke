package handler

import (
	"strconv"
	"time"
	"wk-go/internal/database"
	"wk-go/internal/model"

	"github.com/gin-gonic/gin"
)

// GET /api/admin/huoyuan
func HuoYuanList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("current", c.DefaultQuery("page", "1")))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", c.DefaultQuery("pageSize", "15")))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 15
	}
	pt := c.Query("pt")

	db := database.DB.Model(&model.HuoYuan{})
	if pt != "" {
		db = db.Where("pt = ?", pt)
	}
	var total int64
	db.Count(&total)
	var list []model.HuoYuan
	db.Order("hid DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	PageResult(c, list, total, page, pageSize)
}

// POST /api/admin/huoyuan
func HuoYuanAdd(c *gin.Context) {
	var req model.HuoYuan
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	req.AddTime = time.Now().Format("2006-01-02 15:04:05")
	if req.Status == "" {
		req.Status = "1"
	}
	if err := database.DB.Create(&req).Error; err != nil {
		Fail(c, "添加失败")
		return
	}
	OK(c, req)
}

// PUT /api/admin/huoyuan/:id
func HuoYuanUpdate(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	database.DB.Model(&model.HuoYuan{}).Where("hid = ?", id).Updates(req)
	OK(c, nil)
}

// DELETE /api/admin/huoyuan/:id
func HuoYuanDelete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Where("hid = ?", id).Delete(&model.HuoYuan{})
	OK(c, nil)
}
