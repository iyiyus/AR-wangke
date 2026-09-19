package handler

import (
	"fmt"
	"strconv"
	"time"
	"wk-go/internal/checkorder"
	"wk-go/internal/database"
	"wk-go/internal/model"

	"github.com/gin-gonic/gin"
)

// GET /api/class/list
func ClassList(c *gin.Context) {
	uid, _ := c.Get("uid")
	uidInt := uid.(int64)

	page, _ := atoiSafe(c.DefaultQuery("current", c.DefaultQuery("page", "1")))
	pageSize, _ := atoiSafe(c.DefaultQuery("size", c.DefaultQuery("pageSize", "50")))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	status := c.Query("status")
	fenlei := c.Query("fenlei")

	db := database.DB.Model(&model.Class{})
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if fenlei != "" {
		db = db.Where("fenlei = ?", fenlei)
	}

	var total int64
	db.Count(&total)

	var classes []model.Class
	db.Order("sort ASC, cid ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&classes)

	// 附加用户密价
	type classWithPrice struct {
		model.Class
		MyPrice string `json:"my_price"`
		MyMode  int    `json:"my_mode"`
	}
	result := make([]classWithPrice, 0, len(classes))
	for _, cls := range classes {
		item := classWithPrice{Class: cls}
		var mj model.MiJia
		if err := database.DB.Where("uid = ? AND cid = ?", uidInt, cls.CID).First(&mj).Error; err == nil {
			item.MyPrice = mj.Price
			item.MyMode = mj.Mode
		}
		result = append(result, item)
	}
	PageResult(c, result, total, page, pageSize)
}

// GET /api/class/all  仅返回上架平台（下单用）
func ClassAll(c *gin.Context) {
	var classes []model.Class
	database.DB.Where("status = 1").Order("sort ASC").Find(&classes)
	OK(c, classes)
}

// POST /api/class/add  管理员
func ClassAdd(c *gin.Context) {
	var req model.Class
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	req.AddTime = time.Now().Format("2006-01-02 15:04:05")
	if err := database.DB.Create(&req).Error; err != nil {
		Fail(c, "添加失败")
		return
	}
	OK(c, req)
}

// PUT /api/class/:id  管理员
func ClassUpdate(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	database.DB.Model(&model.Class{}).Where("cid = ?", id).Updates(req)
	OK(c, nil)
}

// DELETE /api/class/:id  管理员
func ClassDelete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&model.Class{}, id)
	OK(c, nil)
}

// GET /api/fenlei/list
func FenLeiList(c *gin.Context) {
	var list []model.FenLei
	database.DB.Where("status = 1").Order("sort ASC").Find(&list)
	OK(c, list)
}

// POST /api/admin/fenlei  管理员
func FenLeiAdd(c *gin.Context) {
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		Fail(c, "参数错误: " + err.Error())
		return
	}
	toStr := func(k string) string {
		if v, ok := raw[k]; ok && v != nil { return fmt.Sprintf("%v", v) }
		return ""
	}
	toInt64 := func(k string) int64 {
		if v, ok := raw[k]; ok && v != nil {
			switch n := v.(type) {
			case float64: return int64(n)
			case int64: return n
			case int: return int64(n)
			}
		}
		return 0
	}
	req := model.FenLei{
		ID:     toInt64("id"),
		Sort:   toStr("sort"),
		Name:   toStr("name"),
		Status: toStr("status"),
		Time:   time.Now().Format("2006-01-02 15:04:05"),
	}
	if req.Status == "" {
		req.Status = "1"
	}
	if req.Sort == "" {
		req.Sort = "0"
	}
	// 如果指定了 ID，检查是否冲突
	if req.ID > 0 {
		var count int64
		database.DB.Model(&model.FenLei{}).Where("id = ?", req.ID).Count(&count)
		if count > 0 {
			Fail(c, "该分类ID已存在")
			return
		}
	}
	if err := database.DB.Create(&req).Error; err != nil {
		Fail(c, "添加失败")
		return
	}
	OK(c, req)
}

// PUT /api/fenlei/:id  管理员
func FenLeiUpdate(c *gin.Context) {
	id := c.Param("id")
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	database.DB.Model(&model.FenLei{}).Where("id = ?", id).Updates(req)
	OK(c, nil)
}

// DELETE /api/fenlei/:id  管理员
func FenLeiDelete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&model.FenLei{}, id)
	OK(c, nil)
}

// GET /api/admin/class/remote?hid=xx  拉取远端平台列表
func ClassRemote(c *gin.Context) {
	hid := c.Query("hid")
	if hid == "" {
		Fail(c, "请指定货源 hid")
		return
	}
	hidInt, _ := strconv.ParseInt(hid, 10, 64)
	list, err := checkorder.FetchRemotePlatforms(hidInt)
	if err != nil {
		Fail(c, "拉取失败: "+err.Error())
		return
	}
	OK(c, list)
}

