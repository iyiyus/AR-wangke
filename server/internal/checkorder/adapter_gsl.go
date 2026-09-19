package checkorder

import "wk-go/internal/model"

// AdapterGSL 哥斯拉 2.0
type AdapterGSL struct{}

func (a *AdapterGSL) base(huo *model.HuoYuan) string {
	return ensureScheme(huo.URL, "http")
}

func (a *AdapterGSL) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	body, err := httpPostForm(a.base(huo)+"/api/unionApi/course", map[string]string{
		"Token":    huo.Token,
		"school":   school,
		"username": user,
		"password": pass,
		"plat":     noun,
	}, huo.Cookie)
	if err != nil {
		return &QueryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	msg, _ := m["msg"].(string)
	if msg != "success" {
		return &QueryResult{Code: -1, Msg: msg}
	}
	r := &QueryResult{Code: 1, Msg: "查询成功"}
	if list, ok := m["data"].([]interface{}); ok {
		for _, item := range list {
			obj, _ := item.(map[string]interface{})
			if v, ok := obj["name"].(string); ok {
				r.Data = append(r.Data, Course{Name: v})
			}
		}
	}
	return r
}

func (a *AdapterGSL) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	body, err := httpPostForm(a.base(huo)+"/api/unionApi/simple", map[string]string{
		"Token":    huo.Token,
		"plat":     noun,
		"school":   school,
		"username": user,
		"password": pass,
		"name":     kcname,
		"mode":     "1",
	}, huo.Cookie)
	if err != nil {
		return &AddResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	msg, _ := m["msg"].(string)
	r := &AddResult{Msg: msg}
	if c, ok := m["code"].(float64); ok && int(c) == 0 {
		r.Code = 1
		if ids, ok := m["course_trade_ids"].([]interface{}); ok && len(ids) > 0 {
			if id, ok := ids[0].(string); ok {
				r.YID = id
			}
		}
	} else {
		r.Code = -1
	}
	return r
}

func (a *AdapterGSL) Progress(huo *model.HuoYuan, user string) []ProgressItem {
	return nil
}

// ProgressByOrder 按 yid 查询单个订单进度
func (a *AdapterGSL) ProgressByOrder(huo *model.HuoYuan, o *model.Order) []ProgressItem {
	if o.YID == "" {
		return nil
	}
	body, err := httpPostForm(a.base(huo)+"/api/unionApi/status/"+o.YID, map[string]string{
		"Token": huo.Token,
	}, huo.Cookie)
	if err != nil {
		return nil
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); !ok || int(c) != 200 {
		return nil
	}
	data, _ := m["data"].(map[string]interface{})
	desc, _ := data["desc"].(string)
	suc, _ := data["suc_progress"].(string)
	remarks, _ := data["remarks"].(string)
	exam, _ := data["exam"].(string)
	if exam != "" {
		remarks = remarks + " | 考试分数：" + exam
	}
	return []ProgressItem{{
		YID:        o.YID,
		KCName:     o.KCName,
		User:       o.UserAccount,
		Pass:       o.Pass,
		StatusText: desc,
		Process:    suc + "%",
		Remarks:    remarks,
	}}
}

func (a *AdapterGSL) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	body, err := httpPostForm(a.base(huo)+"/api/unionApi/retry/"+yid, map[string]string{
		"Token": huo.Token,
	}, huo.Cookie)
	if err != nil {
		return &RetryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	msg, _ := m["msg"].(string)
	if c, ok := m["code"].(float64); ok && int(c) == 200 {
		return &RetryResult{Code: 1, Msg: msg}
	}
	return &RetryResult{Code: -1, Msg: msg}
}

func init() {
	Register("gsl2.0", &AdapterGSL{})
}
