package checkorder

import (
	"strings"
	"wk-go/internal/model"
)

// Adapter29 兼容 29 系统对接协议（week 类型）
//
// 参照 29 系统 checkorder 对接代码：
//   - ckjk.php：查课   POST {url}/api.php?act=get       参数 uid/key/school/user/pass/platform/kcid
//   - xdjk.php：下单   POST {url}/api.php?act=addweek   参数 uid/key/platform/school/user/pass/kcname/kcid
//   - xdjk.php：查单   POST {url}/api.php?act=chadan2   参数 yid（新版，按订单ID精确查询）
//     POST {url}/api/search?           参数 username/kcname/cid（新版按账号查询）
//   - bsjk.php：补刷   POST {url}/api.php?act=budan     参数 uid/key/id
//   - ztjk.php：状态   POST {url}/api.php?act=zt        参数 uid/key/id
//   - xgjk.php：改密   POST {url}/api.php?act=xgmm      参数 uid/key/id/xgmm
//   - msjk.php：秒杀   POST {url}/api.php?act=ms        参数 uid/key/id
//
// 成功约定：act=addweek 返回 code=0 表示成功（统一为 code=1）；
//
//	查单返回 {code:1, data:[{id,kcname,status,process,remarks,courseStartTime,courseEndTime,examStartTime,examEndTime}]}。
type Adapter29 struct {
	HTTPS bool
}

func (a *Adapter29) base(huo *model.HuoYuan) string {
	scheme := "http"
	if a.HTTPS {
		scheme = "https"
	}
	return scheme + "://" + strings.TrimPrefix(strings.TrimPrefix(huo.URL, "http://"), "https://") + "/api.php"
}

// Query 查课 act=get
func (a *Adapter29) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
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

// Add 下单 act=addweek（成功 code=0 → 统一 code=1）
func (a *Adapter29) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	body, err := httpPostForm(a.base(huo)+"?act=addweek", map[string]string{
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

// Progress 按账号查单（act=chadan 兜底，username 查询）
func (a *Adapter29) Progress(huo *model.HuoYuan, user string) []ProgressItem {
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

// ProgressByOrder 新版查单：优先 /api/search（username/kcname/cid），失败退回 chadan
func (a *Adapter29) ProgressByOrder(huo *model.HuoYuan, o *model.Order) []ProgressItem {
	searchURL := strings.TrimSuffix(strings.TrimSuffix(huo.URL, "/"), "api.php")
	searchURL = strings.TrimRight(searchURL, "/") + "/api/search?"
	body, err := httpPostForm(searchURL, map[string]string{
		"username": o.UserAccount,
		"kcname":   o.KCName,
		"cid":      o.Noun,
	}, huo.Cookie)
	if err == nil {
		if items := parseProgressList(body, o.UserAccount); len(items) > 0 {
			return items
		}
	}
	return a.Progress(huo, o.UserAccount)
}

// ProgressByYID 按 yid 精确查进度（29 新版 act=chadan2）
func (a *Adapter29) ProgressByYID(huo *model.HuoYuan, yid string) []ProgressItem {
	body, err := httpPostForm(a.base(huo)+"?act=chadan2", map[string]string{
		"yid": yid,
	}, huo.Cookie)
	if err != nil {
		return nil
	}
	return parseProgressList(body, "")
}

// Retry 补刷 act=budan
func (a *Adapter29) Retry(huo *model.HuoYuan, yid string) *RetryResult {
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

// Zt 状态查询 act=zt
func (a *Adapter29) Zt(huo *model.HuoYuan, yid string) *RetryResult {
	body, err := httpPostForm(a.base(huo)+"?act=zt", map[string]string{
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

// Xgmm 修改密码 act=xgmm
func (a *Adapter29) Xgmm(huo *model.HuoYuan, yid, newpass string) *RetryResult {
	body, err := httpPostForm(a.base(huo)+"?act=xgmm", map[string]string{
		"uid":  huo.User,
		"key":  huo.Pass,
		"id":   yid,
		"xgmm": newpass,
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

// Ms 秒杀 act=ms
func (a *Adapter29) Ms(huo *model.HuoYuan, yid string) *RetryResult {
	body, err := httpPostForm(a.base(huo)+"?act=ms", map[string]string{
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

// parseProgressList 解析查单返回 {code:1,data:[...]}
func parseProgressList(body []byte, user string) []ProgressItem {
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

func init() {
	Register("29", &Adapter29{HTTPS: false})
	Register("29s", &Adapter29{HTTPS: true})
}
