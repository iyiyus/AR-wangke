package handler

import (
	"fmt"
	"regexp"
	"time"
	"wk-go/internal/database"
	"wk-go/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func uuidNew() string {
	return uuid.New().String()
}

// 代理对下级的操作（兼容原 PHP apisub：adduser/userjk/usergj/user_czmm/user_ban/szyqm/ktapi type2）

// getMySub 校验目标 uid 是自己的下级，返回下级用户
func getMySub(c *gin.Context, uidInt int64, subUID int64) *model.User {
	var sub model.User
	if err := database.DB.First(&sub, subUID).Error; err != nil {
		Fail(c, "用户不存在")
		return nil
	}
	if sub.UUID != uidInt {
		Fail(c, "只能操作自己的下级")
		return nil
	}
	return &sub
}

func loadConfMap() map[string]string {
	var configs []model.Config
	database.DB.Find(&configs)
	conf := make(map[string]string)
	for _, cfg := range configs {
		conf[cfg.V] = cfg.K
	}
	return conf
}

// GET /api/proxy/user/list  下级列表
func ProxyUserList(c *gin.Context) {
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
	q := c.Query("qq")
	typ := c.Query("type")

	db := database.DB.Model(&model.User{}).Where("uuid = ?", uidInt)
	if q != "" {
		switch typ {
		case "1":
			db = db.Where("uid = ?", q)
		case "2":
			db = db.Where("user LIKE ?", "%"+q+"%")
		case "3":
			db = db.Where("yqm = ?", q)
		case "4":
			db = db.Where("name LIKE ?", "%"+q+"%")
		case "5":
			db = db.Where("addprice = ?", q)
		case "6":
			db = db.Where("money = ?", q)
		default:
			db = db.Where("user LIKE ? OR name LIKE ?", "%"+q+"%", "%"+q+"%")
		}
	}
	var total int64
	db.Count(&total)
	var users []model.User
	db.Order("uid DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)
	PageResult(c, users, total, page, pageSize)
}

// POST /api/proxy/user/add  代理开户（adduser）
func ProxyUserAdd(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)
	var me model.User
	database.DB.First(&me, uidInt)

	var req struct {
		Name     string  `json:"name" binding:"required"`
		User     string  `json:"user" binding:"required"`
		Pass     string  `json:"pass" binding:"required"`
		AddPrice float64 `json:"addprice" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	conf := loadConfMap()
	if conf["user_htkh"] != "1" {
		Fail(c, "系统未开启代理开户")
		return
	}
	// QQ 号格式（原 PHP 用 QQ 正则校验账号）
	qqRe := regexp.MustCompile(`^[1-9][0-9]{4,11}$`)
	if !qqRe.MatchString(req.User) {
		Fail(c, "账号必须为QQ号")
		return
	}
	if len(req.Pass) < 6 {
		Fail(c, "密码至少6位")
		return
	}
	// 费率必须为 0.05 倍数且不低于上级
	if m := req.AddPrice / 0.05; m-float64(int(m)) > 1e-9 {
		Fail(c, "费率必须为0.05的倍数")
		return
	}
	if req.AddPrice < me.AddPrice {
		Fail(c, "下级费率不能低于你的费率")
		return
	}

	var count int64
	database.DB.Model(&model.User{}).Where("user = ?", req.User).Count(&count)
	if count > 0 {
		Fail(c, "账号已存在")
		return
	}

	// 开户费
	if ktmoney := parsePrice(conf["user_ktmoney"]); ktmoney > 0 {
		if me.Money < ktmoney {
			Fail(c, fmt.Sprintf("余额不足，开户需%.2f元", ktmoney))
			return
		}
		database.DB.Model(&me).UpdateColumn("money", gorm.Expr("money - ?", ktmoney))
		writeLog(me.UID, "代理开户", fmt.Sprintf("开户 %s 扣除%.2f元", req.User, ktmoney), -ktmoney, "")
	}

	now := nowStr()
	user := model.User{
		UUID:     me.UID,
		User:     req.User,
		Pass:     req.Pass,
		Name:     req.Name,
		AddPrice: req.AddPrice,
		Key:      "0",
		AddTime:  now,
		IP:       c.ClientIP(),
		Active:   "1",
	}
	if err := database.DB.Create(&user).Error; err != nil {
		Fail(c, "开户失败")
		return
	}

	// 等级 addkf=1 时自动给下级充值 cz×(上级费率/新费率)
	var dj model.Dengji
	if err := database.DB.Where("rate = ? AND addkf = ?", req.AddPrice, "1").First(&dj).Error; err == nil {
		cz := dj.Money * (me.AddPrice / req.AddPrice)
		if cz > 0 {
			database.DB.Model(&user).UpdateColumn("money", gorm.Expr("money + ?", cz))
			writeLog(user.UID, "代理充值", fmt.Sprintf("开户自动充值%.2f元", cz), cz, "")
		}
	}
	writeLog(user.UID, "开户", fmt.Sprintf("由上级 %s 开户", me.User), 0, c.ClientIP())
	OK(c, nil)
}

// POST /api/proxy/user/recharge  给下级充值（userjk：按费率换算扣费）
func ProxyUserRecharge(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)
	var me model.User
	database.DB.First(&me, uidInt)

	var req struct {
		UID   int64   `json:"uid" binding:"required"`
		Money float64 `json:"money" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	sub := getMySub(c, uidInt, req.UID)
	if sub == nil {
		return
	}
	if req.Money <= 0 {
		Fail(c, "金额错误")
		return
	}
	cost := req.Money * (me.AddPrice / sub.AddPrice)
	if me.Money < cost {
		Fail(c, fmt.Sprintf("余额不足，需%.2f元", cost))
		return
	}
	database.DB.Model(&me).UpdateColumn("money", gorm.Expr("money - ?", cost))
	database.DB.Model(sub).UpdateColumn("money", gorm.Expr("money + ?", req.Money))
	writeLog(me.UID, "代理充值", fmt.Sprintf("给下级 %s 充值%.2f元，扣除%.2f元", sub.User, req.Money, cost), -cost, c.ClientIP())
	writeLog(sub.UID, "上级充值", fmt.Sprintf("上级 %s 给你充值%.2f元", me.User, req.Money), req.Money, c.ClientIP())
	OK(c, nil)
}

// POST /api/proxy/user/price  给下级改费率（usergj：只能上调，换算余额，3元手续费）
func ProxyUserPrice(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)
	var me model.User
	database.DB.First(&me, uidInt)

	var req struct {
		UID      int64   `json:"uid" binding:"required"`
		AddPrice float64 `json:"addprice" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	sub := getMySub(c, uidInt, req.UID)
	if sub == nil {
		return
	}
	if m := req.AddPrice / 0.05; m-float64(int(m)) > 1e-9 {
		Fail(c, "费率必须为0.05的倍数")
		return
	}
	if req.AddPrice < sub.AddPrice {
		Fail(c, "只能上调费率，下调费率请联系管理员")
		return
	}
	if req.AddPrice > me.AddPrice {
		Fail(c, "下级费率不能超过你的费率")
		return
	}

	// 改价手续费（等级 gjkf）
	fee := 3.0
	var dj model.Dengji
	if err := database.DB.Where("rate = ? AND gjkf = ?", me.AddPrice, "1").First(&dj).Error; err == nil {
		fee = dj.Money
	}
	if me.Money < fee {
		Fail(c, fmt.Sprintf("余额不足，改价需%.2f元手续费", fee))
		return
	}

	// 余额换算：余额/旧费率×新费率
	newMoney := sub.Money / sub.AddPrice * req.AddPrice
	database.DB.Model(&me).UpdateColumn("money", gorm.Expr("money - ?", fee))
	database.DB.Model(sub).Updates(map[string]interface{}{
		"addprice": req.AddPrice,
		"money":    newMoney,
	})
	writeLog(me.UID, "代理改价", fmt.Sprintf("给下级 %s 改费率 %.2f→%.2f，扣除%.2f元", sub.User, sub.AddPrice, req.AddPrice, fee), -fee, c.ClientIP())
	OK(c, nil)
}

// POST /api/proxy/user/reset-pass  重置下级密码为 123456（user_czmm）
func ProxyUserResetPass(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)
	var req struct {
		UID int64 `json:"uid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	sub := getMySub(c, uidInt, req.UID)
	if sub == nil {
		return
	}
	database.DB.Model(sub).Update("pass", "123456")
	writeLog(uidInt, "重置密码", fmt.Sprintf("重置下级 %s 密码", sub.User), 0, "")
	OK(c, nil)
}

// POST /api/proxy/user/ban  封禁/解封下级（user_ban）
func ProxyUserBan(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)
	var req struct {
		UID    int64  `json:"uid" binding:"required"`
		Active string `json:"active" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	sub := getMySub(c, uidInt, req.UID)
	if sub == nil {
		return
	}
	database.DB.Model(sub).Update("active", req.Active)
	OK(c, nil)
}

// POST /api/proxy/user/yqm  给下级设邀请码（szyqm）
func ProxyUserYQM(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)
	var req struct {
		UID int64  `json:"uid" binding:"required"`
		YQM string `json:"yqm" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	sub := getMySub(c, uidInt, req.UID)
	if sub == nil {
		return
	}
	var count int64
	database.DB.Model(&model.User{}).Where("yqm = ? AND uid != ?", req.YQM, sub.UID).Count(&count)
	if count > 0 {
		Fail(c, "邀请码已被使用")
		return
	}
	database.DB.Model(sub).Update("yqm", req.YQM)
	OK(c, nil)
}

// POST /api/proxy/user/api-key  给下级开通 API key（ktapi type2：扣 5 元）
func ProxyUserOpenKey(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)
	var me model.User
	database.DB.First(&me, uidInt)
	var req struct {
		UID int64 `json:"uid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	sub := getMySub(c, uidInt, req.UID)
	if sub == nil {
		return
	}
	if sub.Key != "0" {
		Fail(c, "该下级已开通接口")
		return
	}
	if me.Money < 5 {
		Fail(c, "余额不足，开通需扣除5元")
		return
	}
	database.DB.Model(&me).UpdateColumn("money", gorm.Expr("money - ?", 5))
	writeLog(me.UID, "开通接口", fmt.Sprintf("给下级 %s 开通API扣除5元", sub.User), -5, "")
	key := uuidNew()
	database.DB.Model(sub).Update("key", key)
	OK(c, gin.H{"key": key})
}

// POST /api/user/migrate  上级迁移（sjqy：需上级UID+邀请码，原上级7天未登录）
func UserMigrate(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)
	conf := loadConfMap()
	if conf["sjqykg"] != "1" {
		Fail(c, "系统未开启上级迁移")
		return
	}
	var req struct {
		BossUID int64  `json:"boss_uid" binding:"required"`
		YQM     string `json:"yqm" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	var me model.User
	database.DB.First(&me, uidInt)

	var boss model.User
	if err := database.DB.First(&boss, req.BossUID).Error; err != nil {
		Fail(c, "目标上级不存在")
		return
	}
	if boss.YQM != req.YQM {
		Fail(c, "邀请码错误")
		return
	}
	// 原上级 7 天内登录过则不允许迁移
	var oldBoss model.User
	if err := database.DB.First(&oldBoss, me.UUID).Error; err == nil {
		if t, err := time.Parse("2006-01-02 15:04:05", oldBoss.EndTime); err == nil {
			if time.Since(t) < 7*24*time.Hour {
				Fail(c, "原上级7天内有登录，无法迁移")
				return
			}
		}
	}
	database.DB.Model(&me).Update("uuid", boss.UID)
	writeLog(uidInt, "上级迁移", fmt.Sprintf("上级迁移至 %s", boss.User), 0, "")
	OK(c, nil)
}
