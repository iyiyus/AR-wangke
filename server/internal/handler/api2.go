package handler

import (
	"fmt"
	"wk-go/internal/checkorder"
	"wk-go/internal/database"
	"wk-go/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 二套对接 API（兼容原 PHP /api/index.php）与联盟接口（兼容 /lmapi.php）

func API2(c *gin.Context) {
	switch c.Query("act") {
	case "getmoney":
		api2GetMoney(c)
	case "class":
		api2Class(c)
	case "chake":
		api2Chake(c)
	case "cd":
		api2Cd(c)
	case "bd":
		api2Bd(c)
	case "up":
		api2Up(c)
	default:
		c.JSON(200, gin.H{"code": 0, "msg": "未知操作"})
	}
}

func api2GetMoney(c *gin.Context) {
	u := openAuth(c)
	if u == nil {
		return
	}
	c.JSON(200, gin.H{"code": 1, "msg": "查询成功", "user": u.User, "name": u.Name, "money": u.Money})
}

func api2Class(c *gin.Context) {
	u := openAuth(c)
	if u == nil {
		return
	}
	var classes []model.Class
	database.DB.Where("status = 1").Order("sort ASC, cid DESC").Find(&classes)
	data := make([]gin.H, 0, len(classes))
	for _, cls := range classes {
		price := parsePrice(cls.Price)
		data = append(data, gin.H{
			"sort": cls.Sort, "cid": cls.CID, "name": cls.Name,
			"content": cls.Content, "status": cls.Status, "price": cls.Price,
			"price5": fmt.Sprintf("%.2f", price+0.5),
			"jiage":  fmt.Sprintf("%.2f", price*u.AddPrice),
		})
	}
	c.JSON(200, gin.H{"code": 1, "data": data})
}

func api2Chake(c *gin.Context) {
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

	// zdmoney 余额限制
	zdmoney := getConfFloat("zdmoney", 10)
	if u.Money < zdmoney {
		c.JSON(200, gin.H{"code": -2, "msg": fmt.Sprintf("余额小于%v禁止调用查课", zdmoney)})
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

	// 查课/下单比例限制
	var ckCount, tjCount, tjCount1 int64
	database.DB.Model(&model.Log{}).Where("uid = ? AND type = ?", u.UID, "API查课").Count(&ckCount)
	database.DB.Model(&model.Log{}).Where("uid = ? AND type = ?", u.UID, "添加任务").Count(&tjCount)
	database.DB.Model(&model.Log{}).Where("uid = ? AND type = ?", u.UID, "API添加任务").Count(&tjCount1)
	taskTotal := tjCount + tjCount1
	if taskTotal == 0 {
		taskTotal = 1
	}
	bl := float64(ckCount) / float64(taskTotal)
	blLimit := getConfFloat("bl", 999)
	if bl > blLimit {
		c.JSON(200, gin.H{"code": -2, "msg": fmt.Sprintf("API查课比例超%v禁止查课", blLimit)})
		return
	}

	res := checkorder.QueryByClass(cls.CID, school, user, pass)
	if res == nil {
		res = &checkorder.QueryResult{Code: -1, Msg: "查询失败,请联系管理员"}
	}
	writeLog(u.UID, "API查课", fmt.Sprintf("%s-查课信息：%s %s %s", cls.Name, school, user, pass), 0, c.ClientIP())
	database.DB.Model(&model.User{}).Where("uid = ?", uid).UpdateColumn("ck", gorm.Expr("ck + 1"))

	if typ == "xiaochu" && res.Code == 1 {
		names := make([]string, 0, len(res.Data))
		for _, course := range res.Data {
			names = append(names, course.Name)
		}
		c.JSON(200, gin.H{
			"code": res.Code, "msg": res.Msg,
			"data": []string{cls.Name, user, pass, school, joinNames(names)},
			"js":   "", "info": "昔日之苦，安知异日不在尝之? 共勉",
		})
		return
	}
	c.JSON(200, gin.H{"code": res.Code, "msg": res.Msg, "data": res.Data, "userName": res.UserName, "userinfo": school + " " + user + " " + pass})
}

func api2Cd(c *gin.Context) {
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

func api2Bd(c *gin.Context) {
	oid := c.PostForm("id")
	var order model.Order
	if err := database.DB.Where("oid = ?", oid).First(&order).Error; err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": "订单不存在"})
		return
	}
	bsNum, _ := atoiSafe(order.BSNum)
	if bsNum > 20 {
		c.JSON(200, gin.H{"code": -1, "msg": "该订单补刷已超过20次，年轻人，要讲武德，我劝你好自为之"})
		return
	}
	res := checkorder.RetryOrder(order.OID)
	if res.Code == 1 {
		database.DB.Model(&order).Updates(map[string]interface{}{
			"status": "补刷中", "bsnum": gorm.Expr("bsnum + 1"),
		})
		c.JSON(200, gin.H{"code": 1, "msg": res.Msg})
		return
	}
	c.JSON(200, gin.H{"code": -1, "msg": res.Msg})
}

func api2Up(c *gin.Context) {
	oid := c.PostForm("id")
	var order model.Order
	if err := database.DB.Where("oid = ?", oid).First(&order).Error; err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": "订单不存在"})
		return
	}
	items := checkorder.ProgressByOrder(order.OID)
	for _, p := range items {
		database.DB.Model(&model.Order{}).
			Where("user = ? AND pass = ? AND kcname = ?", p.User, order.Pass, p.KCName).
			Updates(map[string]interface{}{
				"yid":             p.YID,
				"status":          p.StatusText,
				"courseStartTime": p.KCStartTime,
				"courseEndTime":   p.KCEndTime,
				"examStartTime":   p.KSStartTime,
				"examEndTime":     p.KSEndTime,
				"process":         p.Process,
				"remarks":         p.Remarks,
			})
	}
	c.JSON(200, gin.H{"code": 1, "msg": "同步成功，请重新查询信息"})
}

// ===== 联盟接口（兼容原 PHP /lmapi.php）=====

func LMAPI(c *gin.Context) {
	switch c.Query("act") {
	case "get":
		lmGet(c)
	case "add":
		lmAdd(c)
	case "chadan":
		lmChaDan(c)
	case "tongbu":
		lmTongBu(c)
	case "budan":
		lmBuDan(c)
	default:
		c.JSON(200, gin.H{"code": 0, "msg": "未知操作"})
	}
}

func lmGet(c *gin.Context) {
	u := openAuth(c)
	if u == nil {
		return
	}
	var cls model.Class
	if err := database.DB.Where("cid = ?", c.PostForm("platform")).First(&cls).Error; err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": "平台不存在"})
		return
	}
	if cls.Status == 0 {
		c.JSON(200, gin.H{"code": -2, "msg": "网课已下架禁止查课！"})
		return
	}
	res := checkorder.QueryByClass(cls.CID, c.PostForm("school"), c.PostForm("user"), c.PostForm("pass"))
	if res == nil {
		res = &checkorder.QueryResult{Code: -1, Msg: "查询失败,请联系管理员"}
	}
	c.JSON(200, gin.H{"code": res.Code, "msg": res.Msg, "data": res.Data})
}

