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

// 对外开放 API（兼容原 PHP api.php）
// GET/POST /open/api?act=xxx

func OpenAPI(c *gin.Context) {
	act := c.Query("act")
	switch act {
	case "getmoney":
		openGetMoney(c)
	case "get":
		openGet(c)
	case "add":
		openAdd(c)
	case "getadd":
		openGetAdd(c)
	case "chadan":
		openChaDan(c)
	case "budan":
		openBuDan(c)
	case "getclass":
		openGetClass(c)
	default:
		c.JSON(200, gin.H{"code": 0, "msg": "未知操作"})
	}
}

// openAuth 校验 uid+key，返回用户；失败时已输出响应
func openAuth(c *gin.Context) *model.User {
	uid := c.PostForm("uid")
	key := c.PostForm("key")
	if uid == "" || key == "" {
		c.JSON(200, gin.H{"code": 0, "msg": "所有项目不能为空"})
		return nil
	}
	var user model.User
	if err := database.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": "用户不存在"})
		return nil
	}
	if user.Key == "0" {
		c.JSON(200, gin.H{"code": -1, "msg": "你还没有开通接口哦"})
		return nil
	}
	if user.Key != key {
		c.JSON(200, gin.H{"code": -2, "msg": "密匙错误"})
		return nil
	}
	return &user
}

func openGetMoney(c *gin.Context) {
	user := openAuth(c)
	if user == nil {
		return
	}
	c.JSON(200, gin.H{"code": 1, "msg": "查询成功", "money": user.Money})
}

func openGet(c *gin.Context) {
	uid := c.PostForm("uid")
	platform := c.PostForm("platform")
	school := c.PostForm("school")
	user := c.PostForm("user")
	pass := c.PostForm("pass")
	typ := c.PostForm("type")

	if platform == "" || school == "" || user == "" || pass == "" {
		c.JSON(200, gin.H{"code": 0, "msg": "所有项目不能为空"})
		return
	}
	u := openAuth(c)
	if u == nil {
		return
	}
	if u.Money <= 10 {
		c.JSON(200, gin.H{"code": -2, "msg": "余额小于10禁止调用查课"})
		return
	}

	var cls model.Class
	if err := database.DB.Where("cid = ?", platform).First(&cls).Error; err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": "平台不存在"})
		return
	}
	if cls.Status == 0 {
		c.JSON(200, gin.H{"code": -2, "msg": "网课已下架禁止查课！"})
		return
	}

	res := checkorder.QueryByClass(cls.CID, school, user, pass)
	if res == nil {
		res = &checkorder.QueryResult{Code: -1, Msg: "查询失败,请联系管理员"}
	}
	out := gin.H{"code": res.Code, "msg": res.Msg, "data": res.Data, "userName": res.UserName, "userinfo": school + " " + user + " " + pass}
	if res.Code == 1 {
		writeLog(u.UID, "API查课", fmt.Sprintf("%s-查课信息：%s %s %s", cls.Name, school, user, pass), 0, c.ClientIP())
		// 更新查课次数与下单率
		var dd int64
		database.DB.Model(&model.Order{}).Where("uid = ?", u.UID).Count(&dd)
		database.DB.Model(&model.User{}).Where("uid = ?", uid).
			Updates(map[string]interface{}{
				"ck":   gorm.Expr("ck + 1"),
				"dd":   dd,
				"xdlv": gorm.Expr("100 * ? / (ck + 1)", dd),
			})
		// xiaochu 紧凑格式
		if typ == "xiaochu" {
			names := make([]string, 0, len(res.Data))
			for _, course := range res.Data {
				names = append(names, course.Name)
			}
			c.JSON(200, gin.H{
				"code": res.Code, "msg": res.Msg,
				"data": []string{cls.Name, user, pass, school, strings.Join(names, ",")},
				"js":   "", "info": "昔日之苦，安知异日不在尝之? 共勉",
			})
			return
		}
	}
	c.JSON(200, out)
}

