package checkorder

import (
	"strings"
	"wk-go/internal/model"
)

// Adapter27 兼容 27 协议（标准 PHP 网课系统对接协议）
type Adapter27 struct {
	HTTPS bool
}

func (a *Adapter27) base(huo *model.HuoYuan) string {
	scheme := "http"
	if a.HTTPS {
		scheme = "https"
	}
	return scheme + "://" + strings.TrimPrefix(strings.TrimPrefix(huo.URL, "http://"), "https://") + "/api.php"
}

func (a *Adapter27) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	body, err := httpPostForm(a.base(huo)+"?act=get", map[string]string{
		"uid":      huo.User,
		"key":      huo.Pass,
		"school":   school,
		"user":     user,
		"pass":     pass,
		"platform": noun,
	}, huo.Cookie)
	if err != nil {
		return &QueryResult{Code: -1, Msg: "网络错误: " + err.Error()}
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
			obj, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
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

func (a *Adapter27) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	body, err := httpPostForm(a.base(huo)+"?act=add", map[string]string{
		"uid":      huo.User,
		"key":      huo.Pass,
		"platform": noun,
		"school":   school,
		"user":     user,
		"pass":     pass,
		"kcid":     kcid,
		"kcname":   kcname,
	}, huo.Cookie)
	if err != nil {
		return &AddResult{Code: -1, Msg: "网络错误: " + err.Error()}
	}
	m := parseJSON(body)
	r := &AddResult{}
	if c, ok := m["code"].(float64); ok {
		// 27 协议：code=0 表示成功，转为统一的 code=1
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

func (a *Adapter27) Progress(huo *model.HuoYuan, user string) []ProgressItem {
	body, err := httpPostForm(a.base(huo)+"?act=chadan", map[string]string{
		"username": user,
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
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
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

func (a *Adapter27) Retry(huo *model.HuoYuan, yid string) *RetryResult {
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
	Register("27", &Adapter27{HTTPS: false})
	Register("27s", &Adapter27{HTTPS: true})
}
