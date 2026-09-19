package handler

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"wk-go/internal/config"
	"wk-go/internal/database"
	"wk-go/internal/model"
)

/* ============================================================
 * 安装系统（复刻 UDID 安装引导逻辑，适配 wk-go）
 * 流程：许可协议 → 环境检测 → 系统配置 → 注册 systemd → 安装完成（自动重启）
 * 判定：install.lock 文件存在即为已安装
 * ============================================================ */

const installLockFile = "install.lock"

// IsInstalled 是否已安装（install.lock 检测，兼容 exe 同目录与当前目录）
func IsInstalled() bool {
	if exe, err := os.Executable(); err == nil {
		if _, err := os.Stat(filepath.Join(filepath.Dir(exe), installLockFile)); err == nil {
			return true
		}
	}
	_, err := os.Stat(installLockFile)
	return err == nil
}

func installLockPath() string {
	if exe, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(exe), installLockFile)
	}
	return installLockFile
}

func workDirOfExe() string {
	if exe, err := os.Executable(); err == nil {
		return filepath.Dir(exe)
	}
	wd, _ := os.Getwd()
	return wd
}

// ---------- 在线重启（复刻 go-restart-solutions 方案三：systemd） ----------

var (
	lastRestartTime time.Time
	restartMutex    sync.Mutex
)

// canRestart 频率限制：5 分钟 1 次
func canRestart() bool {
	restartMutex.Lock()
	defer restartMutex.Unlock()
	if time.Since(lastRestartTime) < 5*time.Minute {
		return false
	}
	lastRestartTime = time.Now()
	return true
}

// ---------- 接口 ----------

// InstallStatus 安装状态 + 环境检测 GET /api/install/status
func InstallStatus(c *gin.Context) {
	// 数据库连通检测
	dbOK := false
	if database.DB != nil {
		sqlDB, err := database.DB.DB()
		if err == nil && sqlDB.Ping() == nil {
			dbOK = true
		}
	}
	// 目录可写检测
	writable := false
	if dir := workDirOfExe(); dir != "" {
		f, err := os.CreateTemp(dir, ".wkwrite_*")
		if err == nil {
			f.Close()
			os.Remove(f.Name())
			writable = true
		}
	}
	c.JSON(http.StatusOK, OKResult(gin.H{
		"installed": IsInstalled(),
		"database":  dbOK,
		"writable":  writable,
		"os":        runtime.GOOS,
	}))
}

// InstallCheckDB 测试数据库连接 POST /api/install/check-db
func InstallCheckDB(c *gin.Context) {
	var req struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	if req.Host == "" || req.User == "" || req.DBName == "" {
		Fail(c, "请填写完整的数据库信息")
		return
	}
	if req.Port == "" {
		req.Port = "3306"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		req.User, req.Password, req.Host, req.Port, req.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil || db.Ping() != nil {
		Fail(c, "数据库连接失败")
		return
	}
	db.Close()
	OK(c, gin.H{"msg": "连接成功"})
}

// InstallRegisterSystemd 注册 systemd 服务 POST /api/install/register-systemd
func InstallRegisterSystemd(c *gin.Context) {
	var req struct {
		ServiceName string `json:"serviceName"`
	}
	c.ShouldBindJSON(&req)

	exePath, err := os.Executable()
	if err != nil {
		Fail(c, "获取程序路径失败")
		return
	}
	exePath, _ = filepath.Abs(exePath)
	workDir := filepath.Dir(exePath)
	envFile := filepath.Join(workDir, ".env")

	serviceName := strings.TrimSpace(req.ServiceName)
	if serviceName == "" {
		serviceName = "wk-go"
	}
	// 仅允许安全的服务名
	if !regexpServiceName(serviceName) {
		Fail(c, "服务名不合法（仅限字母数字-_）")
		return
	}

	serviceContent := fmt.Sprintf(`[Unit]
Description=网课代刷管理系统
After=network.target

[Service]
Type=simple
WorkingDirectory=%s
ExecStart=%s
Restart=always
RestartSec=5
EnvironmentFile=-%s

[Install]
WantedBy=multi-user.target
`, workDir, exePath, envFile)

	servicePath := fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)
	if err := os.WriteFile(servicePath, []byte(serviceContent), 0644); err != nil {
		Fail(c, "写入服务文件失败: "+err.Error())
		return
	}

	exec.Command("systemctl", "daemon-reload").Run()
	exec.Command("systemctl", "enable", serviceName).Run()
	startOut, startErr := exec.Command("systemctl", "start", serviceName).CombinedOutput()
	if startErr != nil {
		if out, err := exec.Command("systemctl", "restart", serviceName).CombinedOutput(); err != nil {
			Fail(c, "服务启动失败: "+string(out))
			return
		}
	}
	_ = startOut

	OK(c, gin.H{
		"serviceName": serviceName,
		"exePath":     exePath,
		"workDir":     workDir,
	})
}

