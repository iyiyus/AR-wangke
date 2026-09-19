package checkorder

import "wk-go/internal/model"

// AdapterOuba 欧巴系列：ouba/obks/obsp
type AdapterOuba struct {
	Task string // 1=视频 2=考试 3=普通
}

func (a *AdapterOuba) base(huo *model.HuoYuan) string {
	return ensureScheme(huo.URL, "https")
}

func (a *AdapterOuba) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	// 欧巴是虚假进度，查课直接返回成功并模拟课程
	return &QueryResult{Code: 1, Msg: "查询成功（欧巴）", Data: []Course{{Name: "课程"}}}
}

func (a *AdapterOuba) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	body, err := httpPostForm(a.base(huo)+"/Openapi/submitCourse.jsp", map[string]string{
		"optoken":    huo.Token,
		"type":       noun,
		"task":       a.Task,
		"school":     school,
		"account":    user,
		"pwd":        pass,
		"courseName": kcname,
		"num":        "1",
	}, huo.Cookie)
	if err != nil {
		return &AddResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	msg, _ := m["message"].(string)
	if c, ok := m["code"].(float64); ok && int(c) == 0 {
		return &AddResult{Code: 1, Msg: "添加成功"}
	}
	return &AddResult{Code: -1, Msg: msg}
}

func (a *AdapterOuba) Progress(huo *model.HuoYuan, user string) []ProgressItem {
	return nil
}

// ProgressByOrder 虚假进度（与原 PHP 一致）
func (a *AdapterOuba) ProgressByOrder(huo *model.HuoYuan, o *model.Order) []ProgressItem {
	return []ProgressItem{{
		YID: "8888", KCName: o.KCName, User: o.UserAccount, Pass: o.Pass,
		StatusText: "已完成", Process: "已经加入服务器",
		Remarks: "完成日常任务中,有更新点补刷",
	}}
}

func (a *AdapterOuba) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	return &RetryResult{Code: -1, Msg: "需订单信息补刷"}
}

// RetryByOrder 欧巴重新提交（会扣费）
func (a *AdapterOuba) RetryByOrder(huo *model.HuoYuan, o *model.Order) *RetryResult {
	body, err := httpPostForm(a.base(huo)+"/Openapi/submitCourse.jsp", map[string]string{
		"optoken": huo.Token, "type": o.Noun, "task": a.Task,
		"school": o.School, "account": o.UserAccount, "pwd": o.Pass,
		"courseName": o.KCName, "num": "1",
	}, huo.Cookie)
	if err != nil {
		return &RetryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); ok && int(c) == 0 {
		return &RetryResult{Code: 1, Msg: "添加成功"}
	}
	msg, _ := m["message"].(string)
	return &RetryResult{Code: -1, Msg: msg}
}

func init() {
	Register("ouba", &AdapterOuba{Task: "3"})
	Register("obks", &AdapterOuba{Task: "2"})
	Register("obsp", &AdapterOuba{Task: "1"})
}
