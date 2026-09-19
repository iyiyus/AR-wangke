package checkorder

import (
	"strings"
	"wk-go/internal/model"
)

// AdapterHidden 新隐藏接口(00)：查课/下单/进度/补刷
type AdapterHidden struct{}

func (a *AdapterHidden) base(huo *model.HuoYuan) string {
	return ensureScheme(huo.URL, "http")
}

func (a *AdapterHidden) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	body, err := httpPostForm(a.base(huo)+"/querycoursemsg", map[string]string{
		"school": school, "account": user, "password": pass,
		"platform": noun, "token": huo.Token,
	}, huo.Cookie)
	if err != nil {
		return &QueryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); !ok || int(c) != 0 {
		return &QueryResult{Code: -1, Msg: "信息错误或者重试"}
	}
	r := &QueryResult{Code: 1, Msg: "查询成功"}
	if msgObj, ok := m["msg"].(map[string]interface{}); ok {
		if list, ok := msgObj["data"].([]interface{}); ok {
			for _, item := range list {
				obj, _ := item.(map[string]interface{})
				if v, ok := obj["coursename"].(string); ok {
					r.Data = append(r.Data, Course{Name: v})
				}
			}
		}
	}
	return r
}

func (a *AdapterHidden) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	return addTask(huo, noun, school, user, pass, kcname, "/addtask", 0)
}

func (a *AdapterHidden) Progress(huo *model.HuoYuan, user string) []ProgressItem { return nil }

// ProgressByOrder 状态归一化：进行中/已完成/异常
func (a *AdapterHidden) ProgressByOrder(huo *model.HuoYuan, o *model.Order) []ProgressItem {
	body, err := httpPostForm(a.base(huo)+"/getcourseprocess", map[string]string{
		"token": huo.Token, "platform": o.Noun, "school": o.School,
		"account": o.UserAccount, "password": o.Pass, "coursename": o.KCName,
	}, huo.Cookie)
	if err != nil {
		return nil
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); !ok || int(c) != 0 {
		return nil
	}
	list, _ := m["msg"].([]interface{})
	items := make([]ProgressItem, 0, len(list))
	for _, item := range list {
		obj, _ := item.(map[string]interface{})
		kcname := anyToStr(obj["coursename"])
		status := anyToStr(obj["status"])
		per := anyToStr(obj["examstatus"])
		ksks := anyToStr(obj["examstarttime"])
		remarks := anyToStr(obj["extra"])
		status = normalizeHiddenStatus(status)
		items = append(items, ProgressItem{
			YID: "", KCName: kcname, User: o.UserAccount, Pass: o.Pass,
			StatusText: status, Process: per + "  " + remarks, Remarks: per,
			KSStartTime: ksks,
		})
	}
	return items
}

func (a *AdapterHidden) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	return &RetryResult{Code: -1, Msg: "需订单信息补刷"}
}

// RetryByOrder 重新下单
func (a *AdapterHidden) RetryByOrder(huo *model.HuoYuan, o *model.Order) *RetryResult {
	return retryOrder(huo, o, "/order", 0)
}

// AdapterZS 止水接口：查课/下单/进度/补刷
type AdapterZS struct{}

func (a *AdapterZS) base(huo *model.HuoYuan) string {
	return ensureScheme(huo.URL, "http")
}

func (a *AdapterZS) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	body, err := httpPostForm(a.base(huo)+"/querycourse", map[string]string{
		"school": school, "account": user, "password": pass,
		"token": huo.Token, "platform": noun,
	}, huo.Cookie)
	if err != nil {
		return &QueryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); !ok || int(c) != 0 {
		return &QueryResult{Code: -1, Msg: "信息错误或者重试"}
	}
	r := &QueryResult{Code: 1, Msg: "查询成功"}
	if msgObj, ok := m["msg"].(map[string]interface{}); ok {
		if list, ok := msgObj["data"].([]interface{}); ok {
			for _, item := range list {
				obj, _ := item.(map[string]interface{})
				if v, ok := obj["coursename"].(string); ok {
					r.Data = append(r.Data, Course{Name: v})
				}
			}
		}
	}
	return r
}

func (a *AdapterZS) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	return addTask(huo, noun, school, user, pass, kcname, "/order", 1)
}

