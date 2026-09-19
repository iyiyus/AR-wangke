package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// POST /api/upload 上传文件，返回 {code:0, url: "/uploads/xxx.png"}
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "请选择文件"})
		return
	}
	// 限制 5MB
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "文件不能超过5MB"})
		return
	}
	ext := filepath.Ext(file.Filename)
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dir := "./uploads"
	os.MkdirAll(dir, 0755)
	dst := filepath.Join(dir, name)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": -1, "msg": "保存失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "上传成功", "data": gin.H{"url": "/uploads/" + name}})
}