func openAdd(c *gin.Context) {
	platform := c.PostForm("platform")
	school := c.PostForm("school")
	userAcc := c.PostForm("user")
	pass := c.PostForm("pass")
	kcid := c.PostForm("kcid")
	kcname := c.PostForm("kcname")

	if platform == "" || school == "" || userAcc == "" || pass == "" || kcname == "" {
		c.JSON(200, gin.H{"code": 0, "msg": "所有项目不能为空"})
		return
	}

	u := openAuth(c)
	if u == nil {
		return
	}

	var cls model.Class
	if err := database.DB.Where("cid = ?", platform).First(&cls).Error; err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": "平台不存在"})
		return
	}
	if cls.Status == 0 {
		c.JSON(200, gin.H{"code": -2, "msg": "小老弟，商品都下架了你还下什么单呢！"})
		return
	}

	danjia := calcPrice(cls, *u)
	if danjia <= 0 || u.AddPrice < 0.1 {
		c.JSON(200, gin.H{"code": -1, "msg": "大佬，我得罪不起您，我小本生意，有哪里得罪之处，还望多多包涵"})
		return
	}

	// wkm4 平台需要先查课并完整匹配课程名
	if pt := getHuoYuanPT(cls.Docking); pt == "wkm4" {
		m4 := checkorder.QueryByClass(cls.CID, school, userAcc, pass)
		if m4 != nil && m4.Code == 1 {
			found := false
			for _, course := range m4.Data {
				if course.Name == kcname {
					found = true
					if course.ID != "" {
						kcid = course.ID
					}
					break
				}
			}
			if !found {
				c.JSON(200, gin.H{"code": -1, "msg": "请完整输入课程名字", "status": -1, "message": "请完整输入课程名字"})
				return
			}
		}
	}

	kcNames := strings.Split(kcname, ",")
	kcIDs := strings.Split(kcid, ",")
	if u.Money < danjia*float64(len(kcNames)) {
		c.JSON(200, gin.H{"code": -1, "msg": "余额不足以本次提交"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	ok := false
	for i, name := range kcNames {
		kc := ""
		if i < len(kcIDs) {
			kc = kcIDs[i]
		}
		dockStatus := "0"
		if cls.Docking == "0" {
			dockStatus = "99"
		}
		// 重复下单（兼容 PHP：仍入库扣费，标记 dockstatus=3）
		var dup int64
		database.DB.Model(&model.Order{}).Where(
			"ptname=? AND school=? AND user=? AND pass=? AND kcid=? AND kcname=?",
			cls.Name, school, userAcc, pass, kcid, kcname,
		).Count(&dup)
		if dup > 0 {
			dockStatus = "3"
		}

		order := model.Order{
			UID: u.UID, CID: cls.CID,
			HID: func() int64 { v, _ := strconv.ParseInt(cls.Docking, 10, 64); return v }(),
			PTName: cls.Name, School: school, UserAccount: userAcc, Pass: pass,
			KCID: kc, KCName: name,
			Fees: fmt.Sprintf("%.2f", danjia), Noun: cls.Noun,
			AddTime: now, IP: c.ClientIP(), DockStatus: dockStatus, Status: "待处理",
		}
		if err := database.DB.Create(&order).Error; err == nil {
			database.DB.Model(&model.User{}).Where("uid = ?", u.UID).
				UpdateColumn("money", gorm.Expr("money - ?", danjia))
			writeLog(u.UID, "API添加任务", fmt.Sprintf("%s %s %s 扣除%.2f元！", userAcc, pass, name, danjia), -danjia, c.ClientIP())
			ok = true
		}
	}
	if ok {
		c.JSON(200, gin.H{"code": 0, "msg": "提交成功", "status": 0, "message": "提交成功", "id": "订单号登录后台自行查看"})
		return
	}
	c.JSON(200, gin.H{"code": -1, "msg": "请完整输入课程名字", "status": -1, "message": "请完整输入课程名字"})
}

// openGetAdd 查课后按相似度>90% 自动下单
func openGetAdd(c *gin.Context) {
	platform := c.PostForm("platform")
	school := c.PostForm("school")
	userAcc := c.PostForm("user")
	pass := c.PostForm("pass")
	kcname := c.PostForm("kcname")

	if platform == "" || school == "" || userAcc == "" || pass == "" || kcname == "" {
		c.JSON(200, gin.H{"code": 0, "msg": "所有项目不能为空"})
		return
	}
	u := openAuth(c)
	if u == nil {
		return
	}

	var cls model.Class
	if err := database.DB.Where("cid = ?", platform).First(&cls).Error; err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": "平台不存在"})
		return
	}

	danjia := calcPrice(cls, *u)
	if danjia <= 0 || u.AddPrice < 0.1 {
		c.JSON(200, gin.H{"code": -1, "msg": "大佬，我得罪不起您，我小本生意，有哪里得罪之处，还望多多包涵"})
		return
	}
	if u.Money < danjia {
		c.JSON(200, gin.H{"code": -1, "msg": "余额不足"})
		return
	}

	res := checkorder.QueryByClass(cls.CID, school, userAcc, pass)
	if res == nil || res.Code != 1 {
		msg := "查询失败"
		if res != nil {
			msg = res.Msg
		}
		c.JSON(200, gin.H{"code": -1, "msg": msg})
		return
	}

	for _, course := range res.Data {
		if similarText(course.Name, kcname) > 90 {
			dockStatus := "0"
			if cls.Docking == "0" {
				dockStatus = "99"
			}
			kc := course.Name
			order := model.Order{
				UID: u.UID, CID: cls.CID,
				HID: func() int64 { v, _ := strconv.ParseInt(cls.Docking, 10, 64); return v }(),
				PTName: cls.Name, School: school, UserAccount: userAcc, Pass: pass,
				KCID: course.ID, KCName: kc,
				Fees: fmt.Sprintf("%.2f", danjia), Noun: cls.Noun,
				AddTime: time.Now().Format("2006-01-02 15:04:05"), IP: c.ClientIP(),
				DockStatus: dockStatus, Status: "待处理",
			}
			if err := database.DB.Create(&order).Error; err == nil {
				database.DB.Model(&model.User{}).Where("uid = ?", u.UID).
					UpdateColumn("money", gorm.Expr("money - ?", danjia))
				writeLog(u.UID, "API添加任务", fmt.Sprintf("%s %s %s 扣除%.2f元！", userAcc, pass, kcname, danjia), -danjia, c.ClientIP())
				c.JSON(200, gin.H{"code": 0, "msg": "提交成功", "status": 0, "message": "提交成功", "id": "订单号登录后台自行查看"})
				return
			}
		}
	}
	c.JSON(200, gin.H{"code": -1, "msg": "请完整输入课程名字", "status": -1, "message": "请完整输入课程名字"})
}

// getHuoYuanPT 获取对接货源的平台代码
func getHuoYuanPT(docking string) string {
	hid, err := strconv.ParseInt(docking, 10, 64)
	if err != nil || hid == 0 {
		return ""
	}
	var huo model.HuoYuan
	if err := database.DB.Where("hid = ?", hid).First(&huo).Error; err != nil {
		return ""
	}
	return huo.PT
}

// similarText 仿 PHP similar_text 的相似度百分比
func similarText(a, b string) float64 {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	if la == 0 && lb == 0 {
		return 100
	}
	if la == 0 || lb == 0 {
		return 0
	}
	// LCS 长度
	dp := make([][]int, 2)
	dp[0] = make([]int, lb+1)
	dp[1] = make([]int, lb+1)
	for i := 1; i <= la; i++ {
		row := dp[i%2]
		prev := dp[(i-1)%2]
		for j := 1; j <= lb; j++ {
			if ra[i-1] == rb[j-1] {
				row[j] = prev[j-1] + 1
			} else if row[j-1] > prev[j] {
				row[j] = row[j-1]
			} else {
				row[j] = prev[j]
			}
		}
	}
	return float64(dp[la%2][lb]) * 200.0 / float64(la+lb)
}

func openChaDan(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		c.JSON(200, gin.H{"code": -1, "msg": "账号不能为空"})
		return
	}
	var orders []model.Order
	database.DB.Where("user = ?", username).Order("oid DESC").Find(&orders)
	if len(orders) == 0 {
		c.JSON(200, gin.H{"code": -1, "msg": "未查到该账号的下单信息"})
		return
	}
	// 兼容原 PHP 字段格式
	data := make([]gin.H, 0, len(orders))
	for _, o := range orders {
		data = append(data, gin.H{
			"id": o.OID, "ptname": o.PTName, "school": o.School,
			"name": o.Name, "user": o.UserAccount, "kcname": o.KCName,
			"addtime": o.AddTime,
			"courseStartTime": o.CourseStartTime, "courseEndTime": o.CourseEndTime,
			"examStartTime": o.ExamStartTime, "examEndTime": o.ExamEndTime,
			"status": o.Status, "process": o.Process, "remarks": o.Remarks,
		})
	}
	c.JSON(200, gin.H{"code": 1, "data": data})
}

func openBuDan(c *gin.Context) {
	oid := c.PostForm("id")
	var order model.Order
	if err := database.DB.Where("oid = ?", oid).First(&order).Error; err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": "订单不存在"})
		return
	}
	bsNum, _ := strconv.Atoi(order.BSNum)
	if bsNum > 5 {
		c.JSON(200, gin.H{"code": -1, "msg": "该订单补刷已超过5次，年轻人，要讲武德，我劝你好自为之"})
		return
	}
	// 调用上游补刷
	res := checkorder.RetryOrder(order.OID)
	if res.Code == 1 {
		database.DB.Model(&order).Updates(map[string]interface{}{
			"status": "补刷中", "bsnum": bsNum + 1,
		})
		c.JSON(200, gin.H{"code": 1, "msg": res.Msg})
		return
	}
	c.JSON(200, gin.H{"code": -1, "msg": res.Msg})
}

func openGetClass(c *gin.Context) {
	var classes []model.Class
	database.DB.Where("status = 1").Find(&classes)
	var data []gin.H
	for _, cls := range classes {
		data = append(data, gin.H{"cid": cls.CID, "name": cls.Name})
	}
	c.JSON(200, gin.H{"code": 1, "data": data})
}
