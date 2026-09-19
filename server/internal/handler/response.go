package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": data})
}

func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{"code": -1, "msg": msg})
}

func PageResult(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"list":     list,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

func parsePrice(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func atoiSafe(s string) (int, error) {
	return strconv.Atoi(s)
}

func parseInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func nowStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
