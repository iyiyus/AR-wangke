package model

import "time"

// User 用户表
type User struct {
	UID      int64   `gorm:"primaryKey;column:uid" json:"uid"`
	UUID     int64   `gorm:"column:uuid" json:"uuid"`
	User     string  `gorm:"column:user" json:"user"`
	Pass     string  `gorm:"column:pass" json:"-"`
	Name     string  `gorm:"column:name" json:"name"`
	CK       int64   `gorm:"column:ck" json:"ck"`
	XDLv     float64 `gorm:"column:xdlv" json:"xdlv"`
	DD       int64   `gorm:"column:dd" json:"dd"`
	QQOpenID string  `gorm:"column:qq_openid" json:"qq_openid"`
	Nickname string  `gorm:"column:nickname" json:"nickname"`
	FaceImg  string  `gorm:"column:faceimg" json:"faceimg"`
	Email    string  `gorm:"column:email" json:"email"`
	Money    float64 `gorm:"column:money" json:"money"`
	ZCZ      string  `gorm:"column:zcz" json:"zcz"`
	AddPrice float64 `gorm:"column:addprice" json:"addprice"`
	Key      string  `gorm:"column:key" json:"key"`
	YQM      string  `gorm:"column:yqm" json:"yqm"`
	YQPrice  string  `gorm:"column:yqprice" json:"yqprice"`
	Notice   string  `gorm:"column:notice" json:"notice"`
	AddTime  string  `gorm:"column:addtime" json:"addtime"`
	EndTime  string  `gorm:"column:endtime" json:"endtime"`
	IP       string  `gorm:"column:ip" json:"ip"`
	Grade    string  `gorm:"column:grade" json:"grade"`
	Active   string  `gorm:"column:active" json:"active"`
}

func (User) TableName() string { return "qingka_wangke_user" }

// Order 订单表
type Order struct {
	OID                    int64  `gorm:"primaryKey;column:oid" json:"oid"`
	UID                    int64  `gorm:"column:uid" json:"uid"`
	CID                    int64  `gorm:"column:cid" json:"cid"`
	HID                    int64  `gorm:"column:hid" json:"hid"`
	YID                    string `gorm:"column:yid" json:"yid"`
	PTName                 string `gorm:"column:ptname" json:"ptname"`
	School                 string `gorm:"column:school" json:"school"`
	Name                   string `gorm:"column:name" json:"name"`
	UserAccount            string `gorm:"column:user" json:"user"`
	Pass                   string `gorm:"column:pass" json:"pass"`
	Phone                  string `gorm:"column:phone" json:"phone"`
	KCID                   string `gorm:"column:kcid" json:"kcid"`
	KCName                 string `gorm:"column:kcname" json:"kcname"`
	CourseStartTime        string `gorm:"column:courseStartTime" json:"courseStartTime"`
	CourseEndTime          string `gorm:"column:courseEndTime" json:"courseEndTime"`
	ExamStartTime          string `gorm:"column:examStartTime" json:"examStartTime"`
	ExamEndTime            string `gorm:"column:examEndTime" json:"examEndTime"`
	ChapterCount           string `gorm:"column:chapterCount" json:"chapterCount"`
	UnfinishedChapterCount string `gorm:"column:unfinishedChapterCount" json:"unfinishedChapterCount"`
	Cookie                 string `gorm:"column:cookie" json:"-"`
	Fees                   string `gorm:"column:fees" json:"fees"`
	Noun                   string `gorm:"column:noun" json:"noun"`
	MiaoShua               string `gorm:"column:miaoshua" json:"miaoshua"`
	AddTime                string `gorm:"column:addtime" json:"addtime"`
	IP                     string `gorm:"column:ip" json:"ip"`
	DockStatus             string `gorm:"column:dockstatus" json:"dockstatus"`
	LoginStatus            string `gorm:"column:loginstatus" json:"loginstatus"`
	Status                 string `gorm:"column:status" json:"status"`
	Process                string `gorm:"column:process" json:"process"`
	BSNum                  string `gorm:"column:bsnum" json:"bsnum"`
	Remarks                string `gorm:"column:remarks" json:"remarks"`
}

