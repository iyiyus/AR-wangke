package checkorder

import (
	"encoding/json"
	"net/http"
	"strings"
	"wk-go/internal/model"
)

// AdapterJT 鸡腿接口（JSON POST）
type AdapterJT struct{}

func (a *AdapterJT) base(huo *model.HuoYuan) string {
	return ensureScheme(huo.URL, "http")
}

func (a *AdapterJT) postJSON(targetURL string, payload map[string]interface{}) ([]byte, error) {
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", targetURL, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := make([]byte, 0, 4096)
	chunk := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(chunk)
		if n > 0 {
			buf = append(buf, chunk[:n]...)
		}
		if err != nil {
			break
		}
	}
	return buf, nil
}

func (a *AdapterJT) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	body, err := a.postJSON(a.base(huo)+"/queryCourses", map[string]interface{}{
		"platform": noun,
		"text":     school + " " + user + " " + pass,
	})
	if err != nil {
		return &QueryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	courses, _ := m["courses"].([]interface{})
	if len(courses) == 0 {
		return &QueryResult{Code: -1, Msg: "信息错误或者重试"}
	}
	first, _ := courses[0].(map[string]interface{})
	children, _ := first["children"].([]interface{})
	r := &QueryResult{Code: 1, Msg: "查询成功"}
	for _, item := range children {
		obj, _ := item.(map[string]interface{})
		if label, ok := obj["label"].(string); ok {
			r.Data = append(r.Data, Course{Name: label})
		}
	}
	return r
}

func (a *AdapterJT) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	body, err := a.postJSON(a.base(huo)+"/submitOrder", map[string]interface{}{
		"platform":   noun,
		"school":     school,
		"username":   user,
		"password":   pass,
		"courseName": kcname,
	})
	if err != nil {
		return &AddResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	msg, _ := m["msg"].(string)
	if c, ok := m["code"].(float64); ok && int(c) == 0 {
		return &AddResult{Code: 1, Msg: "添加成功"}
	}
	return &AddResult{Code: -1, Msg: msg}
}

func (a *AdapterJT) Progress(huo *model.HuoYuan, user string) []ProgressItem {
	return nil
}

func (a *AdapterJT) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	return &RetryResult{Code: -1, Msg: "不支持补刷"}
}

func init() {
	Register("jt", &AdapterJT{})
}
