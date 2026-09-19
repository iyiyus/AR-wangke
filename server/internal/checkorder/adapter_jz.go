package checkorder

import (
	"strings"
	"wk-go/internal/model"
)

// AdapterJZ 捐赠接口（按平台 noun 拼接路径）
type AdapterJZ struct{}

func (a *AdapterJZ) base(huo *model.HuoYuan) string {
	return ensureScheme(huo.URL, "http")
}

func (a *AdapterJZ) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	body, err := httpPostForm(a.base(huo)+"/api/"+noun, map[string]string{
		"school_name": school,
		"username":    user,
		"password":    pass,
	}, huo.Cookie)
	if err != nil {
		return &QueryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); !ok || int(c) != 0 {
		return &QueryResult{Code: -1, Msg: "信息错误或者重试"}
	}
	r := &QueryResult{Code: 1, Msg: "查询成功"}
	if dataStr, ok := m["data"].(string); ok {
		for _, name := range strings.Split(dataStr, "\n") {
			if name = strings.TrimSpace(name); name != "" {
				r.Data = append(r.Data, Course{Name: name})
			}
		}
	}
	return r
}

func (a *AdapterJZ) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	return &AddResult{Code: -1, Msg: "捐赠接口仅支持查课"}
}

func (a *AdapterJZ) Progress(huo *model.HuoYuan, user string) []ProgressItem {
	return nil
}

func (a *AdapterJZ) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	return &RetryResult{Code: -1, Msg: "不支持补刷"}
}

func init() {
	Register("jz", &AdapterJZ{})
}
