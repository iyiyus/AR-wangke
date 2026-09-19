package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"wk-go/internal/config"
	"wk-go/internal/database"
	"wk-go/internal/middleware"
	"wk-go/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type loginReq struct {
	Account  string `json:"account" binding:"required"`
	Pass     string `json:"pass" binding:"required"`
	AuthPass string `json:"auth_pass"`
}

type registerReq struct {
	Name    string `json:"name" binding:"required"`
	Account string `json:"account" binding:"required"`
	Pass    string `json:"pass" binding:"required"`
	YQM     string `json:"yqm"`
	Email   string `json:"email"`
}

// POST /api/auth/login
func Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	var user model.User
	if err := database.DB.Where("user = ? AND pass = ?", req.Account, req.Pass).First(&user).Error; err != nil {
		Fail(c, "账号或密码错误")
		return
	}
	if user.Active == "0" {
		Fail(c, "账号已被封禁")
		return
	}

	// 管理员需要二次验证
	if user.UID == 1 && req.AuthPass == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1001, "msg": "需要管理员认证"})
		return
	}
	if user.UID == 1 && req.AuthPass != config.Global.Admin.VerifyPassword {
		Fail(c, "管理员认证密码错误")
		return
	}

	token, err := middleware.GenerateToken(user.UID, user.User)
	if err != nil {
		Fail(c, "生成token失败")
		return
	}

	database.DB.Model(&user).Update("endtime", time.Now().Format("2006-01-02 15:04:05"))

	OK(c, gin.H{
		"token": token,
		"user": gin.H{
			"uid":      user.UID,
			"user":     user.User,
			"name":     user.Name,
			"money":    user.Money,
			"addprice": user.AddPrice,
			"key":      user.Key,
			"yqm":      user.YQM,
			"active":   user.Active,
		},
	})
}

// POST /api/auth/register
func Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	// 邀请码必填
	if req.YQM == "" {
		Fail(c, "邀请码不能为空")
		return
	}

	var count int64
	database.DB.Model(&model.User{}).Where("user = ?", req.Account).Count(&count)
	if count > 0 {
		Fail(c, "账号已存在")
		return
	}

	var cfgYQZC model.Config
	database.DB.Where("v = ?", "user_yqzc").First(&cfgYQZC)
	if cfgYQZC.K == "0" {
		Fail(c, "系统未开启邀请码注册")
		return
	}

	var inviter model.User
	if err := database.DB.Where("yqm = ?", req.YQM).First(&inviter).Error; err != nil {
		Fail(c, "邀请码无效")
		return
	}
	if inviter.YQPrice == "" {
		Fail(c, "邀请人未设置邀请费率，无法注册")
		return
	}
	if len(req.Pass) < 6 {
		Fail(c, "密码至少6位")
		return
	}
	uuid_ := inviter.UID

	ip := c.ClientIP()
	// 同 IP 每日注册限制 50
	var ipCount int64
	database.DB.Model(&model.User{}).
		Where("ip = ? AND addtime LIKE ?", ip, time.Now().Format("2006-01-02")+"%").
		Count(&ipCount)
	if ipCount >= 50 {
		Fail(c, "该IP今日注册次数已达上限")
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	faceimg := ""
	if req.Account != "" && len(req.Account) >= 5 && len(req.Account) <= 11 {
		faceimg = fmt.Sprintf("http://q.qlogo.cn/headimg_dl?dst_uin=%s&spec=100&img_type=jpg", req.Account)
	}
	user := model.User{
		UUID:     uuid_,
		User:     req.Account,
		Pass:     req.Pass,
		Name:     req.Name,
		AddPrice: parsePrice(inviter.YQPrice),
		Key:      "0",
		AddTime:  now,
		IP:       ip,
		Active:   "1",
		FaceImg:  faceimg,
		Email:    req.Email,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		Fail(c, "注册失败")
		return
	}
	OK(c, nil)
}

// GET /api/user/info
func GetUserInfo(c *gin.Context) {
	uid, _ := c.Get("uid")
	var user model.User
	if err := database.DB.First(&user, uid).Error; err != nil {
		Fail(c, "用户不存在")
		return
	}
	// 安全校验（兼容原 PHP userinfo）：费率异常复位、余额大于总充值冻结（管理员除外）
	if user.UID != 1 {
		if user.AddPrice < 0.1 {
			database.DB.Model(&user).Update("addprice", 1.0)
			user.AddPrice = 1.0
		}
		if user.Money > parsePrice(user.ZCZ) {
			database.DB.Model(&user).Update("active", "0")
			user.Active = "0"
		}
	}
	lv := 0.0
	if user.CK > 0 {
		lv = float64(user.DD) * 100.0 / float64(user.CK)
	}
	OK(c, gin.H{
		"uid":      user.UID,
		"uuid":     user.UUID,
		"user":     user.User,
		"name":     user.Name,
		"money":    user.Money,
		"zcz":      user.ZCZ,
		"ck":       user.CK,
		"dd":       user.DD,
		"lv":       lv,
		"addprice": user.AddPrice,
		"key":      user.Key,
		"yqm":      user.YQM,
		"yqprice":  user.YQPrice,
		"notice":   user.Notice,
		"active":   user.Active,
		"faceimg":  user.FaceImg,
		"nickname": user.Nickname,
		"email":    user.Email,
	})
}

