package handler

import (
	"fmt"
	"strconv"
	"time"
	"wk-go/internal/database"
	"wk-go/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GET /api/admin/users
func AdminUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("current", c.DefaultQuery("page", "1")))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", c.DefaultQuery("pageSize", "15")))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 15
	}
	typ := c.Query("type")
	qq := c.Query("qq")

	db := database.DB.Model(&model.User{})
	if qq != "" {
		switch typ {
		case "1":
			db = db.Where("uid = ?", qq)
		case "2":
			db = db.Where("user LIKE ?", "%"+qq+"%")
		case "3":
			db = db.Where("yqm = ?", qq)
		case "4":
			db = db.Where("name LIKE ?", "%"+qq+"%")
		case "5":
			db = db.Where("addprice = ?", qq)
		case "6":
			db = db.Where("money = ?", qq)
		default:
			db = db.Where("user LIKE ? OR name LIKE ?", "%"+qq+"%", "%"+qq+"%")
		}
	}

	var total int64
	db.Count(&total)
	var users []model.User
	db.Order("uid DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)
	PageResult(c, users, total, page, pageSize)
}

// POST /api/admin/user/add
func AdminUserAdd(c *gin.Context) {
	var req struct {
		Name     string  `json:"name" binding:"required"`
		User     string  `json:"user" binding:"required"`
		Pass     string  `json:"pass" binding:"required"`
		AddPrice float64 `json:"addprice"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	var count int64
	database.DB.Model(&model.User{}).Where("user = ?", req.User).Count(&count)
	if count > 0 {
		Fail(c, "账号已存在")
		return
	}
	if req.AddPrice <= 0 {
		req.AddPrice = 1.0
	}
	// 费率必须为 0.05 倍数
	if m := req.AddPrice / 0.05; m-float64(int(m)) > 1e-9 {
		Fail(c, "费率必须为0.05的倍数")
		return
	}
	user := model.User{
		UUID:     1,
		User:     req.User,
		Pass:     req.Pass,
		Name:     req.Name,
		AddPrice: req.AddPrice,
		Key:      "0",
		AddTime:  time.Now().Format("2006-01-02 15:04:05"),
		Active:   "1",
	}
	database.DB.Create(&user)
	writeLog(user.UID, "开户", fmt.Sprintf("管理员开户 %s", req.User), 0, "")
	OK(c, nil)
}

// POST /api/admin/user/recharge
func AdminUserRecharge(c *gin.Context) {
	var req struct {
		UID   int64   `json:"uid" binding:"required"`
		Money float64 `json:"money" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	database.DB.Model(&model.User{}).Where("uid = ?", req.UID).
		UpdateColumn("money", gorm.Expr("money + ?", req.Money))
	database.DB.Model(&model.User{}).Where("uid = ?", req.UID).
		UpdateColumn("zcz", gorm.Expr("zcz + ?", req.Money))
	writeLog(req.UID, "管理员充值", fmt.Sprintf("管理员充值%.2f元", req.Money), req.Money, "")
	OK(c, nil)
}

// POST /api/admin/user/ban
func AdminUserBan(c *gin.Context) {
	var req struct {
		UID    int64  `json:"uid" binding:"required"`
		Active string `json:"active" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	database.DB.Model(&model.User{}).Where("uid = ?", req.UID).Update("active", req.Active)
	OK(c, nil)
}

// POST /api/admin/user/reset-pass
func AdminUserResetPass(c *gin.Context) {
	var req struct {
		UID     int64  `json:"uid" binding:"required"`
		NewPass string `json:"new_pass" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	database.DB.Model(&model.User{}).Where("uid = ?", req.UID).Update("pass", req.NewPass)
	OK(c, nil)
}

// POST /api/admin/user/level
func AdminUserLevel(c *gin.Context) {
	var req struct {
		UID      int64   `json:"uid" binding:"required"`
		AddPrice float64 `json:"addprice" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	database.DB.Model(&model.User{}).Where("uid = ?", req.UID).Update("addprice", req.AddPrice)
	OK(c, nil)
}

// POST /api/admin/user/api-key
func AdminUserOpenKey(c *gin.Context) {
	var req struct {
		UID int64 `json:"uid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	key := uuid.New().String()
	database.DB.Model(&model.User{}).Where("uid = ?", req.UID).Update("key", key)
	OK(c, gin.H{"key": key})
}

// POST /api/admin/user/yqm
func AdminUserYQM(c *gin.Context) {
	var req struct {
		UID int64  `json:"uid" binding:"required"`
		YQM string `json:"yqm" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	database.DB.Model(&model.User{}).Where("uid = ?", req.UID).Update("yqm", req.YQM)
	OK(c, nil)
}

// GET /api/admin/config
func AdminGetConfig(c *gin.Context) {
	var configs []model.Config
	database.DB.Find(&configs)
	result := make(map[string]string)
	for _, cfg := range configs {
		result[cfg.V] = cfg.K
	}
	OK(c, result)
}

// POST /api/admin/config
func AdminSaveConfig(c *gin.Context) {
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		Fail(c, "参数错误: "+err.Error())
		return
	}
	for k, v := range raw {
		vStr := fmt.Sprintf("%v", v)
		if v == nil {
			vStr = ""
		}
		var count int64
		database.DB.Model(&model.Config{}).Where("v = ?", k).Count(&count)
		if count > 0 {
			// 已存在则更新
			if err := database.DB.Model(&model.Config{}).Where("v = ?", k).Update("k", vStr).Error; err != nil {
				Fail(c, "保存配置失败: "+err.Error())
				return
			}
		} else {
			// 行不存在则创建（新配置项），避免保存后刷新丢失
			if err := database.DB.Create(&model.Config{V: k, K: vStr}).Error; err != nil {
				Fail(c, "保存配置失败: "+err.Error())
				return
			}
		}
	}
	OK(c, nil)
}

// GET /api/admin/dengji
func DengjiList(c *gin.Context) {
	var list []model.Dengji
	database.DB.Order("sort ASC").Find(&list)
	OK(c, list)
}

// POST /api/admin/dengji
func DengjiAdd(c *gin.Context) {
	var req model.Dengji
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	req.Time = time.Now().Format("2006-01-02 15:04:05")
	database.DB.Create(&req)
	OK(c, req)
}

// PUT /api/admin/dengji/:id
func DengjiUpdate(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	database.DB.Model(&model.Dengji{}).Where("id = ?", id).Updates(req)
	OK(c, nil)
}

// DELETE /api/admin/dengji/:id
func DengjiDelete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&model.Dengji{}, id)
	OK(c, nil)
}

// GET /api/admin/log
func AdminLogList(c *gin.Context) {
	uid := c.Query("uid")
	typ := c.Query("type")
	page, _ := atoiSafe(c.DefaultQuery("current", c.DefaultQuery("page", "1")))
	pageSize, _ := atoiSafe(c.DefaultQuery("size", c.DefaultQuery("pageSize", "20")))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	db := database.DB.Model(&model.Log{})
	if uid != "" {
		db = db.Where("uid = ?", uid)
	}
	if typ != "" {
		db = db.Where("type = ?", typ)
	}
	var total int64
	db.Count(&total)
	var logs []model.Log
	db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)
	PageResult(c, logs, total, page, pageSize)
}

// GET /api/admin/stats  首页统计
func AdminStats(c *gin.Context) {
	var totalUsers, totalOrders, todayOrders int64
	var totalMoney float64

	database.DB.Model(&model.User{}).Count(&totalUsers)
	database.DB.Model(&model.Order{}).Count(&totalOrders)
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&model.Order{}).Where("addtime LIKE ?", today+"%").Count(&todayOrders)
	database.DB.Model(&model.Order{}).Select("COALESCE(SUM(fees),0)").Scan(&totalMoney)

	OK(c, gin.H{
		"total_users":  totalUsers,
		"total_orders": totalOrders,
		"today_orders": todayOrders,
		"total_money":  totalMoney,
	})
}

// GET /api/admin/dashboard  完整数据统计
func AdminDashboard(c *gin.Context) {
	var totalUsers, activeUsers, totalOrders, todayOrders, monthOrders int64
	var totalRevenue, totalConsume, todayRevenue float64
	var pendingOrders, doneOrders, errorOrders, refundedOrders int64

	today := time.Now().Format("2006-01-02")
	monthStart := time.Now().Format("2006-01") + "-01"

	database.DB.Model(&model.User{}).Count(&totalUsers)
	database.DB.Model(&model.User{}).Where("active = ?", "1").Count(&activeUsers)

	database.DB.Model(&model.Order{}).Count(&totalOrders)
	database.DB.Model(&model.Order{}).Where("addtime LIKE ?", today+"%").Count(&todayOrders)
	database.DB.Model(&model.Order{}).Where("addtime >= ?", monthStart).Count(&monthOrders)

	database.DB.Model(&model.Order{}).Where("status = ?", "待处理").Count(&pendingOrders)
	database.DB.Model(&model.Order{}).Where("status = ?", "已完成").Count(&doneOrders)
	database.DB.Model(&model.Order{}).Where("status = ?", "异常").Count(&errorOrders)
	database.DB.Model(&model.Order{}).Where("status = ?", "已退款").Count(&refundedOrders)

	database.DB.Model(&model.Pay{}).Where("status = 1").Select("COALESCE(SUM(money),0)").Scan(&totalRevenue)
	database.DB.Model(&model.Pay{}).Where("status = 1 AND addtime >= ?", today).
		Select("COALESCE(SUM(money),0)").Scan(&todayRevenue)
	database.DB.Model(&model.Order{}).Select("COALESCE(SUM(fees),0)").Scan(&totalConsume)

	OK(c, gin.H{
		"users": gin.H{
			"total":  totalUsers,
			"active": activeUsers,
		},
		"orders": gin.H{
			"total":    totalOrders,
			"today":    todayOrders,
			"month":    monthOrders,
			"pending":  pendingOrders,
			"done":     doneOrders,
			"error":    errorOrders,
			"refunded": refundedOrders,
		},
		"revenue": gin.H{
			"total":   totalRevenue,
			"today":   todayRevenue,
			"consume": totalConsume,
		},
	})
}

// GET /api/order/by-platform  按平台统计订单数
func OrderByPlatform(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)

	type row struct {
		PTName string
		Count  int64
	}
	var rows []row

	db := database.DB.Model(&model.Order{}).
		Select("ptname, COUNT(*) as count").
		Group("ptname").
		Order("count DESC").
		Limit(9)
	if uidInt != 1 {
		db = db.Where("uid = ?", uidInt)
	}
	db.Scan(&rows)

	labels := make([]string, 0, len(rows))
	values := make([]int64, 0, len(rows))
	for _, r := range rows {
		labels = append(labels, r.PTName)
		values = append(values, r.Count)
	}
	OK(c, gin.H{"labels": labels, "values": values})
}

// GET /api/order/monthly  近 7 天每日订单数
func OrderMonthly(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)

	type row struct {
		Ymd   string
		Count int64
	}
	var rows []row

	db := database.DB.Model(&model.Order{}).
		Select("DATE_FORMAT(addtime,'%Y-%m-%d') as ymd, COUNT(*) as count").
		Where("addtime >= ?", time.Now().AddDate(0, 0, -6).Format("2006-01-02")).
		Group("ymd").
		Order("ymd ASC")
	if uidInt != 1 {
		db = db.Where("uid = ?", uidInt)
	}
	db.Scan(&rows)

	now := time.Now()
	labels := make([]string, 7)
	values := make([]int64, 7)
	keys := make([]string, 7)
	for i := 0; i < 7; i++ {
		t := now.AddDate(0, 0, -6+i)
		keys[i] = t.Format("2006-01-02")
		labels[i] = fmt.Sprintf("%d/%d", int(t.Month()), t.Day())
	}
	for _, r := range rows {
		for i, k := range keys {
			if k == r.Ymd {
				values[i] = r.Count
			}
		}
	}
	OK(c, gin.H{"labels": labels, "values": values})
}
