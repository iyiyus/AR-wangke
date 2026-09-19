package checkorder

import (
	"log"
	"strconv"
	"wk-go/internal/database"
	"wk-go/internal/model"
)

// QueryByClass 根据平台 cid 查课
func QueryByClass(cid int64, school, user, pass string) *QueryResult {
	var class model.Class
	if err := database.DB.First(&class, cid).Error; err != nil {
		return &QueryResult{Code: -1, Msg: "平台不存在"}
	}
	hid, _ := strconv.ParseInt(class.QueryPlat, 10, 64)
	var huo model.HuoYuan
	if err := database.DB.Where("hid = ?", hid).First(&huo).Error; err != nil {
		return &QueryResult{Code: -1, Msg: "查询货源未配置"}
	}
	adp := Get(huo.PT)
	if adp == nil {
		return &QueryResult{Code: -1, Msg: "不支持的平台类型: " + huo.PT}
	}
	return adp.Query(&huo, class.GetNoun, school, user, pass)
}

// AddOrder 根据订单 oid 提交对接
func AddOrder(oid int64) *AddResult {
	var order model.Order
	if err := database.DB.First(&order, oid).Error; err != nil {
		return &AddResult{Code: -1, Msg: "订单不存在"}
	}
	var class model.Class
	if err := database.DB.First(&class, order.CID).Error; err != nil {
		return &AddResult{Code: -1, Msg: "平台不存在"}
	}
	hid, _ := strconv.ParseInt(class.Docking, 10, 64)
	var huo model.HuoYuan
	if err := database.DB.Where("hid = ?", hid).First(&huo).Error; err != nil {
		return &AddResult{Code: -1, Msg: "对接货源未配置"}
	}
	adp := Get(huo.PT)
	if adp == nil {
		return &AddResult{Code: -1, Msg: "不支持的平台类型: " + huo.PT}
	}
	return adp.Add(&huo, class.Noun, order.School, order.UserAccount, order.Pass, order.KCID, order.KCName)
}

// ProgressByOrder 同步订单进度
func ProgressByOrder(oid int64) []ProgressItem {
	var order model.Order
	if err := database.DB.First(&order, oid).Error; err != nil {
		return nil
	}
	var huo model.HuoYuan
	if err := database.DB.Where("hid = ?", order.HID).First(&huo).Error; err != nil {
		return nil
	}
	adp := Get(huo.PT)
	if adp == nil {
		return nil
	}
	// 优先走需要订单完整信息的实现
	if d, ok := adp.(DetailedAdapter); ok {
		return d.ProgressByOrder(&huo, &order)
	}
	return adp.Progress(&huo, order.UserAccount)
}

// RetryOrder 补刷
func RetryOrder(oid int64) *RetryResult {
	var order model.Order
	if err := database.DB.First(&order, oid).Error; err != nil {
		return &RetryResult{Code: -1, Msg: "订单不存在"}
	}
	var huo model.HuoYuan
	if err := database.DB.Where("hid = ?", order.HID).First(&huo).Error; err != nil {
		return &RetryResult{Code: -1, Msg: "对接货源未配置"}
	}
	adp := Get(huo.PT)
	if adp == nil {
		return &RetryResult{Code: -1, Msg: "不支持的平台类型: " + huo.PT}
	}
	if d, ok := adp.(DetailedRetryer); ok {
		return d.RetryByOrder(&huo, &order)
	}
	return adp.Retry(&huo, order.YID)
}

// RemotePlatform 远端平台简略信息
type RemotePlatform struct {
	CID     string `json:"cid"`
	Name    string `json:"name"`
	Price   string `json:"price"`
	Noun    string `json:"noun"`
	FenLei  string `json:"fenlei"`
	Content string `json:"content"`
}

// FetchRemotePlatforms 调用货源的 getclass 拉取该平台的所有商品列表
func FetchRemotePlatforms(huoYuanID int64) ([]RemotePlatform, error) {
	result := []RemotePlatform{}
	var huo model.HuoYuan
	if err := database.DB.Where("hid = ?", huoYuanID).First(&huo).Error; err != nil {
		return result, err
	}
	url := ensureScheme(huo.URL, "http") + "/api.php?act=getclass"
	body, err := httpPostForm(url, map[string]string{
		"uid": huo.User,
		"key": huo.Pass,
	}, huo.Cookie)
	if err != nil {
		log.Printf("[getclass] HTTP 错误: %v", err)
		return result, err
	}
	m := parseJSON(body)
	if !isOK(m["code"]) {
		log.Printf("[getclass] code 校验失败: %v (类型 %T)", m["code"], m["code"])
		return result, nil
	}
	list, ok := m["data"].([]interface{})
	if !ok {
		log.Printf("[getclass] data 不是数组: %T", m["data"])
		return result, nil
	}
	log.Printf("[getclass] 解析到 %d 条", len(list))
	for _, item := range list {
		obj, _ := item.(map[string]interface{})
		p := RemotePlatform{}
		p.CID = anyToStr(obj["cid"])
		p.Name = anyToStr(obj["name"])
		p.Price = anyToStr(obj["price"])
		p.Noun = anyToStr(obj["noun"])
		p.FenLei = anyToStr(obj["fenlei"])
		p.Content = anyToStr(obj["content"])
		if p.CID != "" || p.Name != "" {
			result = append(result, p)
		}
	}
	return result, nil
}

func anyToStr(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		if x == float64(int64(x)) {
			return formatInt(int64(x))
		}
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		if x {
			return "1"
		}
		return "0"
	}
	return ""
}

func isOK(v interface{}) bool {
	switch x := v.(type) {
	case float64:
		return int(x) == 1 || int(x) == 0
	case string:
		return x == "1" || x == "0"
	}
	return false
}

func formatInt(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
