package router

import (
	"wk-go/internal/handler"
	"wk-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.Use(middleware.CORS())

	// 安装引导接口（无需认证；install.lock 判定已安装）
	install := r.Group("/api/install")
	{
		install.GET("/status", handler.InstallStatus)
		install.POST("/check-db", handler.InstallCheckDB)
		install.POST("/register-systemd", handler.InstallRegisterSystemd)
		install.POST("/run", handler.InstallRun)
	}

	r.GET("/api/site", handler.SiteInfo)
	// 文件上传
	r.POST("/api/upload", handler.Upload)

	// 开放 API（兼容原 PHP）
	r.Any("/open/api", handler.OpenAPI)
	// 兼容原 PHP 根路径 api.php 与二套 API /api/index.php、联盟 /lmapi.php
	r.Any("/api.php", handler.OpenAPI)
	r.Any("/api/index.php", handler.API2)
	r.Any("/lmapi.php", handler.LMAPI)

	// 支付回调（无需鉴权）
	r.POST("/api/pay/notify", handler.PayNotify)

	// 认证
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)
		auth.POST("/forget-password", handler.ForgetPassword)
		auth.POST("/reset-password", handler.ResetPassword)
	}

	// 需要登录
	api := r.Group("/api", middleware.JWTAuth())
	{
		// 用户
		api.GET("/user/info", handler.GetUserInfo)
		api.POST("/user/passwd", handler.ChangePassword)
		api.POST("/user/avatar", handler.UpdateAvatar)
		api.POST("/user/email", handler.UpdateEmail)
		api.POST("/user/api-key/open", handler.OpenAPIKey)
		api.POST("/user/api-key/reset", handler.ResetAPIKey)
		api.POST("/user/notice", handler.UpdateNotice)
		api.POST("/user/invite", handler.SetInvite)
		api.GET("/user/boss/:uid", handler.GetBossInfo)

		// 订单
		api.GET("/order/list", handler.OrderList)
		api.POST("/order/add", handler.OrderAdd)
		api.POST("/order/cancel", handler.OrderCancel)
		api.POST("/order/restart", handler.OrderRestart)
		api.GET("/order/export", handler.OrderExport)
		api.GET("/order/monthly", handler.OrderMonthly)
		api.GET("/order/by-platform", handler.OrderByPlatform)
		api.POST("/order/query", handler.OrderQuery)

		// 平台
		api.GET("/class/all", handler.ClassAll)
		api.GET("/class/list", handler.ClassList)
		api.GET("/fenlei/list", handler.FenLeiList)

		// 充值
		api.POST("/pay/create", handler.PayCreate)
		api.GET("/pay/list", handler.PayList)

		// 工单
		api.GET("/gongdan/list", handler.GongDanList)
		api.POST("/gongdan/add", handler.GongDanAdd)
		api.POST("/gongdan/toanswer", handler.GongDanToAnswer)

		// 代理对下级的操作
		api.GET("/proxy/user/list", handler.ProxyUserList)
		api.POST("/proxy/user/add", handler.ProxyUserAdd)
		api.POST("/proxy/user/recharge", handler.ProxyUserRecharge)
		api.POST("/proxy/user/price", handler.ProxyUserPrice)
		api.POST("/proxy/user/reset-pass", handler.ProxyUserResetPass)
		api.POST("/proxy/user/ban", handler.ProxyUserBan)
		api.POST("/proxy/user/yqm", handler.ProxyUserYQM)
		api.POST("/proxy/user/api-key", handler.ProxyUserOpenKey)
		api.POST("/user/migrate", handler.UserMigrate)

		// 密价（自己查看）
		api.GET("/myprice/list", handler.MyPriceList)

		// 实习打卡（daka 对接，源台协议与 PHP sxdk/api.php 一致）
		api.Any("/daka/api", handler.DakaAPI)

		// 在线重启（管理员，复刻 UDID 自动重启；5 分钟频率限制）
		api.POST("/install/restart", middleware.AdminOnly(), handler.InstallRestart)

		// 盖章/病历（taowa 对接）
		api.GET("/taowa/companies", handler.TaowaCompanies)
		api.GET("/taowa/templates", handler.TaowaTemplates)
		api.POST("/taowa/spec", handler.TaowaSpec)
		api.GET("/taowa/delivery", handler.TaowaDelivery)
		api.GET("/taowa/notice", handler.TaowaNotice)
		api.POST("/taowa/order/add", handler.TaowaOrderAdd)
		api.POST("/taowa/order/add-bl", handler.TaowaOrderAddBL)
		api.GET("/taowa/order/list", handler.TaowaOrderList)
		api.POST("/taowa/order/cancel", handler.TaowaOrderCancel)
		api.POST("/taowa/order/remark", handler.TaowaOrderRemark)
		api.POST("/taowa/order/gxx", handler.TaowaOrderGxx)
		api.POST("/taowa/order/status", handler.TaowaOrderStatus)
		api.POST("/taowa/order/tk", handler.TaowaOrderTk)
		api.POST("/taowa/order/bjtk", handler.TaowaOrderBjtk)
		api.POST("/taowa/order/del", handler.TaowaOrderDel)
		api.GET("/taowa/order/ddlog", handler.TaowaDdlog)

		// 等级列表（注册/下单时用）
		api.GET("/dengji/list", handler.DengjiList)

		// 统计（管理员）
		api.GET("/admin/stats", middleware.AdminOnly(), handler.AdminStats)
		api.GET("/admin/dashboard", middleware.AdminOnly(), handler.AdminDashboard)

		// 管理员 - 用户管理
		admin := api.Group("/admin", middleware.AdminOnly())
		{
			admin.GET("/users", handler.AdminUserList)
			admin.POST("/user/add", handler.AdminUserAdd)
			admin.POST("/user/recharge", handler.AdminUserRecharge)
			admin.POST("/user/ban", handler.AdminUserBan)
			admin.POST("/user/reset-pass", handler.AdminUserResetPass)
			admin.POST("/user/level", handler.AdminUserLevel)
			admin.POST("/user/api-key", handler.AdminUserOpenKey)
			admin.POST("/user/yqm", handler.AdminUserYQM)

			// 系统配置
			admin.GET("/config", handler.AdminGetConfig)
			admin.POST("/config", handler.AdminSaveConfig)

			// 等级管理
			admin.POST("/dengji", handler.DengjiAdd)
			admin.PUT("/dengji/:id", handler.DengjiUpdate)
			admin.DELETE("/dengji/:id", handler.DengjiDelete)

			// 平台管理
			admin.POST("/class", handler.ClassAdd)
			admin.PUT("/class/:id", handler.ClassUpdate)
			admin.DELETE("/class/:id", handler.ClassDelete)
			admin.GET("/class/remote", handler.ClassRemote)
			admin.POST("/class/batch-import", handler.ClassBatchImport)
			admin.POST("/taowa/sync", handler.TaowaSync)

			// 分类管理
			admin.POST("/fenlei", handler.FenLeiAdd)
			admin.PUT("/fenlei/:id", handler.FenLeiUpdate)
			admin.DELETE("/fenlei/:id", handler.FenLeiDelete)

			// 货源管理
			admin.GET("/huoyuan", handler.HuoYuanList)
			admin.POST("/huoyuan", handler.HuoYuanAdd)
			admin.PUT("/huoyuan/:id", handler.HuoYuanUpdate)
			admin.DELETE("/huoyuan/:id", handler.HuoYuanDelete)

			// 订单管理
			admin.POST("/order/status", handler.OrderBatchStatus)
			admin.POST("/order/dock", handler.OrderBatchDock)
			admin.POST("/order/refund", handler.OrderRefund)
			admin.POST("/order/manual-dock", handler.OrderManualDock)
			admin.POST("/order/sync-progress", handler.OrderSyncProgress)

			// 工单回复
			admin.POST("/gongdan/reply", handler.GongDanReply)
			admin.POST("/gongdan/state", handler.GongDanState)

			// 密价管理
			admin.POST("/myprice", handler.MyPriceSet)
			admin.DELETE("/myprice/:id", handler.MyPriceDelete)

			// 日志
			admin.GET("/log", handler.AdminLogList)

			// 批量删除（硬删）
			admin.POST("/class/batch-delete", handler.ClassBatchDelete)
			admin.POST("/fenlei/batch-delete", handler.FenLeiBatchDelete)
			admin.POST("/dengji/batch-delete", handler.DengjiBatchDelete)
			admin.POST("/huoyuan/batch-delete", handler.HuoYuanBatchDelete)
			admin.POST("/myprice/batch-delete", handler.MyPriceBatchDelete)
			admin.POST("/gongdan/batch-delete", handler.GongDanBatchDelete)
			admin.POST("/users/batch-delete", handler.UsersBatchDelete)
			admin.POST("/order/batch-delete", handler.OrderBatchDelete)
			admin.POST("/log/batch-delete", handler.LogBatchDelete)
			admin.POST("/pay/batch-delete", handler.PayBatchDelete)
		}
	}
}
