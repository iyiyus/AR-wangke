package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"wk-go/internal/database"
	"wk-go/internal/model"
)

/* ============================================================
 * daka 实习打卡对接（对应 PHP 对接包 sxdk/api.php）
 * 配置项（config 表，管理员在系统设置填写）：
 *   daka_token      源台 token（默认 "您的token"）
 *   daka_admin      源台账号（默认 "您的TaiShan账号"）
 *   daka_ts_version 版本号（默认 260901）
 *   daka_url_list   源台地址 JSON 数组
 *   daka_del_return 删除订单是否退款（true/false，默认 true）
 * ============================================================ */

// ---------- 配置 ----------

func dakaConf(key, def string) string {
	var c model.Config
	if err := database.DB.Where("v = ?", key).First(&c).Error; err == nil && c.K != "" {
		return c.K
	}
	return def
}

func dakaURLList() []string {
	s := dakaConf("daka_url_list", "")
	if s == "" {
		return []string{
			"http://location.copilotai.top:4007/copilot/",
			"http://82.156.247.209:4007/copilot/",
			"http://location.tspost.top:4007/copilot/",
		}
	}
	var list []string
	if json.Unmarshal([]byte(s), &list) != nil || len(list) == 0 {
		return []string{"http://location.copilotai.top:4007/copilot/"}
	}
	return list
}

// 倍率前单价（元/天），各平台单独改
func dakaPlatformPrice(platform string) float64 {
	prices := map[string]float64{
		"zxjy":   0.6, // 职校家园
		"qzt":    0.6, // 黔职通
		"xyb":    0.6, // 校友帮
		"gxy":    0.6, // 工学云
		"xxy":    0.6, // 习讯云
		"xxt":    0.6, // 学习通
		"hzj":    0.6, // 慧职教
		"gxzy":   0.6, // 广西职业
		"jxzhjy": 0.6, // 江西智慧教育
		"cxy":    0.6, // 成学云
		"bx":     0.6, // 博行
	}
	if v, ok := prices[platform]; ok {
		return v
	}
	return 10 // 默认
}

// 前台计价：平台单价 × 当前登录用户 addprice
func dakaQPrice(platform string, addprice float64) float64 {
	return math.Round(addprice*dakaPlatformPrice(platform)*100) / 100
}

// 天数计算：结束日期距今天之间、打卡周期内的天数（移植 PHP timeCalcTrueday）
func dakaCalcDays(now time.Time, endTime, checkWeek string) int {
	if endTime == "" {
		return 0
	}
	// 打卡周期 0-6
	weekStrs := strings.Split(checkWeek, ",")
	weeks := []int{}
	for _, w := range weekStrs {
		if n, err := strconv.Atoi(strings.TrimSpace(w)); err == nil {
			weeks = append(weeks, n)
		}
	}
	sort.Ints(weeks)
	if len(weeks) == 0 {
		weeks = []int{0, 1, 2, 3, 4, 5, 6}
	}

	end, err := time.ParseInLocation("2006-01-02", endTime, time.Local)
	if err != nil {
		return 0
	}
	endSjc := end.Add(24*time.Hour - time.Second) // 结束日 23:59:59
	if endSjc.Before(now) {
		return 0
	}
	nowWeekDay := (int(now.Weekday()) + 6) % 7 // 周一=0 ... 周日=6
	// 本周周末 23:59:59
	weekEndSjc := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local).Add(time.Duration(6-nowWeekDay) * 24 * time.Hour)

	countLE := func(list []int, day int) int {
		c := 0
		for _, w := range list {
			if w <= day {
				c++
			}
		}
		return c
	}

	nowWeekLast := []int{}
	for _, w := range weeks {
		if w >= nowWeekDay {
			nowWeekLast = append(nowWeekLast, w)
		}
	}
	endWeekDay := (int(endSjc.Weekday()) + 6) % 7

	if !endSjc.After(weekEndSjc) {
		// 结束时间在本周内
		lastWeekLast := 0
		for _, w := range nowWeekLast {
			if w <= endWeekDay {
				lastWeekLast++
			}
		}
		return lastWeekLast
	}
	// 结束时间不在本周内
	endWeekLast := countLE(weeks, endWeekDay)
	intSjc := endSjc.Sub(weekEndSjc.Add(time.Duration(endWeekDay+1) * 24 * time.Hour))
	whole := int(intSjc/(7*24*time.Hour)) * len(weeks)
	return len(nowWeekLast) + whole + endWeekLast
}