func (Order) TableName() string { return "qingka_wangke_order" }

// Class 网课平台表
type Class struct {
	CID       int64  `gorm:"primaryKey;column:cid" json:"cid"`
	Sort      int    `gorm:"column:sort" json:"sort"`
	Name      string `gorm:"column:name" json:"name"`
	GetNoun   string `gorm:"column:getnoun" json:"getnoun"`
	Noun      string `gorm:"column:noun" json:"noun"`
	Price     string `gorm:"column:price" json:"price"`
	QueryPlat string `gorm:"column:queryplat" json:"queryplat"`
	Docking   string `gorm:"column:docking" json:"docking"`
	YunSuan   string `gorm:"column:yunsuan" json:"yunsuan"`
	Content   string `gorm:"column:content" json:"content"`
	AddTime   string `gorm:"column:addtime" json:"addtime"`
	Status    int    `gorm:"column:status" json:"status"`
	FenLei    string `gorm:"column:fenlei" json:"fenlei"`
}

func (Class) TableName() string { return "qingka_wangke_class" }

// Config 系统配置表
type Config struct {
	V string `gorm:"primaryKey;column:v" json:"v"`
	K string `gorm:"column:k" json:"k"`
}

func (Config) TableName() string { return "qingka_wangke_config" }

// Pay 支付记录表
type Pay struct {
	OID        int64      `gorm:"primaryKey;column:oid" json:"oid"`
	OutTradeNo string     `gorm:"column:out_trade_no" json:"out_trade_no"`
	TradeNo    string     `gorm:"column:trade_no" json:"trade_no"`
	Type       string     `gorm:"column:type" json:"type"`
	UID        int64      `gorm:"column:uid" json:"uid"`
	Num        int        `gorm:"column:num" json:"num"`
	AddTime    *time.Time `gorm:"column:addtime" json:"addtime"`
	EndTime    *time.Time `gorm:"column:endtime" json:"endtime"`
	Name       string     `gorm:"column:name" json:"name"`
	Money      string     `gorm:"column:money" json:"money"`
	IP         string     `gorm:"column:ip" json:"ip"`
	Domain     string     `gorm:"column:domain" json:"domain"`
	Status     int        `gorm:"column:status" json:"status"`
}

func (Pay) TableName() string { return "qingka_wangke_pay" }

// Log 日志表
type Log struct {
	ID      int64     `gorm:"primaryKey;column:id" json:"id"`
	UID     int64     `gorm:"column:uid" json:"uid"`
	Type    string    `gorm:"column:type" json:"type"`
	Text    string    `gorm:"column:text" json:"text"`
	Money   string    `gorm:"column:money" json:"money"`
	SMoney  string    `gorm:"column:smoney" json:"smoney"`
	IP      string    `gorm:"column:ip" json:"ip"`
	AddTime time.Time `gorm:"column:addtime" json:"addtime"`
}

func (Log) TableName() string { return "qingka_wangke_log" }

// Dengji 等级表
type Dengji struct {
	ID     int64   `gorm:"primaryKey;column:id" json:"id"`
	Sort   string  `gorm:"column:sort" json:"sort"`
	Name   string  `gorm:"column:name" json:"name"`
	Rate   float64 `gorm:"column:rate" json:"rate"`
	Money  float64 `gorm:"column:money" json:"money"`
	AddKF  string  `gorm:"column:addkf" json:"addkf"`
	GJKF   string  `gorm:"column:gjkf" json:"gjkf"`
	Status string  `gorm:"column:status" json:"status"`
	Time   string  `gorm:"column:time" json:"time"`
}

func (Dengji) TableName() string { return "qingka_wangke_dengji" }

// FenLei 分类表
type FenLei struct {
	ID     int64  `gorm:"primaryKey;column:id" json:"id"`
	Sort   string `gorm:"column:sort" json:"sort"`
	Name   string `gorm:"column:name" json:"name"`
	Status string `gorm:"column:status" json:"status"`
	Time   string `gorm:"column:time" json:"time"`
}

