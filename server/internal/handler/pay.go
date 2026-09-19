package handler

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
	"wk-go/internal/database"
	"wk-go/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// POST /api/pay/create
func PayCreate(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)

	var req struct {
		Money string `json:"money" binding:"required"`
		Type  string `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	// 仅直属代理可在线充值（兼容原 PHP：uuid==1）
	var user model.User
	database.DB.First(&user, uidInt)
	if user.UUID != 1 {
		Fail(c, "仅直属代理支持在线充值，请联系上级充值")
		return
	}

	// 获取支付配置
	conf := loadConfMap()

	// 最低充值限制
	if min := parsePrice(conf["zdpay"]); min > 0 && parsePrice(req.Money) < min {
		Fail(c, fmt.Sprintf("最低充值%.2f元", min))
		return
	}

	outTradeNo := time.Now().Format("20060102150405") + fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	now := time.Now()

	pay := model.Pay{
		OutTradeNo: outTradeNo,
		UID:        uidInt,
		Num:        1,
		Type:       req.Type,
		AddTime:    &now,
		Name:       "账户充值",
		Money:      req.Money,
		IP:         c.ClientIP(),
		Domain:     c.Request.Host,
		Status:     0,
	}
	database.DB.Create(&pay)

	// 构造易支付跳转参数
	params := map[string]string{
		"pid":          conf["epay_pid"],
		"type":         req.Type,
		"out_trade_no": outTradeNo,
		"notify_url":   "https://" + c.Request.Host + "/api/pay/notify",
		"return_url":   "https://" + c.Request.Host + "/pay/result",
		"name":         "账户充值",
		"money":        req.Money,
	}
	sign := epaySign(params, conf["epay_key"])
	params["sign"] = sign
	params["sign_type"] = "MD5"

	payURL := conf["epay_api"] + "submit.php?" + buildQuery(params)
	OK(c, gin.H{"pay_url": payURL, "out_trade_no": outTradeNo})
}

// POST /api/pay/notify  易支付异步回调
func PayNotify(c *gin.Context) {
	c.Request.ParseForm()
	params := make(map[string]string)
	for k, v := range c.Request.Form {
		if k != "sign" && k != "sign_type" {
			params[k] = v[0]
		}
	}

	var configs []model.Config
	database.DB.Find(&configs)
	conf := make(map[string]string)
	for _, cfg := range configs {
		conf[cfg.V] = cfg.K
	}

	sign := epaySign(params, conf["epay_key"])
	if sign != c.PostForm("sign") {
		c.String(200, "fail")
		return
	}

	outTradeNo := c.PostForm("out_trade_no")
	tradeNo := c.PostForm("trade_no")
	tradeStatus := c.PostForm("trade_status")

	if tradeStatus != "TRADE_SUCCESS" {
		c.String(200, "fail")
		return
	}

	var pay model.Pay
	if err := database.DB.Where("out_trade_no = ? AND status = 0", outTradeNo).First(&pay).Error; err != nil {
		c.String(200, "success")
		return
	}

	now := time.Now()
	database.DB.Model(&pay).Updates(map[string]interface{}{
		"trade_no": tradeNo,
		"status":   1,
		"endtime":  now,
	})

	money := 0.0
	fmt.Sscanf(pay.Money, "%f", &money)
	database.DB.Model(&model.User{}).Where("uid = ?", pay.UID).
		UpdateColumn("money", gorm.Expr("money + ?", money))
	database.DB.Model(&model.User{}).Where("uid = ?", pay.UID).
		UpdateColumn("zcz", gorm.Expr("zcz + ?", money))
	writeLog(pay.UID, "充值", fmt.Sprintf("充值%.2f元，订单号%s", money, outTradeNo), money, pay.IP)

	c.String(200, "success")
}

// GET /api/pay/list
func PayList(c *gin.Context) {
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

	db := database.DB.Model(&model.Pay{})
	if uidInt != 1 {
		db = db.Where("uid = ?", uidInt)
	}
	var total int64
	db.Count(&total)
	var pays []model.Pay
	db.Order("oid DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&pays)
	PageResult(c, pays, total, page, pageSize)
}

func epaySign(params map[string]string, key string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		if params[k] != "" {
			parts = append(parts, k+"="+params[k])
		}
	}
	str := strings.Join(parts, "&") + key
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func buildQuery(params map[string]string) string {
	vals := url.Values{}
	for k, v := range params {
		vals.Set(k, v)
	}
	return vals.Encode()
}