// POST /api/admin/class/batch-import  批量导入平台
func ClassBatchImport(c *gin.Context) {
	var req struct {
		HID         int64  `json:"hid" binding:"required"`
		FenLei      string `json:"fenlei"`
		Price       string `json:"price"`
		UseRemote   bool   `json:"use_remote"`
		MarkupMode  string `json:"markup_mode"`
		PriceMarkup string `json:"price_markup"`
		Items       []struct {
			CID     string `json:"cid"`
			Name    string `json:"name"`
			Price   string `json:"price"`
			Noun    string `json:"noun"`
			FenLei  string `json:"fenlei"`
			Content string `json:"content"`
		} `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, "参数错误")
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	hidStr := formatInt64(req.HID)
	markup := 0.0
	if req.PriceMarkup != "" {
		if v, err := strconv.ParseFloat(req.PriceMarkup, 64); err == nil {
			markup = v
		}
	}
	if req.MarkupMode == "" {
		req.MarkupMode = "multiply"
	}

	// 一次性查询已存在的同 docking 同名称记录
	names := make([]string, 0, len(req.Items))
	fenLeiSet := map[string]bool{}
	for _, item := range req.Items {
		names = append(names, item.Name)
		if req.UseRemote && item.FenLei != "" {
			fenLeiSet[item.FenLei] = true
		}
	}
	existing := map[string]bool{}
	if len(names) > 0 {
		var rows []struct{ Name string }
		database.DB.Model(&model.Class{}).
			Select("name").
			Where("docking = ? AND name IN ?", hidStr, names).
			Scan(&rows)
		for _, r := range rows {
			existing[r.Name] = true
		}
	}

	// 自动创建缺失的分类（远端模式才需要）
	if len(fenLeiSet) > 0 {
		ids := make([]int64, 0, len(fenLeiSet))
		for k := range fenLeiSet {
			if v, err := strconv.ParseInt(k, 10, 64); err == nil && v > 0 {
				ids = append(ids, v)
			}
		}
		if len(ids) > 0 {
			var existIDs []int64
			database.DB.Model(&model.FenLei{}).
				Where("id IN ?", ids).
				Pluck("id", &existIDs)
			existMap := map[int64]bool{}
			for _, id := range existIDs {
				existMap[id] = true
			}
			toCreate := []model.FenLei{}
			for _, id := range ids {
				if !existMap[id] {
					toCreate = append(toCreate, model.FenLei{
						ID:     id,
						Name:   fmt.Sprintf("分类 %d", id),
						Sort:   "0",
						Status: "1",
						Time:   now,
					})
				}
			}
			if len(toCreate) > 0 {
				database.DB.Create(&toCreate)
			}
		}
	}

	// 构造批量插入数据
	toInsert := make([]model.Class, 0, len(req.Items))
	for _, item := range req.Items {
		if existing[item.Name] {
			continue
		}
		// 价格
		price := req.Price
		if req.UseRemote && item.Price != "" {
			if p, err := strconv.ParseFloat(item.Price, 64); err == nil {
				var final float64
				if req.MarkupMode == "add" {
					final = p + markup
				} else {
					if markup <= 0 {
						markup = 1
					}
					final = p * markup
				}
				if final < 0 {
					final = 0
				}
				price = strconv.FormatFloat(final, 'f', 2, 64)
			} else {
				price = item.Price
			}
		}
		fenlei := req.FenLei
		if req.UseRemote && item.FenLei != "" {
			fenlei = item.FenLei
		}
		noun := item.CID
		if item.Noun != "" {
			noun = item.Noun
		}
		toInsert = append(toInsert, model.Class{
			Sort:      10,
			Name:      item.Name,
			GetNoun:   noun,
			Noun:      noun,
			Price:     price,
			QueryPlat: hidStr,
			Docking:   hidStr,
			YunSuan:   "*",
			Content:   item.Content,
			AddTime:   now,
			Status:    1,
			FenLei:    fenlei,
		})
	}

	// 分批 INSERT，每批 200 条
	imported := 0
	const batchSize = 200
	for i := 0; i < len(toInsert); i += batchSize {
		end := i + batchSize
		if end > len(toInsert) {
			end = len(toInsert)
		}
		if err := database.DB.Create(toInsert[i:end]).Error; err == nil {
			imported += end - i
		}
	}
	OK(c, gin.H{"imported": imported, "total": len(req.Items)})
}

func formatInt64(n int64) string {
	return strconv.FormatInt(n, 10)
}