// 本站业务日志
func dakaWlog(uid int64, title, msg string, money float64) {
	if uid <= 0 {
		return
	}
	writeLog(uid, "sxdk", title+"："+msg, money, "")
}

// ---------- 源台通信（移植 sendPostRequest）----------

func dakaPost(urlEnd string, data gin.H, uid int64) gin.H {
	if uid <= 0 {
		return gin.H{"code": 1, "msg": "未登录"}
	}
	admin := dakaConf("daka_admin", "")
	token := dakaConf("daka_token", "")
	if admin == "" || token == "" || admin == "您的TaiShan账号" || token == "您的token" {
		return gin.H{"code": 1, "msg": "源台未配置，请在系统设置中填写账号与Token"}
	}
	payload := gin.H{
		"admin":      admin,
		"token":      token,
		"ts_version": dakaConf("daka_ts_version", "260901"),
	}
	for k, v := range data {
		payload[k] = v
	}
	body, _ := json.Marshal(payload)
	urls := dakaURLList()
	client := &http.Client{Timeout: 120 * time.Second}

	for i := 0; i < len(urls)*3; i++ {
		url := urls[i/3] + urlEnd
		req, err := http.NewRequest("POST", url, bytes.NewReader(body))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		resp, err := client.Do(req)
		if err != nil {
			dakaWlog(uid, "TaiShan-网络异常", safeURL(url), 0)
			continue
		}
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if len(respBody) > 0 {
			var decoded gin.H
			if json.Unmarshal(respBody, &decoded) == nil {
				return decoded
			}
		}
		dakaWlog(uid, "TaiShan-网络异常", safeURL(url), 0)
	}
	// addOrder 网络异常回查
	if strings.Contains(urlEnd, "addOrder") {
		dakaWlog(uid, "TaiShan-下单网络异常回查", fmt.Sprintf("%v %v", data["platform"], data["phone"]), 0)
		for i := 0; i < 3; i++ {
			if i > 0 {
				time.Sleep(time.Second)
			}
			sel := dakaPost("selectOrderById", data, uid)
			if code, _ := sel["code"].(float64); code == 0 {
				if arr, ok := sel["data"].([]interface{}); ok && len(arr) == 1 {
					row, _ := arr[0].(map[string]interface{})
					ct, _ := row["createTime"].(string)
					if ctT, err := time.ParseInLocation("2006-01-02 15:04:05", ct, time.Local); err == nil {
						diff := math.Abs(time.Since(ctT).Seconds())
						if diff <= 300 {
							dakaWlog(uid, "TaiShan-下单网络异常判定成功", fmt.Sprintf("%v %v 时间差%.0f秒", data["platform"], data["phone"], diff), 0)
							return gin.H{"code": 0, "msg": "网络异常，判定为添加成功", "selectOrderById": sel}
						}
					}
				}
			}
		}
		dakaWlog(uid, "TaiShan-下单网络异常判定失败", fmt.Sprintf("%v %v 回查未命中", data["platform"], data["phone"]), 0)
		return gin.H{"code": 1, "msg": "添加订单网络异常，下单失败"}
	}
	return gin.H{"code": 1, "msg": "网络异常，源台无法连接，请反馈给管理员"}
}

func safeURL(url string) string {
	if idx := strings.Index(url, "?"); idx >= 0 {
		return url[:idx]
	}
	return url
}

// wxpush 字段处理
func dakaProcessWxpush(wxpush string) gin.H {
	if wxpush == "" {
		return gin.H{"wxpush": nil}
	}
	var decoded gin.H
	if json.Unmarshal([]byte(wxpush), &decoded) == nil {
		return decoded
	}
	return gin.H{"wxpush": wxpush}
}

// 判断当前用户是否可操作该订单（uid=1 站长可操作所有）
func dakaOwnOrder(uid int64, id int64) (model.SxdkOrder, bool) {
	var o model.SxdkOrder
	q := database.DB.Where("id = ?", id)
	if uid != 1 {
		q = q.Where("uid = ?", uid)
	}
	if err := q.First(&o).Error; err != nil {
		return o, false
	}
	return o, true
}

