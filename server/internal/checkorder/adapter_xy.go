package checkorder

import (
	"strings"
	"wk-go/internal/model"
)

// AdapterXY 小月接口（29 协议）
type AdapterXY struct{}

func (a *AdapterXY) base(huo *model.HuoYuan) string {
	url := strings.TrimPrefix(strings.TrimPrefix(huo.URL, "http://"), "https://")
	scheme := "http"
	if strings.HasPrefix(huo.URL, "https://") {
		scheme = "https"
	}
	return scheme + "://" + url + "/api.php"
}

func (a *AdapterXY) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	body, err := httpPostForm(a.base(huo)+"?act=get", map[string]string{
		"uid":      huo.User,
		"key":      huo.Pass,
		"school":   school,
		"user":     user,
		"pass":     pass,
		"platform": noun,
		"kcid":     "",
	}, huo.Cookie)
	if err != nil {
		return &QueryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	r := &QueryResult{}
	if c, ok := m["code"].(float64); ok {
		r.Code = int(c)
	}
	if msg, ok := m["msg"].(string); ok {
		r.Msg = msg
	}
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

// Add 下单（29 协议：act=addweek，code=0 成功）
func (a *AdapterXY) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	body, err := httpPostForm(a.base(huo)+"?act=addweek", map[string]string{
		"uid":      huo.User,
		"key":      huo.Pass,
		"platform": noun,
		"school":   school,
		"user":     user,
		"pass":     pass,
		"kcname":   kcname,
		"kcid":     kcid,
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
	if id, ok := m["id"].(string); ok {
		r.YID = id
	}
	return r
}

// ProgressByYID xy 协议特有：根据 yid 精确查进度
func (a *AdapterXY) ProgressByYID(huo *model.HuoYuan, yid string) []ProgressItem {
	body, err := httpPostForm(a.base(huo)+"?act=chadan2", map[string]string{
		"yid": yid,
	}, huo.Cookie)
	if err != nil {
		return nil
	}
	return a.parseProgress(body, "")
}

func (a *AdapterXY) Progress(huo *model.HuoYuan, user string) []ProgressItem {
	body, err := httpPostForm(a.base(huo)+"?act=chadan", map[string]string{
		"username": user,
	}, huo.Cookie)
	if err != nil {
		return nil
	}
	return a.parseProgress(body, user)
}

func (a *AdapterXY) parseProgress(body []byte, user string) []ProgressItem {
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

func (a *AdapterXY) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	body, err := httpPostForm(a.base(huo)+"?act=budan", map[string]string{
		"uid": huo.User,
		"key": huo.Pass,
		"id":  yid,
	}, huo.Cookie)
	if err != nil {
		return &RetryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	r := &RetryResult{}
	if c, ok := m["code"].(float64); ok {
		r.Code = int(c)
	}
	if msg, ok := m["msg"].(string); ok {
		r.Msg = msg
	}
	return r
}

func init() {
	Register("xy", &AdapterXY{})
}
