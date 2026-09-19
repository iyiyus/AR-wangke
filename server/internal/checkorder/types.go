package checkorder

import "wk-go/internal/model"

// 查课结果中的单个课程
type Course struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// 查课结果
type QueryResult struct {
	Code     int      `json:"code"`
	Msg      string   `json:"msg"`
	Data     []Course `json:"data,omitempty"`
	UserName string   `json:"userName,omitempty"`
}

// 下单结果
type AddResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	YID  string `json:"yid,omitempty"`
}

// 进度查询结果
type ProgressItem struct {
	YID         string `json:"yid"`
	KCName      string `json:"kcname"`
	User        string `json:"user"`
	Pass        string `json:"pass"`
	StatusText  string `json:"status_text"`
	Process     string `json:"process"`
	Remarks     string `json:"remarks"`
	KCStartTime string `json:"kcks"`
	KCEndTime   string `json:"kcjs"`
	KSStartTime string `json:"ksks"`
	KSEndTime   string `json:"ksjs"`
}

// 补刷结果
type RetryResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Adapter 平台对接适配器
type Adapter interface {
	// 查课
	Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult
	// 下单
	Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult
	// 查进度
	Progress(huo *model.HuoYuan, user string) []ProgressItem
	// 补刷
	Retry(huo *model.HuoYuan, yid string) *RetryResult
}

// DetailedAdapter 可选扩展接口：需要订单完整信息(yid/school/kcname/noun)的进度查询
type DetailedAdapter interface {
	ProgressByOrder(huo *model.HuoYuan, o *model.Order) []ProgressItem
}

// DetailedRetryer 可选扩展接口：需要订单完整信息的补刷
type DetailedRetryer interface {
	RetryByOrder(huo *model.HuoYuan, o *model.Order) *RetryResult
}

// 已注册的适配器
var adapters = map[string]Adapter{}

// Register 注册适配器
func Register(pt string, adp Adapter) {
	adapters[pt] = adp
}

// Get 获取适配器
func Get(pt string) Adapter {
	return adapters[pt]
}

// AllPlatforms 所有支持的平台
func AllPlatforms() []string {
	keys := make([]string, 0, len(adapters))
	for k := range adapters {
		keys = append(keys, k)
	}
	return keys
}
