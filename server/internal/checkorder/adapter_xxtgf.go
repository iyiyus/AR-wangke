package checkorder

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"wk-go/internal/model"
)

// AdapterXXTGF 学习通官方查课（仅支持查课，不下单）
type AdapterXXTGF struct{}

var phoneRe = regexp.MustCompile(`^1[3-9]\d{9}$`)

func (a *AdapterXXTGF) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	jar := newCookieJar()
	client := &http.Client{Jar: jar, Timeout: httpClient.Timeout}

	if phoneRe.MatchString(user) {
		// 手机号登录
		form := url.Values{}
		form.Set("uname", user)
		form.Set("code", pass)
		form.Set("loginType", "1")
		form.Set("roleSelect", "true")
		req, _ := http.NewRequest("POST",
			"https://passport2-api.chaoxing.com/v11/loginregister?cx_xxt_passport=json",
			strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		if err != nil {
			return &QueryResult{Code: -1, Msg: "登录失败"}
		}
		resp.Body.Close()
	} else {
		// 学号登录
		schUrl := "https://passport2-api.chaoxing.com/org/searchUnis?filter=" +
			url.QueryEscape(school) + "&product=1&type="
		body, _ := httpGet(schUrl, nil)
		m := parseJSON(body)
		fid := ""
		if froms, ok := m["froms"].([]interface{}); ok && len(froms) > 0 {
			obj, _ := froms[0].(map[string]interface{})
			if v, ok := obj["schoolid"].(string); ok {
				fid = v
			}
		}
		if fid == "" {
			return &QueryResult{Code: -1, Msg: "未找到该学校"}
		}
		form := url.Values{}
		form.Set("pwd", pass)
		form.Set("t", "0")
		req, _ := http.NewRequest("POST",
			"https://passport2-api.chaoxing.com/v6/idNumberLogin?fid="+fid+"&idNumber="+user,
			strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)
		if err != nil {
			return &QueryResult{Code: -1, Msg: "登录失败"}
		}
		resp.Body.Close()
	}

	// 获取课程列表
	resp, err := client.Get("https://mooc1-api.chaoxing.com/mycourse/")
	if err != nil {
		return &QueryResult{Code: -1, Msg: "获取课程失败"}
	}
	defer resp.Body.Close()
	body := readAll(resp.Body)
	m := parseJSON(body)
	if c, ok := m["result"].(float64); !ok || int(c) != 1 {
		return &QueryResult{Code: -1, Msg: "信息错误或者重试"}
	}
	r := &QueryResult{Code: 1, Msg: "查询成功"}
	if list, ok := m["channelList"].([]interface{}); ok {
		for _, item := range list {
			obj, _ := item.(map[string]interface{})
			content, _ := obj["content"].(map[string]interface{})
			course, _ := content["course"].(map[string]interface{})
			if dataList, ok := course["data"].([]interface{}); ok && len(dataList) > 0 {
				dataItem, _ := dataList[0].(map[string]interface{})
				name, _ := dataItem["name"].(string)
				id, _ := dataItem["id"].(string)
				if name != "" {
					r.Data = append(r.Data, Course{ID: id, Name: name})
				}
			}
		}
	}
	return r
}

func (a *AdapterXXTGF) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	return &AddResult{Code: -1, Msg: "学习通官方接口仅支持查课"}
}

func (a *AdapterXXTGF) Progress(huo *model.HuoYuan, user string) []ProgressItem {
	return nil
}

func (a *AdapterXXTGF) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	return &RetryResult{Code: -1, Msg: "不支持"}
}

func init() {
	Register("xxtgf", &AdapterXXTGF{})
}
