package handler

import (
	"time"
	"wk-go/internal/database"
	"wk-go/internal/model"

	"github.com/gin-gonic/gin"
)

// GET /api/gongdan/list
func GongDanList(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)
	page, _ := atoiSafe(c.DefaultQuery("current", c.DefaultQuery("page", "1")))
	pageSize, _ := atoiSafe(c.DefaultQuery("size", c.DefaultQuery("pageSize", "15")))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 15
	}

	db := database.DB.Model(&model.GongDan{})
	if uidInt != 1 {
		db = db.Where("uid = ?", uidInt)
	}
	state := c.Query("state")
	if state != "" {
		db = db.Where("state = ?", state)
	}
	var total int64
	db.Count(&total)
	var list []model.GongDan
	db.Order("gid DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	PageResult(c, list, total, page, pageSize)
}

// POST /api/gongdan/add
func GongDanAdd(c *gin.Context) {
	uid, _ := c.Get("uid")
	var req struct {
		OID     int64  `json:"oid"`
		Region  string `json:"region" binding:"required"`
		Title   string `json:"title"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	gd := model.GongDan{
		UID:     uid.(int64),
		OID:     req.OID,
		Region:  req.Region,
		Title:   req.Title,
		Content: req.Content,
		State:   "待回复",
		AddTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	database.DB.Create(&gd)
	OK(c, nil)
}

// POST /api/gongdan/toanswer  用户追答（state 回到待回复）
func GongDanToAnswer(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)
	var req struct {
		GID     int64  `json:"gid" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	var gd model.GongDan
	if err := database.DB.First(&gd, req.GID).Error; err != nil {
		Fail(c, "工单不存在")
		return
	}
	if uidInt != 1 && gd.UID != uidInt {
		Fail(c, "无权操作")
		return
	}
	database.DB.Model(&gd).Updates(map[string]interface{}{
		"content": gd.Content + "\n[追答] " + req.Content,
		"state":   "待回复",
	})
	OK(c, nil)
}

// POST /api/gongdan/reply  管理员回复（state：已回复/已驳回/已关闭/不做处理）
func GongDanReply(c *gin.Context) {
	var req struct {
		GID    int64  `json:"gid" binding:"required"`
		Answer string `json:"answer" binding:"required"`
		State  string `json:"state"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	updates := map[string]interface{}{"answer": req.Answer}
	if req.State != "" {
		updates["state"] = req.State
	} else {
		updates["state"] = "已回复"
	}
	database.DB.Model(&model.GongDan{}).Where("gid = ?", req.GID).Updates(updates)
	OK(c, nil)
}

// POST /api/gongdan/state  工单状态流转（关闭/不做处理等）
func GongDanState(c *gin.Context) {
	var req struct {
		GID   int64  `json:"gid" binding:"required"`
		State string `json:"state" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	database.DB.Model(&model.GongDan{}).Where("gid = ?", req.GID).Update("state", req.State)
	OK(c, nil)
}
