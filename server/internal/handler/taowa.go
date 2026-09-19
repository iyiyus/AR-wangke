package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"wk-go/internal/database"
	"wk-go/internal/model"
)

/* ============================================================
 * taowa 盖章/病历对接（对应 PHP 对接包 sxzsapi.php / sxzsconfig.php）
 * 配置项（config 表，管理员在系统设置填写）：
 *   taowa_url       对接平台地址（erkaiapi_url，默认 https://www.sxzsjk.top）
 *   taowa_uid       你的 UID（erkaiapi 账号）
 *   taowa_key       你的 KEY（erkaiapi 密钥）
 *   taowa_jg        盖章价格倍数（默认 2）
 *   taowa_fenlei    盖章分类 ID（源台分类设置里查看）
 *   taowa_hid       盖章平台 ID
 *   taowa_bljg      病历价格倍数（默认 2）
 *   taowa_blfenlei  病历分类 ID（源台分类设置里查看）
 *   taowa_blhid     病历平台 ID
 *   taowa_danjia    源台价格乘数（默认 1.05）
 *   taowa_newdanjia 源台改价乘数（默认 1.08）
 * 说明：公司/模板/规格/配送/物流等查询走源台公开接口；
 *       下单走 erkaiapi 协议（act=add/addbl，带 uid/key），
 *       同时本地记录订单（order 表 kcid=gz/bl）并扣费。
 * ============================================================ */

func taowaConf(key, def string) string {
	var c model.Config
	if err := database.DB.Where("v = ?", key).First(&c).Error; err == nil && c.K != "" {
		return c.K
	}
	return def
}

func taowaBase() string {
	u := taowaConf("taowa_url", "https://www.sxzsjk.top")
	return strings.TrimRight(u, "/")
}

// taowaCredsReady 是否已填写源台凭证
func taowaCredsReady() bool {
	uid := taowaConf("taowa_uid", "")
	key := taowaConf("taowa_key", "")
	return uid != "" && key != "" && uid != "你的UID" && key != "你的KEY"
}