func (FenLei) TableName() string { return "qingka_wangke_fenlei" }

// GongDan 工单表
type GongDan struct {
	GID     int64  `gorm:"primaryKey;column:gid" json:"gid"`
	UID     int64  `gorm:"column:uid" json:"uid"`
	OID     int64  `gorm:"column:oid" json:"oid"`
	Region  string `gorm:"column:region" json:"region"`
	Title   string `gorm:"column:title" json:"title"`
	Content string `gorm:"column:content" json:"content"`
	Answer  string `gorm:"column:answer" json:"answer"`
	State   string `gorm:"column:state" json:"state"`
	AddTime string `gorm:"column:addtime" json:"addtime"`
}

func (GongDan) TableName() string { return "qingka_wangke_gongdan" }

// MiJia 密价表
type MiJia struct {
	MID     int64  `gorm:"primaryKey;column:mid" json:"mid"`
	UID     int64  `gorm:"column:uid" json:"uid"`
	CID     int64  `gorm:"column:cid" json:"cid"`
	Mode    int    `gorm:"column:mode" json:"mode"`
	Price   string `gorm:"column:price" json:"price"`
	AddTime string `gorm:"column:addtime" json:"addtime"`
}

func (MiJia) TableName() string { return "qingka_wangke_mijia" }

// HuoYuan 货源表
type HuoYuan struct {
	HID     int64  `gorm:"primaryKey;column:hid" json:"hid"`
	PT      string `gorm:"column:pt" json:"pt"`
	Name    string `gorm:"column:name" json:"name"`
	URL     string `gorm:"column:url" json:"url"`
	User    string `gorm:"column:user" json:"user"`
	Pass    string `gorm:"column:pass" json:"pass"`
	Token   string `gorm:"column:token" json:"token"`
	IP      string `gorm:"column:ip" json:"ip"`
	Cookie  string `gorm:"column:cookie" json:"cookie"`
	Money   string `gorm:"column:money" json:"money"`
	Status  string `gorm:"column:status" json:"status"`
	AddTime string `gorm:"column:addtime" json:"addtime"`
	EndTime string `gorm:"column:endtime" json:"endtime"`
}

func (HuoYuan) TableName() string { return "qingka_wangke_huoyuan" }

// SxdkOrder 实习打卡订单表（daka 对接，对应 PHP qingka_wangke_sxdk）
type SxdkOrder struct {
	ID          int64  `gorm:"primaryKey;column:id" json:"id"`
	SxdkID      int64  `gorm:"column:sxdkId" json:"sxdkId"`
	UID         int64  `gorm:"column:uid" json:"uid"`
	Platform    string `gorm:"column:platform" json:"platform"`
	Phone       string `gorm:"column:phone" json:"phone"`
	Password    string `gorm:"column:password" json:"password"`
	Code        int    `gorm:"column:code" json:"code"`
	WxPush      string `gorm:"column:wxpush" json:"wxpush"`
	Name        string `gorm:"column:name" json:"name"`
	Address     string `gorm:"column:address" json:"address"`
	UpCheckTime string `gorm:"column:up_check_time" json:"up_check_time"`
	DownCheck   string `gorm:"column:down_check_time" json:"down_check_time"`
	CheckWeek   string `gorm:"column:check_week" json:"check_week"`
	EndTime     string `gorm:"column:end_time" json:"end_time"`
	DayPaper    int    `gorm:"column:day_paper" json:"day_paper"`
	WeekPaper   int    `gorm:"column:week_paper" json:"week_paper"`
	MonthPaper  int    `gorm:"column:month_paper" json:"month_paper"`
	CreateTime  string `gorm:"column:createTime" json:"createTime"`
	UpdateTime  string `gorm:"column:updateTime" json:"updateTime"`
}

func (SxdkOrder) TableName() string { return "qingka_wangke_sxdk" }