func (a *AdapterZS) Progress(huo *model.HuoYuan, user string) []ProgressItem { return nil }

// ProgressByOrder 中文键「课程/简洁进度/参考进度」
func (a *AdapterZS) ProgressByOrder(huo *model.HuoYuan, o *model.Order) []ProgressItem {
	body, err := httpPostForm(a.base(huo)+"/getcoursemsg", map[string]string{
		"token": huo.Token, "value": o.Noun, "school": o.School,
		"account": o.UserAccount, "password": o.Pass, "coursename": o.KCName,
	}, huo.Cookie)
	if err != nil {
		return nil
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); !ok || int(c) != 0 {
		return nil
	}
	list, _ := m["msg"].([]interface{})
	items := make([]ProgressItem, 0, len(list))
	for _, item := range list {
		obj, _ := item.(map[string]interface{})
		kcname := anyToStr(obj["课程"])
		status := anyToStr(obj["简洁进度"])
		per := anyToStr(obj["参考进度"])
		ksks := anyToStr(obj["考试开始时间"])
		switch {
		case status == "进行中" || status == "上号中" || status == "执行中" || status == "考试中" ||
			status == "平时分进行中" || status == "平时分中" || status == "重刷中" || status == "待补时长中":
			status = "进行中"
		case strings.Contains(status, "完成") || strings.Contains(status, "h"):
			status = "已完成"
		default:
			status = "异常"
		}
		items = append(items, ProgressItem{
			YID: "", KCName: kcname, User: o.UserAccount, Pass: o.Pass,
			StatusText: status, Process: per,
			KSStartTime: ksks,
		})
	}
	return items
}

func (a *AdapterZS) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	return &RetryResult{Code: -1, Msg: "需订单信息补刷"}
}

// RetryByOrder 止水补单：重新下单 code==1 成功
func (a *AdapterZS) RetryByOrder(huo *model.HuoYuan, o *model.Order) *RetryResult {
	return retryOrder(huo, o, "/order", 1)
}

// addTask 00/zs 下单：00 期望 code==0，zs 期望 code==1
func addTask(huo *model.HuoYuan, noun, school, user, pass, kcname, path string, okCode int) *AddResult {
	base := ensureScheme(huo.URL, "http")
	body, err := httpPostForm(base+path, map[string]string{
		"school": school, "account": user, "password": pass,
		"platform": noun, "token": huo.Token, "course": kcname,
	}, huo.Cookie)
	if err != nil {
		return &AddResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	msg, _ := m["msg"].(string)
	if c, ok := m["code"].(float64); ok && int(c) == okCode {
		return &AddResult{Code: 1, Msg: "添加成功"}
	}
	return &AddResult{Code: -1, Msg: msg}
}

// retryOrder 00/zs 补刷：重新下单
func retryOrder(huo *model.HuoYuan, o *model.Order, path string, okCode int) *RetryResult {
	body, err := httpPostForm(ensureScheme(huo.URL, "http")+path, map[string]string{
		"token": huo.Token, "platform": o.Noun, "school": o.School,
		"account": o.UserAccount, "password": o.Pass, "coursename": o.KCName,
	}, huo.Cookie)
	if err != nil {
		return &RetryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	msg, _ := m["msg"].(string)
	if c, ok := m["code"].(float64); ok && int(c) == okCode {
		return &RetryResult{Code: 1, Msg: msg}
	}
	return &RetryResult{Code: -1, Msg: msg}
}

// normalizeHiddenStatus 00 平台状态归一化
func normalizeHiddenStatus(s string) string {
	switch {
	case s == "进行中" || s == "上号中" || s == "执行中" || s == "考试中" ||
		s == "平时分进行中" || s == "平时分" || s == "重刷中" || s == "待补时长中":
		return "进行中"
	case strings.Contains(s, "完成") || strings.Contains(s, "h"):
		return "已完成"
	case strings.Contains(s, "禁用") || strings.Contains(s, "错误") || strings.Contains(s, "退款"):
		return "异常"
	}
	return "进行中"
}

func init() {
	Register("00", &AdapterHidden{})
	Register("zs", &AdapterZS{})
}
