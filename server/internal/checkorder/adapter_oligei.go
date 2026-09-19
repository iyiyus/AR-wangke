package checkorder

import "wk-go/internal/model"

// AdapterOligei oligei/benz 接口
type AdapterOligei struct{}

func (a *AdapterOligei) base(huo *model.HuoYuan) string {
	return ensureScheme(huo.URL, "http")
}

func (a *AdapterOligei) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	body, err := httpPostForm(a.base(huo)+"/api/query", map[string]string{
		"token":  huo.Token,
		"school": school,
		"user":   user,
		"pass":   pass,
		"ptid":   noun,
	}, huo.Cookie)
	if err != nil {
		return &QueryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); ok && int(c) == -1 {
		msg, _ := m["msg"].(string)
		return &QueryResult{Code: -1, Msg: msg}
	}
	r := &QueryResult{Code: 1, Msg: "查询成功"}
	if list, ok := m["data"].([]interface{}); ok {
		for _, item := range list {
			obj, _ := item.(map[string]interface{})
			c := Course{}
			if v, ok := obj["id"].(string); ok {
				c.ID = v
			}
			if v, ok := obj["name"].(string); ok {
				c.Name = v
			}
			r.Data = append(r.Data, c)
		}
	}
	if v, ok := m["userName"].(string); ok {
		r.UserName = v
	}
	return r
}

func (a *AdapterOligei) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	body, err := httpPostForm(a.base(huo)+"/api/add", map[string]string{
		"token":    huo.Token,
		"ptid":     noun,
		"school":   school,
		"user":     user,
		"pass":     pass,
		"kcname":   kcname,
		"kcid":     kcid,
		"miaoshua": "0",
	}, huo.Cookie)
	if err != nil {
		return &AddResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	r := &AddResult{}
	if c, ok := m["code"].(float64); ok {
		if int(c) == 0 {
			r.Code = 1
		} else {
			r.Code = int(c)
		}
	}
	if msg, ok := m["msg"].(string); ok {
		r.Msg = msg
	}
	return r
}

func (a *AdapterOligei) Progress(huo *model.HuoYuan, user string) []ProgressItem {
	body, err := httpPostForm(a.base(huo)+"/api/order", map[string]string{
		"token": huo.Token,
		"user":  user,
	}, huo.Cookie)
	if err != nil {
		return nil
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); !ok || int(c) != 1 {
		return nil
	}
	list, _ := m["data"].([]interface{})
	var items []ProgressItem
	for _, item := range list {
		obj, _ := item.(map[string]interface{})
		p := ProgressItem{User: user}
		if v, ok := obj["id"].(string); ok {
			p.YID = v
		}
		if v, ok := obj["kcname"].(string); ok {
			p.KCName = v
		}
		if v, ok := obj["status"].(string); ok {
			p.StatusText = v
		}
		if v, ok := obj["process"].(string); ok {
			p.Process = v
		}
		if v, ok := obj["remarks"].(string); ok {
			p.Remarks = v
		}
		if v, ok := obj["courseStartTime"].(string); ok {
			p.KCStartTime = v
		}
		if v, ok := obj["courseEndTime"].(string); ok {
			p.KCEndTime = v
		}
		if v, ok := obj["examStartTime"].(string); ok {
			p.KSStartTime = v
		}
		if v, ok := obj["examEndTime"].(string); ok {
			p.KSEndTime = v
		}
		items = append(items, p)
	}
	return items
}

func (a *AdapterOligei) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	body, err := httpPostForm(a.base(huo)+"/api/reset", map[string]string{
		"token": huo.Token,
		"id":    yid,
	}, huo.Cookie)
	if err != nil {
		return &RetryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); ok && int(c) == 1 {
		return &RetryResult{Code: 1, Msg: "操作成功"}
	}
	return &RetryResult{Code: -1, Msg: "操作失败"}
}

func init() {
	Register("oligei", &AdapterOligei{})
}
