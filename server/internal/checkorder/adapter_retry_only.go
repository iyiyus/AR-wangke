package checkorder

import "wk-go/internal/model"

// 仅支持补刷的平台（与原 PHP bsjk.php 一致）

// AdapterNiuwa 牛蛙补刷
type AdapterNiuwa struct{}

func (a *AdapterNiuwa) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	return &QueryResult{Code: -1, Msg: "不支持查课"}
}
func (a *AdapterNiuwa) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	return &AddResult{Code: -1, Msg: "不支持下单"}
}
func (a *AdapterNiuwa) Progress(huo *model.HuoYuan, user string) []ProgressItem { return nil }
func (a *AdapterNiuwa) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	body, err := httpPostForm("https://"+huo.URL+"/api/reprogress.php", map[string]string{
		"token": huo.Token, "id": yid,
	}, huo.Cookie)
	if err != nil {
		return &RetryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); ok && int(c) == 0 {
		return &RetryResult{Code: 1, Msg: "操作成功"}
	}
	return &RetryResult{Code: -1, Msg: "接口异常，请联系管理员"}
}

// AdapterShenwei 神威补刷
type AdapterShenwei struct{}

func (a *AdapterShenwei) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	return &QueryResult{Code: -1, Msg: "不支持查课"}
}
func (a *AdapterShenwei) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	return &AddResult{Code: -1, Msg: "不支持下单"}
}
func (a *AdapterShenwei) Progress(huo *model.HuoYuan, user string) []ProgressItem { return nil }
func (a *AdapterShenwei) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	body, err := httpGet("http://"+huo.URL+"/system/oc/redraw?token="+huo.Token+"&id="+yid, nil)
	if err != nil {
		return &RetryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	if c, ok := m["code"].(float64); ok && int(c) == 0 {
		return &RetryResult{Code: 1, Msg: "操作成功"}
	}
	return &RetryResult{Code: -1, Msg: "接口异常，请联系管理员"}
}

// AdapterYunshang 云上补刷
type AdapterYunshang struct{}

func (a *AdapterYunshang) Query(huo *model.HuoYuan, noun, school, user, pass string) *QueryResult {
	return &QueryResult{Code: -1, Msg: "不支持查课"}
}
func (a *AdapterYunshang) Add(huo *model.HuoYuan, noun, school, user, pass, kcid, kcname string) *AddResult {
	return &AddResult{Code: -1, Msg: "不支持下单"}
}
func (a *AdapterYunshang) Progress(huo *model.HuoYuan, user string) []ProgressItem { return nil }
func (a *AdapterYunshang) Retry(huo *model.HuoYuan, yid string) *RetryResult {
	body, err := httpGet("http://"+huo.URL+"/apis.php?s=reOrderState&token="+huo.Token+"&ids="+yid, nil)
	if err != nil {
		return &RetryResult{Code: -1, Msg: "网络错误"}
	}
	m := parseJSON(body)
	switch c, ok := m["code"].(float64); {
	case ok && int(c) == 0:
		return &RetryResult{Code: 1, Msg: "操作成功"}
	case ok && int(c) == 1:
		return &RetryResult{Code: 1, Msg: "订单完成无需操作，有问题联系管理"}
	}
	return &RetryResult{Code: -1, Msg: "接口异常，请联系管理员"}
}

func init() {
	Register("niuwa", &AdapterNiuwa{})
	Register("shenwei", &AdapterShenwei{})
	Register("yunshang", &AdapterYunshang{})
}
