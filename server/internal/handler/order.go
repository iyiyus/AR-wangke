package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"wk-go/internal/checkorder"
	"wk-go/internal/database"
	"wk-go/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type orderQuery struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
	Current    int    `form:"current"`
	Size       int    `form:"size"`
	OID        string `form:"oid"`
	UID        string `form:"uid"`
	QQ         string `form:"qq"`
	CID        string `form:"cid"`
	StatusText string `form:"status_text"`
	Dock       string `form:"dock"`
}

// GET /api/order/list
func OrderList(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)

	var q orderQuery
	c.ShouldBindQuery(&q)
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 15
	}
	// 兼容 current/size
	if q.Current > 0 {
		q.Page = q.Current
	}
	if q.Size > 0 {
		q.PageSize = q.Size
	}

	db := database.DB.Model(&model.Order{})

	// 非管理员只看自己的
	if uidInt != 1 {
		db = db.Where("uid = ?", uidInt)
	} else if q.UID != "" {
		db = db.Where("uid = ?", q.UID)
	}

	if q.OID != "" {
		db = db.Where("oid = ?", q.OID)
	}
	if q.QQ != "" {
		db = db.Where("user LIKE ? OR school LIKE ?", "%"+q.QQ+"%", "%"+q.QQ+"%")
	}
	if q.CID != "" {
		db = db.Where("cid = ?", q.CID)
	}
	if q.StatusText != "" {
		db = db.Where("status = ?", q.StatusText)
	}
	if q.Dock != "" {
		db = db.Where("dockstatus = ?", q.Dock)
	}

	var total int64
	db.Count(&total)

	var orders []model.Order
	db.Order("oid DESC").Offset((q.Page - 1) * q.PageSize).Limit(q.PageSize).Find(&orders)

	PageResult(c, orders, total, q.Page, q.PageSize)
}

// POST /api/order/add
func OrderAdd(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)

	var req struct {
		CID    int64  `json:"cid" binding:"required"`
		School string `json:"school" binding:"required"`
		User   string `json:"user" binding:"required"`
		Pass   string `json:"pass" binding:"required"`
		KCID   string `json:"kcid"`
		KCName string `json:"kcname" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数不能为空")
		return
	}

	var user model.User
	database.DB.First(&user, uidInt)

	var class model.Class
	if err := database.DB.First(&class, req.CID).Error; err != nil {
		Fail(c, "平台不存在")
		return
	}
	if class.Status == 0 {
		Fail(c, "该平台已下架")
		return
	}

	// 计算单价
	danjia := calcPrice(class, user)
	if danjia <= 0 {
		Fail(c, "价格异常，请联系管理员")
		return
	}

	kcNames := strings.Split(req.KCName, ",")
	kcIDs := strings.Split(req.KCID, ",")

	if user.Money < danjia*float64(len(kcNames)) {
		Fail(c, "余额不足")
		return
	}

	ip := c.ClientIP()
	now := time.Now().Format("2006-01-02 15:04:05")

	for i, name := range kcNames {
		kcid := ""
		if i < len(kcIDs) {
			kcid = kcIDs[i]
		}
		dockStatus := "0"
		if class.Docking == "0" {
			dockStatus = "99"
		}
		// 检查重复
		var dup int64
		database.DB.Model(&model.Order{}).Where(
			"ptname=? AND school=? AND user=? AND pass=? AND kcname=?",
			class.Name, req.School, req.User, req.Pass, name,
		).Count(&dup)
		if dup > 0 {
			dockStatus = "3"
		}

		order := model.Order{
			UID:         uidInt,
			CID:         req.CID,
			HID:         func() int64 { v, _ := strconv.ParseInt(class.Docking, 10, 64); return v }(),
			PTName:      class.Name,
			School:      req.School,
			UserAccount: req.User,
			Pass:        req.Pass,
			KCID:        kcid,
			KCName:      name,
			Fees:        fmt.Sprintf("%.2f", danjia),
			Noun:        class.Noun,
			AddTime:     now,
			IP:          ip,
			DockStatus:  dockStatus,
			Status:      "待处理",
		}
		if err := database.DB.Create(&order).Error; err == nil {
			database.DB.Model(&model.User{}).Where("uid = ?", uidInt).
				UpdateColumn("money", gorm.Expr("money - ?", danjia))
			writeLog(uidInt, "下单", fmt.Sprintf("%s %s %s 扣除%.2f元", req.User, req.Pass, name, danjia), -danjia, ip)
		}
	}
	OK(c, nil)
}

// POST /api/order/cancel
func OrderCancel(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)
	var req struct {
		OID int64 `json:"oid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	var order model.Order
	if err := database.DB.First(&order, req.OID).Error; err != nil {
		Fail(c, "订单不存在")
		return
	}
	if uidInt != 1 && order.UID != uidInt {
		Fail(c, "无权操作")
		return
	}
	database.DB.Model(&order).Updates(map[string]interface{}{
		"status": "已取消", "dockstatus": "4",
	})
	OK(c, nil)
}

