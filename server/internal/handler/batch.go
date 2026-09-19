package handler

import (
	"wk-go/internal/database"
	"wk-go/internal/model"

	"github.com/gin-gonic/gin"
)

type batchDeleteReq struct {
	IDs []int64 `json:"ids" binding:"required"`
}

func batchDelete(c *gin.Context, m interface{}, idCol string) {
	var req batchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	if len(req.IDs) == 0 {
		Fail(c, "请选择要删除的项")
		return
	}
	res := database.DB.Where(idCol+" IN ?", req.IDs).Delete(m)
	OK(c, gin.H{"deleted": res.RowsAffected})
}

// POST /api/admin/class/batch-delete
func ClassBatchDelete(c *gin.Context) {
	batchDelete(c, &model.Class{}, "cid")
}

// POST /api/admin/fenlei/batch-delete
func FenLeiBatchDelete(c *gin.Context) {
	batchDelete(c, &model.FenLei{}, "id")
}

// POST /api/admin/dengji/batch-delete
func DengjiBatchDelete(c *gin.Context) {
	batchDelete(c, &model.Dengji{}, "id")
}

// POST /api/admin/huoyuan/batch-delete
func HuoYuanBatchDelete(c *gin.Context) {
	batchDelete(c, &model.HuoYuan{}, "hid")
}

// POST /api/admin/myprice/batch-delete
func MyPriceBatchDelete(c *gin.Context) {
	batchDelete(c, &model.MiJia{}, "mid")
}

// POST /api/admin/gongdan/batch-delete
func GongDanBatchDelete(c *gin.Context) {
	batchDelete(c, &model.GongDan{}, "gid")
}

// POST /api/admin/users/batch-delete
func UsersBatchDelete(c *gin.Context) {
	batchDelete(c, &model.User{}, "uid")
}

// POST /api/admin/order/batch-delete
func OrderBatchDelete(c *gin.Context) {
	batchDelete(c, &model.Order{}, "oid")
}

// POST /api/admin/log/batch-delete
func LogBatchDelete(c *gin.Context) {
	batchDelete(c, &model.Log{}, "id")
}

// POST /api/admin/pay/batch-delete
func PayBatchDelete(c *gin.Context) {
	batchDelete(c, &model.Pay{}, "oid")
}
