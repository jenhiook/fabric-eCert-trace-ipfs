package controller

import (
	"backend/pkg"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GovGetHistory 政府端证照历史（带分页、时间区间）
func GovGetHistory(c *gin.Context) {
	certId := c.Param("certId")
	if certId == "" {
		certId = c.Query("certId")
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	res, err := pkg.ChaincodeQuery("QueryCertificateHistory", certId)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "查询历史失败"})
		return
	}
	var history []map[string]interface{}
	if err := json.Unmarshal([]byte(res), &history); err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "解析历史数据失败"})
		return
	}

	// 时间筛选
	filtered := make([]map[string]interface{}, 0)
	for _, h := range history {
		t, ok := h["time"].(string)
		if !ok {
			continue
		}
		ts, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			continue
		}
		if startTime != "" {
			start, _ := strconv.ParseInt(startTime, 10, 64)
			if ts < start {
				continue
			}
		}
		if endTime != "" {
			end, _ := strconv.ParseInt(endTime, 10, 64)
			if ts > end {
				continue
			}
		}
		filtered = append(filtered, h)
	}

	total := len(filtered)
	start := (page - 1) * limit
	end := start + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	paged := filtered[start:end]

	c.JSON(200, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"history": paged,
			"total":   total,
			"page":    page,
			"limit":   limit,
		},
	})
}

// GovGetAuditLog 政府端审计日志（复用历史）
func GovGetAuditLog(c *gin.Context) {
	certId := c.Param("certId")
	res, _ := pkg.ChaincodeQuery("QueryCertificateHistory", certId)
	c.JSON(200, gin.H{
		"code": 200,
		"data": json.RawMessage(res),
	})
}

// GovExportEvidencePack 导出证据包（JSON）
func GovExportEvidencePack(c *gin.Context) {
	certId := c.Param("certId")
	res, err := pkg.ChaincodeQuery("QueryCertificateHistory", certId)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=evidence_%s.json", certId))
	c.String(200, res)
}

// GovAuditReport 导出审计报表（CSV，带 UTF-8 BOM）
func GovAuditReport(c *gin.Context) {
	certId := c.Query("certId")
	res, err := pkg.ChaincodeQuery("QueryCertificateHistory", certId)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	var history []map[string]interface{}
	json.Unmarshal([]byte(res), &history)

	csvData := "时间,事件类型,操作人,角色,证件编号,状态变更前,状态变更后,交易哈希,CID\n"
	for _, h := range history {
		csvData += fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
			h["time"], h["eventType"], h["operator"], h["role"],
			h["certNumber"], h["statusBefore"], h["statusAfter"],
			h["txid"], h["cid"])
	}
	// 添加 UTF-8 BOM 防止 Excel 打开中文乱码
	bom := []byte{0xEF, 0xBB, 0xBF}
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=audit_%s.csv", certId))
	c.Data(200, "text/csv; charset=utf-8", append(bom, []byte(csvData)...))
}

// GovVerifyEvidence 对外核验接口（验证CID+txid）
func GovVerifyEvidence(c *gin.Context) {
	var req struct {
		CertId     string `json:"certId"`
		TxId       string `json:"txId"`
		Cid        string `json:"cid"`
		Fingerprint string `json:"fingerprint"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	res, err := pkg.ChaincodeQuery("QueryCertificateHistory", req.CertId)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	var history []map[string]interface{}
	json.Unmarshal([]byte(res), &history)
	valid := false
	for _, h := range history {
		if h["txid"] == req.TxId && h["cid"] == req.Cid {
			valid = true
			break
		}
	}
	if valid {
		c.JSON(200, gin.H{"code": 200, "valid": true, "message": "核验通过"})
	} else {
		c.JSON(200, gin.H{"code": 200, "valid": false, "message": "核验失败，记录不匹配"})
	}
}

// GovGetStats 政府端统计：全局证照统计（所有证照）
func GovGetStats(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未登录"})
		return
	}
	userType, _ := c.Get("userType")
	if userType != "政务部门" {
		c.JSON(403, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	res, err := pkg.ChaincodeQuery("GetAllFruitInfo", "")
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "查询证照失败"})
		return
	}
	var allCerts []map[string]interface{}
	if err := json.Unmarshal([]byte(res), &allCerts); err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "解析证照数据失败"})
		return
	}

	stageDist := map[string]int{"已上链": 0, "已审核": 0, "已备案": 0, "已核验": 0}
	certTypeDist := make(map[string]int)

	for _, cert := range allCerts {
		farmer, _ := cert["farmer_input"].(map[string]interface{})
		govt, _ := cert["factory_input"].(map[string]interface{})
		ent, _ := cert["driver_input"].(map[string]interface{})
		tech, _ := cert["shop_input"].(map[string]interface{})

		stage := "已上链"
		if tech != nil {
			if txid, ok := tech["sh_txid"]; ok && txid != nil && txid != "" {
				stage = "已核验"
			}
		}
		if stage == "已上链" && ent != nil {
			if txid, ok := ent["dr_txid"]; ok && txid != nil && txid != "" {
				stage = "已备案"
			}
		}
		if stage == "已上链" && govt != nil {
			if txid, ok := govt["fac_txid"]; ok && txid != nil && txid != "" {
				stage = "已审核"
			}
		}
		stageDist[stage]++

		certType, _ := farmer["fa_fruitName"].(string)
		if certType == "" {
			certType = "未分类"
		}
		certTypeDist[certType]++
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"totalCerts":   len(allCerts),
			"stageDist":    stageDist,
			"certTypeDist": certTypeDist,
		},
	})
}