// POST /api/order/restart  补刷
func OrderRestart(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)
	var req struct {
		OID int64 `json:"oid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	var order model.Order
	if err := database.DB.First(&order, req.OID).Error; err != nil {
		Fail(c, "订单不存在")
		return
	}
	if uidInt != 1 && order.UID != uidInt {
		Fail(c, "无权操作")
		return
	}
	bsNum, _ := strconv.Atoi(order.BSNum)
	if bsNum >= 5 {
		Fail(c, "补刷已超过5次")
		return
	}
	// 自营订单仅排队
	if order.DockStatus == "99" {
		database.DB.Model(&order).Update("status", "待处理")
		OK(c, nil)
		return
	}
	// 调用上游补刷
	res := checkorder.RetryOrder(order.OID)
	if res.Code == 1 {
		database.DB.Model(&order).Updates(map[string]interface{}{
			"status": "补刷中", "bsnum": bsNum + 1,
		})
		OK(c, nil)
		return
	}
	Fail(c, res.Msg)
}

// POST /api/order/status  批量修改任务状态（管理员）
func OrderBatchStatus(c *gin.Context) {
	var req struct {
		OIDs   []int64 `json:"oids" binding:"required"`
		Status string  `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	database.DB.Model(&model.Order{}).Where("oid IN ?", req.OIDs).Update("status", req.Status)
	OK(c, nil)
}

// POST /api/order/dock  批量修改处理状态（管理员）
func OrderBatchDock(c *gin.Context) {
	var req struct {
		OIDs []int64 `json:"oids" binding:"required"`
		Dock string  `json:"dock" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	database.DB.Model(&model.Order{}).Where("oid IN ?", req.OIDs).Update("dockstatus", req.Dock)
	OK(c, nil)
}

// POST /api/order/refund  退款（管理员）
func OrderRefund(c *gin.Context) {
	var req struct {
		OIDs []int64 `json:"oids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	for _, oid := range req.OIDs {
		var order model.Order
		if err := database.DB.First(&order, oid).Error; err != nil {
			continue
		}
		fees, _ := strconv.ParseFloat(order.Fees, 64)
		database.DB.Model(&model.User{}).Where("uid = ?", order.UID).
			UpdateColumn("money", gorm.Expr("money + ?", fees))
		database.DB.Model(&order).Updates(map[string]interface{}{
			"status": "已退款", "dockstatus": "4",
		})
		writeLog(order.UID, "退款", fmt.Sprintf("订单%d退款%.2f元", oid, fees), fees, "")
	}
	OK(c, nil)
}

// GET /api/order/export
func OrderExport(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)

	gs := c.Query("gs")
	cid := c.Query("cid")
	statusText := c.Query("status_text")

	db := database.DB.Model(&model.Order{})
	if uidInt != 1 {
		db = db.Where("uid = ?", uidInt)
	}
	if cid != "" {
		db = db.Where("cid = ?", cid)
	}
	if statusText != "" {
		db = db.Where("status = ?", statusText)
	}

	var orders []model.Order
	db.Find(&orders)

	var lines []string
	for _, o := range orders {
		var line string
		switch gs {
		case "1":
			line = fmt.Sprintf("%s+%s+%s+%s", o.School, o.UserAccount, o.Pass, o.KCName)
		case "2":
			line = fmt.Sprintf("%s+%s+%s", o.UserAccount, o.Pass, o.KCName)
		case "3":
			line = fmt.Sprintf("%s+%s+%s", o.School, o.UserAccount, o.Pass)
		case "4":
			line = fmt.Sprintf("%s+%s", o.UserAccount, o.Pass)
		default:
			line = fmt.Sprintf("%s+%s+%s+%s", o.School, o.UserAccount, o.Pass, o.KCName)
		}
		lines = append(lines, line)
	}
	OK(c, strings.Join(lines, "\n"))
}

