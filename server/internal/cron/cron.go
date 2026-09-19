package cron

import (
	"log"
	"time"
	"wk-go/internal/checkorder"
	"wk-go/internal/database"
	"wk-go/internal/model"
)

// Start 启动定时任务
func Start() {
	go runDocking()
	go runProgress()
}

// 自动对接：每 30 秒处理一批 dockstatus=0 的订单
func runDocking() {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for range t.C {
		dockPending()
	}
}

func dockPending() {
	var orders []model.Order
	database.DB.Where("dockstatus = ? AND status != ?", "0", "已取消").
		Limit(10).Find(&orders)
	for _, order := range orders {
		// 简单参数校验
		if order.UserAccount == "" || order.School == "" {
			database.DB.Model(&order).Updates(map[string]interface{}{
				"status": "请检查账号或学校", "dockstatus": "2",
			})
			continue
		}
		res := checkorder.AddOrder(order.OID)
		if res.Code == 1 {
			var class model.Class
			database.DB.First(&class, order.CID)
			database.DB.Model(&order).Updates(map[string]interface{}{
				"hid":        class.Docking,
				"status":     "进行中",
				"dockstatus": "1",
				"yid":        res.YID,
			})
			log.Printf("[cron] 订单 %d 对接成功", order.OID)
		} else {
			database.DB.Model(&order).Updates(map[string]interface{}{
				"dockstatus": "2",
				"remarks":    res.Msg,
			})
			log.Printf("[cron] 订单 %d 对接失败: %s", order.OID, res.Msg)
		}
	}
}

// 自动同步进度：每 2 分钟同步一次进行中/补刷中订单
func runProgress() {
	t := time.NewTicker(2 * time.Minute)
	defer t.Stop()
	for range t.C {
		syncProgress()
	}
}

func syncProgress() {
	var orders []model.Order
	database.DB.Where("dockstatus = ? AND status IN ?", "1",
		[]string{"进行中", "补刷中", "待处理"}).
		Limit(20).Find(&orders)

	for _, order := range orders {
		var items []checkorder.ProgressItem
		// xy/29 协议有 yid 时走精确查询（act=chadan2）
		var huo model.HuoYuan
		if database.DB.Where("hid = ?", order.HID).First(&huo).Error == nil &&
			order.YID != "" {
			if adp, ok := checkorder.Get(huo.PT).(*checkorder.AdapterXY); ok {
				items = adp.ProgressByYID(&huo, order.YID)
			}
			if items == nil {
				if adp29, ok := checkorder.Get(huo.PT).(*checkorder.Adapter29); ok {
					items = adp29.ProgressByYID(&huo, order.YID)
				}
			}
		}
		if items == nil {
			items = checkorder.ProgressByOrder(order.OID)
		}
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
				})
		}
	}
}
