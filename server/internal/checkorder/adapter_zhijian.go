package checkorder

import "wk-go/internal/model"

// AdapterZhijian 指尖接口
type AdapterZhijian struct{}

func (a *AdapterZhijian) base(huo *model.HuoYuan) string {
	return ensureScheme(huo.URL, "http")
}

func (a *AdapterZhijian) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	body, err := httpPostForm(a.base(huo)+"/api/Course/search", map[string]string{
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
	if state, ok := m["state"].(float64); !ok || int(state) != 200 {
		msg, _ := m["msg"].(string)
		return &QueryResult{Code: -1, Msg: msg}
	}
	r := &QueryResult{Code: 1, Msg: "查询成功"}
	if data, ok := m["data"].(map[string]interface{}); ok {
		if list, ok := data["list"].([]interface{}); ok {
			for _, item := range list {
				obj, _ := item.(map[string]interface{})
				if v, ok := obj["name"].(string); ok {
					r.Data = append(r.Data, Course{Name: v})
				}
			}
		}
	}
	return r
}

func (a *AdapterZhijian) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	body, err := httpPostForm(a.base(huo)+"/api/Orders/add", map[string]string{
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
	if state, ok := m["state"].(float64); ok && int(state) == 200 {
		r := &AddResult{Code: 1, Msg: "添加成功"}
		if data, ok := m["data"].(map[string]interface{}); ok {
			if id, ok := data["id"].(string); ok {
				r.YID = id
			}
		}
		return r
	}
	return &AddResult{Code: -1, Msg: msg}
}

func (a *AdapterZhijian) Progress(huo *model.HuoYuan, user string) []ProgressItem {
	return nil
}

func (a *AdapterZhijian) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	body, err := httpPostForm(a.base(huo)+"/api/Orders/refresh/", map[string]string{
		"token": huo.Token,
		"id":    yid,
	}, huo.Cookie)
	if err != nil {
		return &RetryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	if state, ok := m["state"].(float64); ok && int(state) == 200 {
		return &RetryResult{Code: 1, Msg: "操作成功"}
	}
	return &RetryResult{Code: -1, Msg: "接口异常，请联系管理员"}
}

func init() {
	Register("zijian", &AdapterZhijian{})
}