// POST /api/user/passwd
func ChangePassword(c *gin.Context) {
	uid, _ := c.Get("uid")
	var req struct {
		OldPass string `json:"old_pass" binding:"required"`
		NewPass string `json:"new_pass" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	var user model.User
	database.DB.First(&user, uid)
	if user.Pass != req.OldPass {
		Fail(c, "原密码错误")
		return
	}
	database.DB.Model(&user).Update("pass", req.NewPass)
	OK(c, nil)
}

// POST /api/user/api-key/open  自开通（兼容 ktapi type1：余额<300 扣 10 元，否则免费）
func OpenAPIKey(c *gin.Context) {
	uid, _ := c.Get("uid")
	var user model.User
	database.DB.First(&user, uid)
	if user.Key != "0" {
		Fail(c, "已开通")
		return
	}
	if user.Money < 300 {
		if user.Money < 10 {
			Fail(c, "余额不足10元，无法扣除开通费用")
			return
		}
		database.DB.Model(&user).UpdateColumn("money", gorm.Expr("money - ?", 10))
		writeLog(user.UID, "开通接口", "开通API接口扣除10元", -10, "")
	}
	key := uuid.New().String()
	database.DB.Model(&user).Update("key", key)
	OK(c, gin.H{"key": key})
}

// POST /api/user/api-key/reset
func ResetAPIKey(c *gin.Context) {
	uid, _ := c.Get("uid")
	key := uuid.New().String()
	database.DB.Model(&model.User{}).Where("uid = ?", uid).Update("key", key)
	OK(c, gin.H{"key": key})
}

// POST /api/user/notice
func UpdateNotice(c *gin.Context) {
	uid, _ := c.Get("uid")
	var req struct {
		Notice string `json:"notice"`
	}
	c.ShouldBindJSON(&req)
	database.DB.Model(&model.User{}).Where("uid = ?", uid).Update("notice", req.Notice)
	OK(c, nil)
}

// POST /api/user/invite  设置邀请码/邀请费率（费率须为 0.05 倍数）
func SetInvite(c *gin.Context) {
	uid, _ := c.Get("uid")
	var user model.User
	database.DB.First(&user, uid)

	var req struct {
		YQM     string `json:"yqm"`
		YQPrice string `json:"yqprice"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	// 邀请码：未设置时允许自动生成
	yqm := req.YQM
	if yqm == "" && user.YQM == "" {
		yqm = randomString(5)
	}
	if yqm != "" {
		if len(yqm) < 4 {
			Fail(c, "邀请码至少4位")
			return
		}
		var count int64
		database.DB.Model(&model.User{}).Where("yqm = ? AND uid != ?", yqm, uid).Count(&count)
		if count > 0 {
			Fail(c, "邀请码已被使用")
			return
		}
		updates["yqm"] = yqm
	}
	if req.YQPrice != "" {
		p := parsePrice(req.YQPrice)
		if p <= 0 {
			Fail(c, "邀请费率错误")
			return
		}
		// 必须是 0.05 的倍数
		if m := p / 0.05; m-float64(int(m)) > 1e-9 {
			Fail(c, "邀请费率必须为0.05的倍数")
			return
		}
		updates["yqprice"] = req.YQPrice
	}
	if len(updates) == 0 {
		Fail(c, "无修改内容")
		return
	}
	database.DB.Model(&model.User{}).Where("uid = ?", uid).Updates(updates)
	OK(c, gin.H{"yqm": yqm})
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GET /api/user/boss/:uid
func GetBossInfo(c *gin.Context) {
	uid := c.Param("uid")
	var user model.User
	if err := database.DB.Select("uid,user,name,yqm,notice").Where("uid = ?", uid).First(&user).Error; err != nil {
		Fail(c, "用户不存在")
		return
	}
	OK(c, gin.H{
		"uid":    user.UID,
		"user":   user.User,
		"name":   user.Name,
		"yqm":    user.YQM,
		"notice": user.Notice,
	})
}

// POST /api/user/avatar 更新头像
func UpdateAvatar(c *gin.Context) {
	uidVal, exists := c.Get("uid")
	if !exists || uidVal == nil {
		Fail(c, "未登录")
		return
	}
	var body struct {
		Avatar string `json:"avatar"`
	}
	if c.ShouldBindJSON(&body) != nil || body.Avatar == "" {
		Fail(c, "参数错误")
		return
	}
	database.DB.Model(&model.User{}).Where("uid = ?", uidVal).Update("faceimg", body.Avatar)
	OK(c, gin.H{"msg": "头像更新成功", "faceimg": body.Avatar})
}

// POST /api/user/email 更新邮箱
func UpdateEmail(c *gin.Context) {
	uidVal, exists := c.Get("uid")
	if !exists || uidVal == nil {
		Fail(c, "未登录")
		return
	}
	var body struct {
		Email string `json:"email"`
	}
	if c.ShouldBindJSON(&body) != nil || body.Email == "" {
		Fail(c, "参数错误")
		return
	}
	database.DB.Model(&model.User{}).Where("uid = ?", uidVal).Update("email", body.Email)
	OK(c, gin.H{"msg": "邮箱更新成功", "email": body.Email})
}