func calcPrice(class model.Class, user model.User) float64 {
	price, _ := strconv.ParseFloat(class.Price, 64)
	var danjia float64
	if class.YunSuan == "+" {
		danjia = price + user.AddPrice
	} else {
		danjia = price * user.AddPrice
	}

	// 密价
	var mijia model.MiJia
	if err := database.DB.Where("uid = ? AND cid = ?", user.UID, class.CID).First(&mijia).Error; err == nil {
		p, _ := strconv.ParseFloat(mijia.Price, 64)
		switch mijia.Mode {
		case 0:
			danjia = danjia - p
		case 1:
			danjia = (price - p) * user.AddPrice
		case 2:
			danjia = p
		}
		if danjia <= 0 {
			danjia = 0
		}
		orig := price * user.AddPrice
		if danjia > orig {
			danjia = orig
		}
	}
	return danjia
}

func writeLog(uid int64, typ, text string, money float64, ip string) {
	var user model.User
	database.DB.First(&user, uid)
	log := model.Log{
		UID:    uid,
		Type:   typ,
		Text:   text,
		Money:  fmt.Sprintf("%.2f", money),
		SMoney: fmt.Sprintf("%.2f", user.Money),
		IP:     ip,
	}
	database.DB.Create(&log)
}

// POST /api/order/query  查课（提交任务前预查询）
func OrderQuery(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)

	var req struct {
		CID    int64  `json:"cid" binding:"required"`
		School string `json:"school" binding:"required"`
		User   string `json:"user" binding:"required"`
		Pass   string `json:"pass" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数不能为空")
		return
	}

	var user model.User
	database.DB.First(&user, uidInt)

	// zdmoney 余额限制
	zdmoney := getConfFloat("zdmoney", 10)
	if user.Money < zdmoney {
		Fail(c, fmt.Sprintf("余额小于%v禁止查课", zdmoney))
		return
	}
	// 查课/下单比例限制
	var ckCount, orderCount int64
	database.DB.Model(&model.Log{}).Where("uid = ? AND type = ?", uidInt, "API查课").Count(&ckCount)
	database.DB.Model(&model.Log{}).Where("uid = ? AND type IN ?", uidInt, []string{"添加任务", "API添加任务", "下单"}).Count(&orderCount)
	taskTotal := orderCount
	if taskTotal == 0 {
		taskTotal = 1
	}
	blLimit := getConfFloat("bl", 999)
	if float64(user.CK)/float64(taskTotal) > blLimit {
		Fail(c, fmt.Sprintf("查课比例超%v禁止查课", blLimit))
		return
	}

	res := checkorder.QueryByClass(req.CID, req.School, req.User, req.Pass)
	if res == nil {
		Fail(c, "查询失败")
		return
	}

	// 增加查课计数
	database.DB.Model(&model.User{}).Where("uid = ?", uidInt).
		UpdateColumn("ck", gorm.Expr("ck + 1"))
	writeLog(uidInt, "查课", fmt.Sprintf("查询 %s %s", req.School, req.User), 0, c.ClientIP())

	OK(c, res)
}

// POST /api/order/manual-dock  立即对接（管理员）
func OrderManualDock(c *gin.Context) {
	var req struct {
		OID int64 `json:"oid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	res := checkorder.AddOrder(req.OID)
	if res.Code == 1 {
		var order model.Order
		database.DB.First(&order, req.OID)
		var class model.Class
		database.DB.First(&class, order.CID)
		database.DB.Model(&order).Updates(map[string]interface{}{
			"hid":        class.Docking,
			"status":     "进行中",
			"dockstatus": "1",
			"yid":        res.YID,
		})
		OK(c, res)
		return
	}
	Fail(c, res.Msg)
}

// POST /api/order/sync-progress  同步进度（管理员）
func OrderSyncProgress(c *gin.Context) {
	var req struct {
		OID int64 `json:"oid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	items := checkorder.ProgressByOrder(req.OID)
	for _, p := range items {
		database.DB.Model(&model.Order{}).
			Where("oid = ?", req.OID).
			Updates(map[string]interface{}{
				"yid":             p.YID,
				"status":          p.StatusText,
				"courseStartTime": p.KCStartTime,
				"courseEndTime":   p.KCEndTime,
				"examStartTime":   p.KSStartTime,
				"examEndTime":     p.KSEndTime,
				"process":         p.Process,
			})
	}
	OK(c, items)
}