func lmAdd(c *gin.Context) {
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
	if cls.Status == 0 {
		c.JSON(200, gin.H{"code": -2, "msg": "商品已下架"})
		return
	}
	// 联盟固定单价 1（兼容原 PHP）
	danjia := 1.0
	if u.Money < danjia {
		c.JSON(200, gin.H{"code": -1, "msg": "余额不足"})
		return
	}
	dockStatus := "0"
	if cls.Docking == "0" {
		dockStatus = "99"
	}
	order := model.Order{
		UID: u.UID, CID: cls.CID,
		HID: parseInt64(cls.Docking), PTName: cls.Name,
		School: school, UserAccount: userAcc, Pass: pass,
		KCName: kcname, Fees: fmt.Sprintf("%.2f", danjia), Noun: cls.Noun,
		AddTime: nowStr(), IP: c.ClientIP(), DockStatus: dockStatus, Status: "待处理",
	}
	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": "提交失败"})
		return
	}
	database.DB.Model(&model.User{}).Where("uid = ?", u.UID).
		UpdateColumn("money", gorm.Expr("money - ?", danjia))
	writeLog(u.UID, "API添加任务", fmt.Sprintf("%s %s %s 扣除%.2f元！", userAcc, pass, kcname, danjia), -danjia, c.ClientIP())
	c.JSON(200, gin.H{"code": 0, "msg": "提交成功", "status": 0, "message": "提交成功"})
}

func lmChaDan(c *gin.Context) {
	api2Cd(c)
}

func lmTongBu(c *gin.Context) {
	api2Up(c)
}

func lmBuDan(c *gin.Context) {
	oid := c.PostForm("id")
	var order model.Order
	if err := database.DB.Where("oid = ?", oid).First(&order).Error; err != nil {
		c.JSON(200, gin.H{"code": -1, "msg": "订单不存在"})
		return
	}
	bsNum, _ := atoiSafe(order.BSNum)
	if bsNum > 5 {
		c.JSON(200, gin.H{"code": -1, "msg": "该订单补刷已超过5次"})
		return
	}
	res := checkorder.RetryOrder(order.OID)
	if res.Code == 1 {
		database.DB.Model(&order).Updates(map[string]interface{}{
			"status": "补刷中", "bsnum": gorm.Expr("bsnum + 1"),
		})
		c.JSON(200, gin.H{"code": 1, "msg": res.Msg})
		return
	}
	c.JSON(200, gin.H{"code": -1, "msg": res.Msg})
}

// ===== 工具 =====

func getConfFloat(key string, def float64) float64 {
	var cfg model.Config
	if err := database.DB.Where("v = ?", key).First(&cfg).Error; err != nil {
		return def
	}
	v := parsePrice(cfg.K)
	if v == 0 {
		return def
	}
	return v
}

func joinNames(names []string) string {
	out := ""
	for i, n := range names {
		if i == 0 {
			out = n
		} else {
			out += "," + n
		}
	}
	return out
}