// taowaProxy 调源台 erkaiapi：POST form {base}/sxzstwapi.php?act=xxx {uid,key,...}
func taowaProxy(act string, params gin.H) gin.H {
	if !taowaCredsReady() {
		return gin.H{"code": 1, "msg": "未配置对接账号，请在系统设置中填写UID与KEY"}
	}
	form := url.Values{}
	form.Set("uid", taowaConf("taowa_uid", ""))
	form.Set("key", taowaConf("taowa_key", ""))
	for k, v := range params {
		form.Set(k, fmt.Sprintf("%v", v))
	}
	target := taowaBase() + "/sxzstwapi.php?act=" + url.QueryEscape(act)
	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest("POST", target, strings.NewReader(form.Encode()))
	if err != nil {
		return gin.H{"code": 1, "msg": "请求构造失败"}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		return gin.H{"code": 1, "msg": "源台连接失败"}
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var out gin.H
	if json.Unmarshal(respBody, &out) != nil {
		return gin.H{"code": 1, "msg": "源台响应异常"}
	}
	return out
}

// taowaPublicGet 源台公开 GET 接口（class.php / bl.php / jg.php / sf.php / chaxun.php / gonggao.php）
func taowaPublicGet(path string, query string) ([]byte, error) {
	url := path
	if !strings.HasPrefix(url, "http") {
		url = taowaBase() + path
	}
	if query != "" {
		url += "?" + query
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// 解析源台标准响应 {code, msg, body}
func taowaBody(res gin.H) (gin.H, bool) {
	code, _ := res["code"].(float64)
	if code != 1 {
		return nil, false
	}
	if body, ok := res["body"].(map[string]interface{}); ok {
		return body, true
	}
	if body, ok := res["body"].(gin.H); ok {
		return body, true
	}
	return gin.H{}, true
}

// ---------- 查询类 ----------

// GET /api/taowa/companies 盖章公司列表（直接调上游 class.php）
func TaowaCompanies(c *gin.Context) {
	data, err := taowaPublicGet("https://www.sxzsjk.top/class.php", "")
	if err != nil {
		OK(c, gin.H{"companies": []gin.H{}, "configured": true, "msg": "上游不可达"})
		return
	}
	var parsed struct {
		Code int             `json:"code"`
		Data []interface{}   `json:"data"`
	}
	if json.Unmarshal(data, &parsed) != nil {
		OK(c, gin.H{"companies": []gin.H{}, "configured": true, "msg": "上游数据格式错误"})
		return
	}
	OK(c, gin.H{"companies": parsed.Data, "configured": true})
}

// GET /api/taowa/templates 病历模板列表（直接调上游 bl.php）
func TaowaTemplates(c *gin.Context) {
	data, err := taowaPublicGet("https://www.sxzsjk.top/bl.php", "")
	if err != nil {
		OK(c, gin.H{"templates": []gin.H{}, "configured": true, "msg": "上游不可达"})
		return
	}
	var parsed struct {
		Code int           `json:"code"`
		Data []interface{} `json:"data"`
	}
	if json.Unmarshal(data, &parsed) != nil {
		OK(c, gin.H{"templates": []gin.H{}, "configured": true, "msg": "上游数据格式错误"})
		return
	}
	OK(c, gin.H{"templates": parsed.Data, "configured": true})
}

// POST /api/taowa/spec 规格查询（源台 getxuanx）
func TaowaSpec(c *gin.Context) {
	var req struct {
		CompanyName string `json:"companyName"`
		Userinfo    string `json:"userinfo"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	if req.CompanyName == "" {
		Fail(c, "请选择公司")
		return
	}
	res := taowaProxy("getxuanx", gin.H{
		"companyName": req.CompanyName,
		"userinfo":    req.Userinfo,
	})
	if body, ok := taowaBody(res); ok {
		OK(c, body["data"])
		return
	}
	Fail(c, res["msg"].(string))
}

// GET /api/taowa/delivery 配送方式（源台公开 chaxun.php）
func TaowaDelivery(c *gin.Context) {
	data, err := taowaPublicGet("/chaxun.php", "")
	if err != nil {
		// 源台不可达时降级返回空，避免前端报错
		OK(c, gin.H{"delivery": []gin.H{}, "msg": "配送方式获取失败"})
		return
	}
	// 校验响应为合法 JSON，否则降级
	var parsed interface{}
	if json.Unmarshal(data, &parsed) != nil {
		OK(c, gin.H{"delivery": []gin.H{}, "msg": "配送方式获取失败"})
		return
	}
	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
}

// GET /api/taowa/notice 源台公告（直接调上游 gonggao.php）
func TaowaNotice(c *gin.Context) {
	data, err := taowaPublicGet("https://www.sxzsjk.top/gonggao.php", "")
	if err != nil || len(data) == 0 {
		OK(c, gin.H{"notice": taowaConf("tcgonggao", "")})
		return
	}
	var parsed struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			ID      string `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"data"`
	}
	if json.Unmarshal(data, &parsed) == nil && len(parsed.Data) > 0 {
		// 拼接所有公告
		notice := ""
		for _, g := range parsed.Data {
			notice += "【" + g.Title + "】\n" + g.Content + "\n\n"
		}
		OK(c, gin.H{"notice": strings.TrimSpace(notice)})
		return
	}
	OK(c, gin.H{"notice": taowaConf("tcgonggao", "")})
}

// GET /api/taowa/order/ddlog 物流查询（源台公开 sf.php?sfdh=）
func TaowaDdlog(c *gin.Context) {
	sfdh := strings.TrimSpace(c.Query("sfdh"))
	if sfdh == "" {
		Fail(c, "请输入快递单号")
		return
	}
	data, err := taowaPublicGet("/sf.php", "sfdh="+sfdh)
	if err != nil {
		Fail(c, "物流查询失败")
		return
	}
	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
}

// ---------- 下单 ----------

// 本地写盖章/病历订单
func taowaCreateOrder(c *gin.Context, typ string, uid int64, ptName, school, user, pass, phone string,
	kcName string, specData interface{}, fileName string, price float64, remark string) (model.Order, string) {
	userRow := model.User{}
	if err := database.DB.First(&userRow, uid).Error; err != nil {
		return model.Order{}, "用户不存在"
	}
	bei := parsePrice(taowaConf("taowa_jg", "2"))
	hid := taowaConf("taowa_hid", "")
	fenlei := taowaConf("taowa_fenlei", "")
	if typ == "bl" {
		bei = parsePrice(taowaConf("taowa_bljg", "2"))
		hid = taowaConf("taowa_blhid", "")
		fenlei = taowaConf("taowa_blfenlei", "")
	}
	if bei <= 0 {
		bei = 2
	}
	money := round2(price * bei)
	if userRow.Money < money {
		return model.Order{}, "余额不足"
	}
	// 源台下单
	var res gin.H
	if typ == "gz" {
		res = taowaProxy("add", gin.H{
			"companyName": ptName,
			"userinfo":    school,
			"userinfo2":   user,
			"userinfo3":   pass,
			"data":        specData,
			"fileName":    fileName,
			"fenlei":      fenlei,
		})
	} else {
		res = taowaProxy("addbl", gin.H{
			"templateName":  ptName,
			"patientName":   school,
			"patientGender": "",
			"patientAge":    "",
			"department":    "",
			"condition":     "",
			"treatment":     "",
			"diagnosisDate": "",
			"hospital":      "",
			"email":         user,
			"data":          specData,
			"fenlei":        fenlei,
		})
	}
	dockStatus := "0"
	if code, _ := res["code"].(float64); code != 1 {
		dockStatus = "0" // 待对接（未配置或失败留在本地队列）
	}
	specJSON, _ := json.Marshal(specData)
	if specJSON == nil {
		specJSON = []byte("[]")
	}
	order := model.Order{
		UID:             uid,
		HID:             parseInt64(hid),
		PTName:          ptName,
		School:          school,
		UserAccount:     user,
		Pass:            pass,
		Phone:           phone,
		KCID:            typ,
		KCName:          kcName,
		CourseStartTime: string(specJSON),
		CourseEndTime:   fileName,
		Fees:            fmt.Sprintf("%.2f", money),
		Noun:            fmt.Sprintf("%.2f", price),
		Remarks:         remark,
		AddTime:         nowStr(),
		IP:              c.ClientIP(),
		DockStatus:      dockStatus,
		Status:          "待处理",
	}
	if err := database.DB.Create(&order).Error; err != nil {
		return model.Order{}, "订单保存失败"
	}
	database.DB.Model(&userRow).Update("money", userRow.Money-money)
	writeLog(uid, "taowa", fmt.Sprintf("%s下单 %s %s 扣除%.2f", typ, ptName, school, money), -money, c.ClientIP())
	return order, ""
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// POST /api/taowa/order/add 盖章下单
func TaowaOrderAdd(c *gin.Context) {
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(int64)
	var req struct {
		CompanyName string        `json:"companyName"`
		Name        string        `json:"name"`
		Phone       string        `json:"phone"`
		Address     string        `json:"address"`
		Remark      string        `json:"remark"`
		Data        []interface{} `json:"data"`
		FileName    string        `json:"fileName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	if req.CompanyName == "" || req.Name == "" || req.Phone == "" || req.Address == "" {
		Fail(c, "请填写完整信息")
		return
	}
	if len(req.Data) == 0 {
		Fail(c, "请选择规格")
		return
	}
	// 算价：规格价格合计
	price := 0.0
	names := []string{}
	for _, d := range req.Data {
		if m, ok := d.(map[string]interface{}); ok {
			if dd, ok := m["data"].(map[string]interface{}); ok {
				names = append(names, toStr(dd["name"]))
				price += toFloat(dd["price"])
			}
		}
	}
	if price <= 0 {
		// 尝试从源台查价
		Fail(c, "规格价格无效")
		return
	}
	remark := req.Remark
	if remark == "" {
		remark = req.Address
	}
	order, msg := taowaCreateOrder(c, "gz", uid, req.CompanyName, req.Name, req.Phone,
		req.Address, req.Phone, strings.Join(names, ","), req.Data, req.FileName, price, remark)
	if msg != "" {
		Fail(c, msg)
		return
	}
	OK(c, gin.H{"oid": order.OID, "msg": fmt.Sprintf("下单成功，扣费%s元！", order.Fees)})
}

// POST /api/taowa/order/add-bl 病历下单
func TaowaOrderAddBL(c *gin.Context) {
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(int64)
	var req struct {
		TemplateName  string  `json:"templateName"`
		PatientName   string  `json:"patientName"`
		Gender        string  `json:"gender"`
		Age           string  `json:"age"`
		Department    string  `json:"department"`
		Condition     string  `json:"condition"`
		Treatment     string  `json:"treatment"`
		DiagnosisDate string  `json:"diagnosisDate"`
		Hospital      string  `json:"hospital"`
		Email         string  `json:"email"`
		FileName      string  `json:"fileName"`
		Price         float64 `json:"price"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	if req.TemplateName == "" || req.PatientName == "" || req.Email == "" {
		Fail(c, "请填写完整信息")
		return
	}
	if req.Price <= 0 {
		Fail(c, "模板价格无效")
		return
	}
	info := fmt.Sprintf("性别:%s|年龄:%s|科室:%s|病情:%s|治疗:%s|诊断日期:%s|医院:%s|邮箱:%s",
		req.Gender, req.Age, req.Department, req.Condition, req.Treatment, req.DiagnosisDate, req.Hospital, req.Email)
	data := []interface{}{
		gin.H{"userName": req.PatientName, "data": gin.H{"id": 0, "name": "病历提交"}},
	}
	order, msg := taowaCreateOrder(c, "bl", uid, req.TemplateName, req.PatientName, req.Email,
		info, req.Email, req.TemplateName, data, req.FileName, req.Price, "")
	if msg != "" {
		Fail(c, msg)
		return
	}
	OK(c, gin.H{"oid": order.OID, "msg": fmt.Sprintf("下单成功，扣费%s元！", order.Fees)})
}

// ---------- 订单管理 ----------

// GET /api/taowa/order/list?type=gz|bl&page=&pageSize=&status=
func TaowaOrderList(c *gin.Context) {
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(int64)
	typ := c.DefaultQuery("type", "gz")
	if typ != "gz" && typ != "bl" {
		typ = "gz"
	}
	page := parseInt64(c.DefaultQuery("page", "1"))
	pageSize := parseInt64(c.DefaultQuery("pageSize", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	q := database.DB.Model(&model.Order{}).Where("kcid = ?", typ)
	if uid != 1 {
		q = q.Where("uid = ?", uid)
	}
	if st := c.Query("status"); st != "" {
		q = q.Where("status = ?", st)
	}
	var total int64
	q.Count(&total)
	var list []model.Order
	q.Order("oid desc").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&list)
	PageResult(c, list, total, int(page), int(pageSize))
}

// POST /api/taowa/order/cancel 取消订单（本地，全额退费）
func TaowaOrderCancel(c *gin.Context) {
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(int64)
	oid := parseInt64(c.PostForm("oid"))
	var order model.Order
	q := database.DB.Where("oid = ? AND kcid IN ?", oid, []string{"gz", "bl"})
	if uid != 1 {
		q = q.Where("uid = ?", uid)
	}
	if err := q.First(&order).Error; err != nil {
		Fail(c, "订单不存在")
		return
	}
	if order.Status == "已取消" || order.Status == "已完成" {
		Fail(c, "当前状态不可取消")
		return
	}
	money := parsePrice(order.Fees)
	database.DB.Model(&model.User{}).Where("uid = ?", order.UID).Update("money", gorm.Expr("money + ?", money))
	database.DB.Model(&order).Update("status", "已取消")
	writeLog(order.UID, "taowa", fmt.Sprintf("取消订单 %d 退回%.2f", oid, money), money, c.ClientIP())
	OK(c, gin.H{"msg": fmt.Sprintf("订单已取消，退回%.2f元", money)})
}

// POST /api/taowa/order/remark 改备注（管理员或本人）
func TaowaOrderRemark(c *gin.Context) {
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(int64)
	oid := parseInt64(c.PostForm("oid"))
	remark := strings.TrimSpace(c.PostForm("remark"))
	var order model.Order
	q := database.DB.Where("oid = ? AND kcid IN ?", oid, []string{"gz", "bl"})
	if uid != 1 {
		q = q.Where("uid = ?", uid)
	}
	if err := q.First(&order).Error; err != nil {
		Fail(c, "订单不存在")
		return
	}
	database.DB.Model(&order).Update("remarks", remark)
	OK(c, nil)
}

// POST /api/taowa/order/gxx 补充文件
func TaowaOrderGxx(c *gin.Context) {
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(int64)
	oid := parseInt64(c.PostForm("oid"))
	fileName := strings.TrimSpace(c.PostForm("fileName"))
	var order model.Order
	q := database.DB.Where("oid = ? AND kcid IN ?", oid, []string{"gz", "bl"})
	if uid != 1 {
		q = q.Where("uid = ?", uid)
	}
	if err := q.First(&order).Error; err != nil {
		Fail(c, "订单不存在")
		return
	}
	old := order.CourseEndTime
	if old != "" {
		fileName = old + "," + fileName
	}
	database.DB.Model(&order).Update("courseEndTime", fileName)
	OK(c, gin.H{"msg": "补充成功"})
}

// POST /api/taowa/order/status 改状态（管理员）
func TaowaOrderStatus(c *gin.Context) {
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(int64)
	if uid != 1 {
		Fail(c, "权限不足")
		return
	}
	oid := parseInt64(c.PostForm("oid"))
	status := strings.TrimSpace(c.PostForm("status"))
	var order model.Order
	if err := database.DB.Where("oid = ? AND kcid IN ?", oid, []string{"gz", "bl"}).First(&order).Error; err != nil {
		Fail(c, "订单不存在")
		return
	}
	database.DB.Model(&order).Update("status", status)
	OK(c, nil)
}

// POST /api/taowa/order/tk 全额退款（管理员）
func TaowaOrderTk(c *gin.Context) {
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(int64)
	if uid != 1 {
		Fail(c, "权限不足")
		return
	}
	oid := parseInt64(c.PostForm("oid"))
	var order model.Order
	if err := database.DB.Where("oid = ? AND kcid IN ?", oid, []string{"gz", "bl"}).First(&order).Error; err != nil {
		Fail(c, "订单不存在")
		return
	}
	if order.Status == "已退款" {
		Fail(c, "订单已退款")
		return
	}
	money := parsePrice(order.Fees)
	database.DB.Model(&model.User{}).Where("uid = ?", order.UID).Update("money", gorm.Expr("money + ?", money))
	database.DB.Model(&order).Update("status", "已退款")
	writeLog(order.UID, "taowa", fmt.Sprintf("全额退款订单 %d 退回%.2f", oid, money), money, c.ClientIP())
	OK(c, gin.H{"msg": fmt.Sprintf("已全额退款%.2f元", money)})
}

// POST /api/taowa/order/bjtk 半价退款（管理员）
func TaowaOrderBjtk(c *gin.Context) {
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(int64)
	if uid != 1 {
		Fail(c, "权限不足")
		return
	}
	oid := parseInt64(c.PostForm("oid"))
	var order model.Order
	if err := database.DB.Where("oid = ? AND kcid IN ?", oid, []string{"gz", "bl"}).First(&order).Error; err != nil {
		Fail(c, "订单不存在")
		return
	}
	if order.Status == "已退款" {
		Fail(c, "订单已退款")
		return
	}
	money := parsePrice(order.Fees) / 2
	database.DB.Model(&model.User{}).Where("uid = ?", order.UID).Update("money", gorm.Expr("money + ?", money))
	database.DB.Model(&order).Update("status", "已退款")
	writeLog(order.UID, "taowa", fmt.Sprintf("半价退款订单 %d 退回%.2f", oid, money), money, c.ClientIP())
	OK(c, gin.H{"msg": fmt.Sprintf("已半价退款%.2f元", money)})
}

// POST /api/taowa/order/del 删除订单（管理员）
func TaowaOrderDel(c *gin.Context) {
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(int64)
	if uid != 1 {
		Fail(c, "权限不足")
		return
	}
	oid := parseInt64(c.PostForm("oid"))
	var order model.Order
	if err := database.DB.Where("oid = ? AND kcid IN ?", oid, []string{"gz", "bl"}).First(&order).Error; err != nil {
		Fail(c, "订单不存在")
		return
	}
	database.DB.Delete(&order)
	OK(c, gin.H{"msg": "已删除"})
}

// POST /api/admin/taowa/sync  自动从源台拉分类+公司/模板，创建本地分类和商品，并回填配置
func TaowaSync(c *gin.Context) {
	if !taowaCredsReady() {
		Fail(c, "请先填写源台 UID 和 KEY")
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")

	// 1) 拉盖章公司列表
	gzRes := taowaProxy("sxgetclass", gin.H{})
	gzBody, gzOk := taowaBody(gzRes)
	if !gzOk {
		msg, _ := gzRes["msg"].(string)
		Fail(c, "拉取盖章公司失败: " + msg)
		return
	}
	gzCompanies, _ := gzBody["data"].([]interface{})

	// 2) 创建/获取"实习盖章"分类
	gzFenleiName := "实习盖章"
	var gzFL model.FenLei
	if err := database.DB.Where("name = ?", gzFenleiName).First(&gzFL).Error; err != nil {
		gzFL = model.FenLei{Name: gzFenleiName, Sort: "0", Status: "1", Time: now}
		if err := database.DB.Create(&gzFL).Error; err != nil {
			Fail(c, "创建盖章分类失败: " + err.Error())
			return
		}
	}
	gzFLIDStr := strconv.FormatInt(gzFL.ID, 10)

	// 3) 批量插入盖章商品（docking=taowa_gz）
	gzCount := 0
	for _, item := range gzCompanies {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := m["name"].(string)
		if name == "" {
			continue
		}
		// 跳过已存在
		var cnt int64
		database.DB.Model(&model.Class{}).Where("docking = ? AND name = ?", "taowa_gz", name).Count(&cnt)
		if cnt > 0 {
			continue
		}
		price, _ := m["price"].(string)
		if price == "" {
			if f, ok := m["price"].(float64); ok {
				price = strconv.FormatFloat(f, 'f', 2, 64)
			}
		}
		content, _ := m["content"].(string)
		cls := model.Class{
			Sort:      10,
			Name:      name,
			GetNoun:   "taowa_gz",
			Noun:      "taowa_gz",
			Price:     price,
			QueryPlat: "taowa_gz",
			Docking:   "taowa_gz",
			YunSuan:   "*",
			Content:   content,
			AddTime:   now,
			Status:    1,
			FenLei:    gzFLIDStr,
		}
		if err := database.DB.Create(&cls).Error; err == nil {
			gzCount++
		}
	}

	// 4) 拉病历模板列表
	blRes := taowaProxy("getclassbl", gin.H{})
	blBody, blOk := taowaBody(blRes)
	blTemplates, _ := blBody["data"].([]interface{})
	blFLIDStr := ""
	blCount := 0
	if blOk {
		blFenleiName := "实习病历"
		var blFL model.FenLei
		if err := database.DB.Where("name = ?", blFenleiName).First(&blFL).Error; err != nil {
			blFL = model.FenLei{Name: blFenleiName, Sort: "0", Status: "1", Time: now}
			if err := database.DB.Create(&blFL).Error; err == nil {
				blFLIDStr = strconv.FormatInt(blFL.ID, 10)
			}
		} else {
			blFLIDStr = strconv.FormatInt(blFL.ID, 10)
		}
		for _, item := range blTemplates {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			name, _ := m["name"].(string)
			if name == "" {
				continue
			}
			var cnt int64
			database.DB.Model(&model.Class{}).Where("docking = ? AND name = ?", "taowa_bl", name).Count(&cnt)
			if cnt > 0 {
				continue
			}
			price, _ := m["price"].(string)
			if price == "" {
				if f, ok := m["price"].(float64); ok {
					price = strconv.FormatFloat(f, 'f', 2, 64)
				}
			}
			cls := model.Class{
				Sort:      10,
				Name:      name,
				GetNoun:   "taowa_bl",
				Noun:      "taowa_bl",
				Price:     price,
				QueryPlat: "taowa_bl",
				Docking:   "taowa_bl",
				YunSuan:   "*",
				Content:   "",
				AddTime:   now,
				Status:    1,
				FenLei:    blFLIDStr,
			}
			if err := database.DB.Create(&cls).Error; err == nil {
				blCount++
			}
		}
	}

	// 5) 回填配置项
	setConfig := func(k, v string) {
		database.DB.Model(&model.Config{}).Where("v = ?", k).Update("k", v)
	}
	setConfig("taowa_fenlei", gzFLIDStr)
	setConfig("taowa_blfenlei", blFLIDStr)
	// 平台ID = 货源表的 hid（qywk.top 那条记录）
	var huoYuan struct {
		HID int64 `gorm:"column:hid"`
	}
	database.DB.Model(&model.HuoYuan{}).Where("url LIKE ?", "%qywk.top%").First(&huoYuan)
	if huoYuan.HID > 0 {
		hidStr := strconv.FormatInt(huoYuan.HID, 10)
		setConfig("taowa_hid", hidStr)
		setConfig("taowa_blhid", hidStr)
	}

	OK(c, gin.H{
		"gz_fenlei_id": gzFLIDStr,
		"gz_count":     gzCount,
		"gz_total":     len(gzCompanies),
		"bl_fenlei_id": blFLIDStr,
		"bl_count":     blCount,
		"bl_total":     len(blTemplates),
		"huoyuan_hid":  huoYuan.HID,
	})
}