// ---------- 入口 ----------

// DakaAPI 实习打卡统一入口 GET /api/daka/api
func DakaAPI(c *gin.Context) {
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(int64)
	if uid <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "未登录"})
		return
	}
	var user model.User
	if err := database.DB.First(&user, uid).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户不存在"})
		return
	}
	act := c.Query("act")
	if act == "" {
		act = c.PostForm("act")
	}
	form := gin.H{}
	// 支持 form 与 JSON body
	if c.Request.Method == http.MethodPost {
		ct := c.GetHeader("Content-Type")
		if strings.Contains(ct, "application/json") {
			var body gin.H
			if err := c.ShouldBindJSON(&body); err == nil {
				form = body
			}
		} else {
			c.Request.ParseForm()
			for k, v := range c.Request.PostForm {
				form[k] = v[0]
			}
		}
	} else {
		for k, v := range c.Request.URL.Query() {
			form[k] = v[0]
		}
	}
	dakaHandle(c, act, form, user)
}

func dakaHandle(c *gin.Context, act string, form gin.H, user model.User) {
	uid := user.UID
	s := func(k string) string { v, _ := form[k].(string); return strings.TrimSpace(v) }

	switch act {
	case "price":
		// 各平台今日单价
		platforms := []string{"zxjy", "qzt", "xyb", "gxy", "xxy", "xxt", "hzj", "gxzy", "jxzhjy", "cxy", "bx"}
		names := map[string]string{"zxjy": "职校家园", "qzt": "黔职通", "xyb": "校友帮", "gxy": "工学云", "xxy": "习讯云", "xxt": "学习通", "hzj": "慧职教", "gxzy": "广西职业", "jxzhjy": "江西智慧教育", "cxy": "成学云", "bx": "博行"}
		data := gin.H{}
		for _, p := range platforms {
			data[p] = gin.H{"name": names[p], "price": dakaQPrice(p, user.AddPrice), "yuan": dakaPlatformPrice(p)}
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": data})
	case "getNotice":
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": dakaConf("tcgonggao", "")})
	case "order":
		// 我的订单列表 + 平台价格/剩余天数
		var list []model.SxdkOrder
		q := database.DB.Where("uid = ?", uid)
		if s("id") != "" {
			q = q.Where("id = ?", parseInt64(s("id")))
		}
		if s("platform") != "" {
			q = q.Where("platform = ?", s("platform"))
		}
		q.Order("id desc").Find(&list)
		type row struct {
			model.SxdkOrder
			Price   float64 `json:"price"`
			Day     int     `json:"day"`
			Expired bool    `json:"expired"`
		}
		out := []row{}
		for _, o := range list {
			day := dakaCalcDays(time.Now(), o.EndTime, o.CheckWeek)
			expired := false
			if et, err := time.ParseInLocation("2006-01-02", o.EndTime, time.Local); err == nil {
				expired = et.Add(24*time.Hour - time.Second).Before(time.Now())
			}
			out = append(out, row{SxdkOrder: o, Price: dakaQPrice(o.Platform, user.AddPrice), Day: day, Expired: expired})
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": out})
	case "add":
		formD := gin.H{}
		if f, ok := form["form"].(map[string]interface{}); ok {
			formD = f
		} else if f, ok := form["form"].(string); ok {
			json.Unmarshal([]byte(f), &formD)
		}
		if len(formD) == 0 {
			formD = form
		}
		platform := toStr(formD["platform"])
		phone := toStr(formD["phone"])
		password := toStr(formD["password"])
		name := toStr(formD["name"])
		address := toStr(formD["address"])
		checkTime := toStr(formD["check_time"])
		upCheckTime := toStr(formD["up_check_time"])
		downCheckTime := toStr(formD["down_check_time"])
		checkWeek := toStr(formD["check_week"])
		endTime := toStr(formD["end_time"])
		dayPaper := toInt(formD["day_paper"])
		weekPaper := toInt(formD["week_paper"])
		monthPaper := toInt(formD["month_paper"])
		runType := toInt(formD["runType"])

		day := dakaCalcDays(time.Now(), endTime, checkWeek)
		bei := 1
		if platform == "xyb" && runType == 3 {
			bei = 5
		}
		// 已存在检查
		var cnt int64
		database.DB.Model(&model.SxdkOrder{}).Where("uid = ? AND phone = ? AND platform = ?", uid, phone, platform).Count(&cnt)
		if cnt > 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "订单已存在"})
			return
		}
		et, err := time.ParseInLocation("2006-01-02", endTime, time.Local)
		if err != nil || !et.Add(24*time.Hour-time.Second).After(time.Now()) {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "下单天数不符合规范"})
			return
		}
		money := dakaQPrice(platform, user.AddPrice) * float64(day) * float64(bei)
		if money < 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "金额异常"})
			return
		}
		if user.Money < math.Round(money) {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "余额不足"})
			return
		}
		// 源台下单
		res := dakaPost("addOrder", gin.H{
			"platform": platform, "phone": phone, "password": password, "name": name,
			"address": address, "check_time": checkTime, "up_check_time": upCheckTime,
			"down_check_time": downCheckTime, "check_week": checkWeek, "end_time": endTime,
			"day_paper": dayPaper, "week_paper": weekPaper, "month_paper": monthPaper,
			"runType": runType,
		}, uid)
		if code, _ := res["code"].(float64); code != 0 {
			c.JSON(http.StatusOK, res)
			return
		}
		sel, _ := res["selectOrderById"].(gin.H)
		var upID int64
		if arr, ok := sel["data"].([]interface{}); ok && len(arr) == 1 {
			if row0, ok := arr[0].(map[string]interface{}); ok {
				upID = int64(toFloat(row0["id"]))
			}
		}
		if upID <= 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "下单失败，请联系管理员"})
			return
		}
		if upCheckTime == "" {
			upCheckTime = checkTime
		}
		wxpush := gin.H{"wxpush": ""}
		if platform == "xyb" {
			wxpush["runType"] = runType
		}
		wxpushJSON, _ := json.Marshal(wxpush)
		now := time.Now()
		nw := model.SxdkOrder{
			SxdkID: upID, UID: uid, Platform: platform, Phone: phone, Password: password,
			Code: 1, WxPush: string(wxpushJSON), Name: name, Address: address,
			UpCheckTime: upCheckTime, DownCheck: downCheckTime, CheckWeek: checkWeek,
			EndTime: endTime, DayPaper: dayPaper, WeekPaper: weekPaper, MonthPaper: monthPaper,
			CreateTime: now.Format("2006-01-02 15:04:05"), UpdateTime: now.Format("2006-01-02 15:04:05"),
		}
		if err := database.DB.Create(&nw).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "本地记录失败"})
			return
		}
		database.DB.Model(&user).Update("money", user.Money-money)
		dakaWlog(uid, "TaiShan-本台添加成功", fmt.Sprintf("%s %s 对接id：%d 天数：%d 结束：%s 扣除%.2f", platform, phone, upID, day, endTime, money), -money)
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": fmt.Sprintf("订单添加成功，扣除%.2f元！", money)})
	case "searchPhoneInfo":
		res := dakaPost("searchPhoneInfo", gin.H{
			"platform": s("platform"), "phone": s("phone"), "password": s("password"),
		}, uid)
		c.JSON(http.StatusOK, res)
	case "del":
		id := parseInt64(s("id"))
		o, ok := dakaOwnOrder(uid, id)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "您无此订单"})
			return
		}
		res := dakaPost("deleteOrder", gin.H{"id": o.SxdkID, "platform": o.Platform, "check_week": o.CheckWeek, "end_time": o.EndTime, "wxpush": o.WxPush}, uid)
		if code, _ := res["code"].(float64); code == 0 {
			otherMsg := ""
			refund := 0.0
			day := dakaCalcDays(time.Now(), o.EndTime, o.CheckWeek)
			if dakaConf("daka_del_return", "true") == "true" {
				wx := dakaProcessWxpush(o.WxPush)
				bei := 1
				if o.Platform == "xyb" {
					if rt, ok := wx["runType"].(float64); ok && rt == 3 {
						bei = 5
					}
				}
				refund = dakaQPrice(o.Platform, user.AddPrice) * float64(day) * float64(bei)
				if refund > 0 {
					if et, err := time.ParseInLocation("2006-01-02", o.EndTime, time.Local); err == nil && et.Add(24*time.Hour-time.Second).After(time.Now()) {
						otherMsg = fmt.Sprintf("，订单未到期，已退款：%.2f", refund)
						database.DB.Model(&user).Update("money", user.Money+refund)
					} else {
						otherMsg = "，此订单已到期，无需退款"
						refund = 0
					}
				}
			}
			database.DB.Delete(&o)
			dakaWlog(uid, "TaiShan-删单成功", fmt.Sprintf("本台id：%d 对接id：%d %s 结束：%s 剩余天数：%d%s", id, o.SxdkID, o.Platform, o.EndTime, day, otherMsg), refund)
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功" + otherMsg})
		} else {
			msg, _ := res["msg"].(string)
			dakaWlog(uid, "TaiShan-删单失败", fmt.Sprintf("本台id：%d 原因：%s", id, msg), 0)
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "删除失败，请联系管理员"})
		}
	case "getLog":
		id := parseInt64(s("id"))
		o, ok := dakaOwnOrder(uid, id)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "您无此订单"})
			return
		}
		c.JSON(http.StatusOK, dakaPost("getLog", gin.H{"phone": o.Phone, "platform": o.Platform}, uid))
	case "nowCheck":
		id := parseInt64(s("id"))
		platform := s("platform")
		money := dakaQPrice(platform, user.AddPrice)
		if user.Money < math.Round(money) {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "余额不足"})
			return
		}
		o, ok := dakaOwnOrder(uid, id)
		if !ok {
			if uid == 1 {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "立即打卡涉及扣费，站长无法操作代理订单"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "您无此订单"})
			return
		}
		res := dakaPost("nowCheck", gin.H{"id": o.SxdkID, "platform": o.Platform}, uid)
		if code, _ := res["code"].(float64); code == 0 {
			database.DB.Model(&user).Update("money", user.Money-money)
			dakaWlog(uid, "TaiShan-立即打卡成功", fmt.Sprintf("本台id：%d 平台：%s 对接id：%d 扣除%.2f", id, platform, o.SxdkID, money), -money)
		} else {
			msg, _ := res["msg"].(string)
			dakaWlog(uid, "TaiShan-立即打卡失败", fmt.Sprintf("本台id：%d 原因：%s", id, msg), 0)
		}
		c.JSON(http.StatusOK, res)
	case "buPapers":
		id := parseInt64(s("id"))
		startTime := s("startTime")
		endTime := s("endTime")
		levelName := s("levelName")
		o, ok := dakaOwnOrder(uid, id)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "您无此订单"})
			return
		}
		res := dakaPost("buPapers", gin.H{"id": o.SxdkID, "platform": o.Platform, "startTime": startTime, "endTime": endTime, "type": levelName}, uid)
		if code, _ := res["code"].(float64); code == 0 {
			dakaWlog(uid, "TaiShan-补报告成功", fmt.Sprintf("本台id：%d 类型：%s 区间：%s-%s", id, levelName, startTime, endTime), 0)
		} else {
			msg, _ := res["msg"].(string)
			dakaWlog(uid, "TaiShan-补报告失败", fmt.Sprintf("本台id：%d 类型：%s 原因：%s", id, levelName, msg), 0)
		}
		c.JSON(http.StatusOK, res)
	case "changeCheckCode":
		id := parseInt64(s("id"))
		code := s("code")
		o, ok := dakaOwnOrder(uid, id)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "您无此订单"})
			return
		}
		res := dakaPost("setCheckCode", gin.H{"id": o.SxdkID, "platform": o.Platform, "code": code}, uid)
		if codef, _ := res["code"].(float64); codef == 0 {
			database.DB.Model(&o).Update("code", atoiSafeInt(code))
			dakaWlog(uid, "TaiShan-改状态成功", fmt.Sprintf("本台id：%d code：%s", id, code), 0)
		} else {
			msg, _ := res["msg"].(string)
			dakaWlog(uid, "TaiShan-改状态失败", fmt.Sprintf("本台id：%d 原因：%s", id, msg), 0)
		}
		c.JSON(http.StatusOK, res)
	case "changeHolidayCode":
		id := parseInt64(s("id"))
		code := s("code")
		o, ok := dakaOwnOrder(uid, id)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "您无此订单"})
			return
		}
		res := dakaPost("setHolidayCode", gin.H{"id": o.SxdkID, "platform": o.Platform, "code": code}, uid)
		if codef, _ := res["code"].(float64); codef == 0 {
			dakaWlog(uid, "TaiShan-改节假日成功", fmt.Sprintf("本台id：%d code：%s", id, code), 0)
		} else {
			msg, _ := res["msg"].(string)
			dakaWlog(uid, "TaiShan-改节假日失败", fmt.Sprintf("本台id：%d 原因：%s", id, msg), 0)
		}
		c.JSON(http.StatusOK, res)
	case "getWxPush":
		id := parseInt64(s("id"))
		o, ok := dakaOwnOrder(uid, id)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "您无此订单"})
			return
		}
		c.JSON(http.StatusOK, dakaPost("getWxPush", gin.H{"phone": o.Phone, "platform": o.Platform}, uid))
	case "querySourceOrder":
		id := parseInt64(s("id"))
		o, ok := dakaOwnOrder(uid, id)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "您无此订单"})
			return
		}
		res := dakaPost("selectOrderById", gin.H{"id": o.SxdkID, "platform": o.Platform, "phone": o.Phone}, uid)
		if code, _ := res["code"].(float64); code != 0 {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "订单不存在，请联系管理员"})
			return
		}
		if arr, ok := res["data"].([]interface{}); ok && len(arr) == 1 {
			c.JSON(http.StatusOK, gin.H{"code": 0, "data": arr[0]})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "订单不存在，请联系管理员"})
	case "edit":
		formD := gin.H{}
		if f, ok := form["form"].(map[string]interface{}); ok {
			formD = f
		} else if f, ok := form["form"].(string); ok {
			json.Unmarshal([]byte(f), &formD)
		}
		if len(formD) == 0 {
			formD = form
		}
		id := parseInt64(toStr(formD["id"]))
		platform := toStr(formD["platform"])
		phone := toStr(formD["phone"])
		password := toStr(formD["password"])
		name := toStr(formD["name"])
		address := toStr(formD["address"])
		checkTime := toStr(formD["check_time"])
		upCheckTime := toStr(formD["up_check_time"])
		downCheckTime := toStr(formD["down_check_time"])
		checkWeek := toStr(formD["check_week"])
		endTime := toStr(formD["end_time"])
		dayPaper := toInt(formD["day_paper"])
		weekPaper := toInt(formD["week_paper"])
		monthPaper := toInt(formD["month_paper"])
		runType := toInt(formD["runType"])

		var o model.SxdkOrder
		if err := database.DB.Where("id = ? AND uid = ? AND phone = ?", id, uid, phone).First(&o).Error; err != nil {
			if uid == 1 {
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "改单涉及扣费，站长无法操作代理订单"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "您无此订单"})
			return
		}
		oldEnd, _ := time.ParseInLocation("2006-01-02", o.EndTime, time.Local)
		newEnd, _ := time.ParseInLocation("2006-01-02", endTime, time.Local)
		now := time.Now()
		var day int
		if !oldEnd.Add(24*time.Hour-time.Second).After(now) && newEnd.Add(24*time.Hour-time.Second).After(now) {
			day = dakaCalcDays(now, endTime, checkWeek)
		} else {
			if !newEnd.Add(24*time.Hour-time.Second).After(oldEnd.Add(24*time.Hour-time.Second)) && o.CheckWeek == checkWeek {
				day = 0
			} else {
				oldDay := dakaCalcDays(now, o.EndTime, o.CheckWeek)
				newDay := dakaCalcDays(now, endTime, checkWeek)
				day = newDay - oldDay
				if day < 0 {
					day = 0
				}
			}
		}
		bei := 1
		if platform == "xyb" && runType == 3 {
			bei = 5
		}
		money := dakaQPrice(platform, user.AddPrice) * float64(day) * float64(bei)
		if user.Money < math.Round(money) {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "余额不足"})
			return
		}
		formD["id"] = o.SxdkID
		res := dakaPost("editOrder", formD, uid)
		if code, _ := res["code"].(float64); code != 0 {
			msg, _ := res["msg"].(string)
			dakaWlog(uid, "TaiShan-改单失败", fmt.Sprintf("本台id：%d 对接id：%d 原因：%s", id, o.SxdkID, msg), 0)
			c.JSON(http.StatusOK, res)
			return
		}
		if upCheckTime == "" {
			upCheckTime = checkTime
		}
		wx := dakaProcessWxpush(o.WxPush)
		if platform == "xyb" {
			wx["runType"] = runType
		}
		wxJSON, _ := json.Marshal(wx)
		nowStrV := now.Format("2006-01-02 15:04:05")
		database.DB.Model(&o).Updates(map[string]interface{}{
			"password": password, "name": name, "address": address,
			"up_check_time": upCheckTime, "down_check_time": downCheckTime,
			"check_week": checkWeek, "end_time": endTime, "wxpush": string(wxJSON),
			"day_paper": dayPaper, "week_paper": weekPaper, "month_paper": monthPaper,
			"updateTime": nowStrV,
		})
		database.DB.Model(&user).Update("money", user.Money-money)
		dakaWlog(uid, "TaiShan-改单成功", fmt.Sprintf("%s %s 本台id：%d 对接id：%d 增加天数：%d 原结束：%s 现结束：%s 扣除%.2f", platform, phone, id, o.SxdkID, day, o.EndTime, endTime, money), -money)
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": fmt.Sprintf("订单修改成功,扣费：%.2f", money), "msgs": res["msg"]})
	case "xxtGetSchoolList":
		c.JSON(http.StatusOK, dakaPost("xxtGetSchoolList", gin.H{"filter": s("filter")}, uid))
	case "hzjGetSchoolList":
		client := &http.Client{Timeout: 10 * time.Second}
		req, _ := http.NewRequest("POST", "https://hzj.gzdekan.com/api/getSchoolListWeb?authorization=", nil)
		req.Header.Set("User-Agent", "MyCustomUserAgent/1.0")
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "获取失败"})
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		c.Data(http.StatusOK, "application/json; charset=utf-8", body)
	case "xxyGetSchoolList":
		// 官方接口 + 源台回退
		c.JSON(http.StatusOK, dakaPost("xxyGetSchoolList", gin.H{}, uid))
	case "get_userrow":
		if uid == 1 {
			res := dakaPost("get_userrow", gin.H{}, uid)
			if code, _ := res["code"].(float64); code == 0 {
				c.JSON(http.StatusOK, gin.H{"code": 0, "data": res["data"]})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 1, "data": gin.H{"msg": ""}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "权限不足"})
	case "yunOrder":
		if uid == 1 {
			res := dakaPost("yunOrder", gin.H{}, uid)
			if code, _ := res["code"].(float64); code == 0 {
				fixNum := 0
				if arr, ok := res["data"].([]interface{}); ok {
					for _, it := range arr {
						row, _ := it.(map[string]interface{})
						upID := int64(toFloat(row["id"]))
						platform := toStr(row["platform"])
						code := int(toFloat(row["code"]))
						wxpush := toStr(row["wxpush"])
						endTime := toStr(row["end_time"])
						phone := toStr(row["phone"])
						password := toStr(row["password"])
						var exist model.SxdkOrder
						if database.DB.Where("sxdkId = ? AND platform = ?", upID, platform).First(&exist).Error == nil {
							database.DB.Model(&exist).Updates(map[string]interface{}{"code": code, "wxpush": wxpush, "end_time": endTime})
						}
						// 修补缺失 sxdkId 订单
						var miss model.SxdkOrder
						if database.DB.Where("phone = ? AND password = ? AND platform = ? AND sxdkId = 0", phone, password, platform).First(&miss).Error == nil {
							database.DB.Model(&miss).Updates(map[string]interface{}{"sxdkId": upID, "code": code, "wxpush": wxpush, "end_time": endTime})
							fixNum++
						}
					}
					countNum := len(arr)
					c.JSON(http.StatusOK, gin.H{"code": 0, "msg": fmt.Sprintf("拉取完成！同步：'%d'条成功，修复订单：'%d'条", countNum, fixNum), "data": res["data"]})
					return
				}
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "拉取完成！", "data": res["data"]})
				return
			}
			msg, _ := res["msg"].(string)
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": fmt.Sprintf("拉取失败：'%s'", msg)})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "权限不足"})
	default:
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "未知操作"})
	}
}

// ---------- 辅助 ----------

func toStr(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", t)
	}
}

func toFloat(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	case int:
		return float64(t)
	default:
		return 0
	}
}

func toInt(v interface{}) int {
	return int(toFloat(v))
}

func atoiSafeInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
