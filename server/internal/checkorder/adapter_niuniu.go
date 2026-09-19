package checkorder

import "wk-go/internal/model"

// AdapterNiuNiu 牛牛接口
type AdapterNiuNiu struct{}

func (a *AdapterNiuNiu) base(huo *model.HuoYuan) string {
	return ensureScheme(huo.URL, "http")
}

func (a *AdapterNiuNiu) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	body, err := httpPostForm(a.base(huo)+"/query", map[string]string{
		"token":    huo.Token,
		"plat":     noun,
		"school":   school,
		"username": user,
		"password": pass,
	}, huo.Cookie)
	if err != nil {
		return &QueryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); !ok || int(c) != 0 {
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
	return r
}

func (a *AdapterNiuNiu) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	body, err := httpPostForm(a.base(huo)+"/order", map[string]string{
		"token":      huo.Token,
		"plat":       noun,
		"school":     school,
		"username":   user,
		"password":   pass,
		"courseName": kcname,
		"courseId":   kcid,
	}, huo.Cookie)
	if err != nil {
		return &AddResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	msg, _ := m["msg"].(string)
	r := &AddResult{Msg: msg}
	if c, ok := m["code"].(float64); ok && int(c) == 0 {
		r.Code = 1
		if id, ok := m["id"].(string); ok {
			r.YID = id
		}
	} else {
		r.Code = -1
	}
	return r
}

func (a *AdapterNiuNiu) Progress(huo *model.HuoYuan, user string) []ProgressItem {
	return nil
}

func (a *AdapterNiuNiu) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	body, err := httpPostForm(a.base(huo)+"/redraw", map[string]string{
		"token": huo.Token,
		"id":    yid,
	}, huo.Cookie)
	if err != nil {
		return &RetryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); ok && int(c) == 0 {
		return &RetryResult{Code: 1, Msg: "操作成功"}
	}
	return &RetryResult{Code: -1, Msg: "操作失败"}
}

func init() {
	Register("niuniu", &AdapterNiuNiu{})
}
