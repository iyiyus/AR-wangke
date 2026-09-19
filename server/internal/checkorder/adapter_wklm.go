package checkorder

import (
	"encoding/json"
	"wk-go/internal/model"
)

// AdapterWKLM 网课联盟
type AdapterWKLM struct {
	Tag string // wklm 或 wklmbks
}

func (a *AdapterWKLM) base(huo *model.HuoYuan) string {
	return ensureScheme(huo.URL, "http")
}

func (a *AdapterWKLM) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	body, err := httpPostForm(a.base(huo)+"/order/search", map[string]string{
		"ptid":          noun,
		"school":        school,
		"studentnumber": user,
		"studentpwd":    pass,
	}, huo.Cookie)
	if err != nil {
		return &QueryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); !ok || int(c) != 1 {
		return &QueryResult{Code: -1, Msg: "信息错误或者重试"}
	}
	r := &QueryResult{Code: 1, Msg: "查询成功"}
	data, _ := m["data"].(map[string]interface{})
	list, _ := data["courseList"].([]interface{})
	for _, item := range list {
		obj, _ := item.(map[string]interface{})
		if v, ok := obj["name"].(string); ok {
			r.Data = append(r.Data, Course{Name: v})
		}
	}
	return r
}

func (a *AdapterWKLM) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	accountInfo := "考试"
	if a.Tag == "wklmbks" {
		accountInfo = "不考试"
	}
	if noun == "12" {
		kcname = "(职教云)" + kcname
	}
	body, err := httpPostFormHeaders(a.base(huo)+"/orderrecord/add?dialog=1", map[string]string{
		"row[courses_id]_text": "",
		"row[courses_id]":      noun,
		"row[school]":          school,
		"row[studentnumber]":   user,
		"row[studentpwd]":      pass,
		"courseids":            kcid,
		"row[name]":            kcname,
		"row[accountinfo]":     accountInfo,
		"row[hour]":            "0",
	}, huo.Cookie, map[string]string{"X-Requested-With": "XMLHttpRequest"})
	if err != nil {
		return &AddResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	msg, _ := m["message"].(string)
	if c, ok := m["code"].(float64); ok && int(c) == 1 {
		return &AddResult{Code: 1, Msg: "添加成功"}
	}
	return &AddResult{Code: -1, Msg: msg}
}

func (a *AdapterWKLM) Progress(huo *model.HuoYuan, user string) []ProgressItem {
	return nil
}

// ProgressByOrder 按学号查询联盟订单进度
func (a *AdapterWKLM) ProgressByOrder(huo *model.HuoYuan, o *model.Order) []ProgressItem {
	filter, _ := json.Marshal(map[string]string{"studentnumber": o.UserAccount})
	body, err := httpPostFormHeaders(a.base(huo)+"/orderrecord/index", map[string]string{
		"addtabs": "1",
		"sort":    "id",
		"order":   "desc",
		"offset":  "0",
		"limit":   "10",
		"filter":  string(filter),
	}, huo.Cookie, map[string]string{"X-Requested-With": "XMLHttpRequest"})
	if err != nil {
		return nil
	}
	m := parseJSON(body)
	if _, ok := m["total"]; !ok {
		return nil
	}
	rows, _ := m["rows"].([]interface{})
	items := make([]ProgressItem, 0, len(rows))
	for _, row := range rows {
		obj, _ := row.(map[string]interface{})
		yid := anyToStr(obj["id"])
		yuser := anyToStr(obj["studentnumber"])
		ypass := anyToStr(obj["studentpwd"])
		kcname := anyToStr(obj["name"])
		ksks := anyToStr(obj["ksks"])
		ksjs := anyToStr(obj["ksjs"])
		process := anyToStr(obj["process"])
		statusTex := "进行中"
		if s, ok := obj["status"].(float64); ok && int(s) == 2 {
			statusTex = "已完成"
		}
		items = append(items, ProgressItem{
			YID: yid, KCName: kcname, User: yuser, Pass: ypass,
			StatusText: statusTex, Process: statusTex, Remarks: process,
			KSStartTime: ksks, KSEndTime: ksjs,
		})
	}
	return items
}

func (a *AdapterWKLM) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	body, err := httpPostFormHeaders(a.base(huo)+"/orderrecord/retry3/ids/"+yid, map[string]string{},
		huo.Cookie, map[string]string{"X-Requested-With": "XMLHttpRequest"})
	if err != nil {
		return &RetryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); ok && int(c) == 1 {
		return &RetryResult{Code: 1, Msg: "操作成功"}
	}
	return &RetryResult{Code: -1, Msg: "接口异常，请联系管理员"}
}

func init() {
	Register("wklm", &AdapterWKLM{Tag: "wklm"})
	Register("wklmbks", &AdapterWKLM{Tag: "wklmbks"})
}