// InstallRun 执行安装 POST /api/install/run
func InstallRun(c *gin.Context) {
	if IsInstalled() {
		Fail(c, "系统已安装")
		return
	}

	var req struct {
		DBHost      string `json:"dbHost"`
		DBPort      string `json:"dbPort"`
		DBUser      string `json:"dbUser"`
		DBPassword  string `json:"dbPassword"`
		DBName      string `json:"dbName"`
		ServerPort  string `json:"serverPort"`
		ProjectName string `json:"projectName"`
		SiteName    string `json:"siteName"`
		AdminUser   string `json:"adminUser"`
		AdminPwd    string `json:"adminPwd"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	if req.DBHost == "" || req.DBUser == "" || req.DBName == "" {
		Fail(c, "请填写完整的数据库信息")
		return
	}
	if req.DBPort == "" {
		req.DBPort = "3306"
	}
	if req.ServerPort == "" {
		req.ServerPort = "8080"
	}
	if req.ProjectName == "" {
		req.ProjectName = "wk-go"
	}
	if req.SiteName == "" {
		req.SiteName = "网课代刷管理系统"
	}
	if req.AdminUser == "" {
		req.AdminUser = "admin"
	}
	if req.AdminPwd == "" {
		req.AdminPwd = "admin123"
	}

	// 1. 创建数据库（如不存在）
	rootDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		req.DBUser, req.DBPassword, req.DBHost, req.DBPort)
	adminDB, err := sql.Open("mysql", rootDSN)
	if err != nil {
		Fail(c, "数据库连接失败: "+err.Error())
		return
	}
	if _, err := adminDB.Exec(fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", req.DBName)); err != nil {
		adminDB.Close()
		Fail(c, "创建数据库失败: "+err.Error())
		return
	}
	adminDB.Close()

	// 2. 连接目标库并导入基础表结构（database/init.sql）
	gormDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		req.DBUser, req.DBPassword, req.DBHost, req.DBPort, req.DBName)
	if err := database.InitWithDSN(gormDSN); err != nil {
		Fail(c, "数据库初始化失败: "+err.Error())
		return
	}
	if err := execInitSQL(gormDSN); err != nil {
		Fail(c, "初始化数据表失败: "+err.Error())
		return
	}

	// 3. 更新内存配置并写回 config.yaml
	cfg := config.Global
	cfg.Database.Host = req.DBHost
	if p, err := atoiSafe(req.DBPort); err == nil {
		cfg.Database.Port = p
	}
	cfg.Database.User = req.DBUser
	cfg.Database.Password = req.DBPassword
	cfg.Database.DBName = req.DBName
	if p, err := atoiSafe(req.ServerPort); err == nil && p > 0 {
		cfg.Server.Port = p
	}
	// 管理员认证密码与安装密码保持一致，避免安装后不知 auth_pass 无法登录
	cfg.Admin.VerifyPassword = req.AdminPwd
	if err := config.Save(); err != nil {
		Fail(c, "写入配置失败: "+err.Error())
		return
	}

	// 4. 创建管理员账号（明文比对，与登录逻辑一致）
	now := time.Now().Format("2006-01-02 15:04:05")
	adminKey := genKey(16)
	adminYQM := genInviteCode(8)
	if res := database.DB.Exec(`DELETE FROM qingka_wangke_user WHERE user = ?`, req.AdminUser); res.Error != nil {
		Fail(c, "初始化管理员失败: "+res.Error.Error())
		return
	}
	// 注意 key 为 MySQL 保留字，必须加反引号；补齐所有 NOT NULL 无默认值列
	ins := database.DB.Exec("INSERT INTO qingka_wangke_user (uid, uuid, user, pass, name, ck, xdlv, dd, qq_openid, nickname, faceimg, money, zcz, addprice, `key`, yqm, yqprice, notice, addtime, endtime, ip, grade, active) VALUES (1, 0, ?, ?, ?, 0, 0, 0, '', '', '', 0, '0', 0.1, ?, ?, '', '', ?, ?, '', '', '1')",
		req.AdminUser, req.AdminPwd, req.AdminUser, adminKey, adminYQM, now, now)
	if ins.Error != nil {
		Fail(c, "初始化管理员失败: "+ins.Error.Error())
		return
	}

	// 5. 站点名称写配置表
	if res := database.DB.Exec(`INSERT INTO qingka_wangke_config (v, k) VALUES ('site_name', ?) ON DUPLICATE KEY UPDATE k = VALUES(k)`, req.SiteName); res.Error != nil {
		Fail(c, "写入站点配置失败: "+res.Error.Error())
		return
	}

	// 6. 写 .env（供 systemd EnvironmentFile 使用）
	workDir := workDirOfExe()
	envContent := fmt.Sprintf("DB_HOST=%s\nDB_PORT=%s\nDB_USER=%s\nDB_PASSWORD=%s\nDB_NAME=%s\nSERVER_PORT=%s\n",
		req.DBHost, req.DBPort, req.DBUser, req.DBPassword, req.DBName, req.ServerPort)
	os.WriteFile(filepath.Join(workDir, ".env"), []byte(envContent), 0644)

	// 7. 写 install.lock
	os.WriteFile(installLockPath(), []byte("installed at "+now), 0644)

	// 8. 异步自动重启（systemd，2 秒后执行）
	go func() {
		time.Sleep(2 * time.Second)
		exec.Command("systemctl", "restart", req.ProjectName).Run()
	}()

	OK(c, gin.H{
		"adminUser":  req.AdminUser,
		"adminPwd":   req.AdminPwd,
		"serverPort": req.ServerPort,
		"siteName":   req.SiteName,
	})
}

// execInitSQL 执行 database/init.sql 初始化表结构
func execInitSQL(dsn string) error {
	workDir := workDirOfExe()
	candidates := []string{
		filepath.Join(workDir, "database", "init.sql"),
		"database/init.sql",
		"../database/init.sql",
	}
	var sqlPath string
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			sqlPath = p
			break
		}
	}
	if sqlPath == "" {
		// 无 init.sql 时尝试 AutoMigrate 已知模型
		return database.DB.AutoMigrate(&model.User{}, &model.Order{}, &model.Config{}, &model.SxdkOrder{})
	}
	data, err := os.ReadFile(sqlPath)
	if err != nil {
		return err
	}
	// 执行多语句 SQL 需要开启 multiStatements
	db, err := sql.Open("mysql", dsn+"&multiStatements=true")
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec(string(data)); err != nil {
		return err
	}
	return nil
}

// InstallRestart 在线重启服务 POST /api/install/restart
// 复刻 go-restart-solutions 方案三（systemd），带管理员认证 + 5 分钟频率限制
func InstallRestart(c *gin.Context) {
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(int64)
	if uid != 1 {
		Fail(c, "权限不足")
		return
	}
	if !canRestart() {
		Fail(c, "操作过于频繁，请 5 分钟后再试")
		return
	}
	serviceName := strings.TrimSpace(c.PostForm("serviceName"))
	if serviceName == "" {
		serviceName = "wk-go"
	}
	if !regexpServiceName(serviceName) {
		Fail(c, "服务名不合法")
		return
	}
	if runtime.GOOS == "windows" {
		// Windows 开发环境：不支持 systemctl
		Fail(c, "当前系统不支持 systemd 服务，请在部署环境（Linux）使用")
		return
	}
	cmd := exec.Command("systemctl", "restart", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		writeLog(uid, "install", "在线重启失败: "+string(output), 0, c.ClientIP())
		Fail(c, "重启失败: "+string(output))
		return
	}
	writeLog(uid, "install", fmt.Sprintf("在线重启服务 %s 成功", serviceName), 0, c.ClientIP())
	OK(c, gin.H{"msg": "服务重启中..."})
}

// ---------- 辅助 ----------

func regexpServiceName(s string) bool {
	if s == "" || len(s) > 64 {
		return false
	}
	for _, r := range s {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

func genKey(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}

func genInviteCode(n int) string {
	return genKey(n)
}

// OKResult 与 OK 相同的成功结构（供 status 等直接返回）
func OKResult(data gin.H) gin.H {
	return gin.H{"code": 200, "msg": "success", "data": data}
}
