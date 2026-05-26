package controller

import (
	"backend/pkg"
	"encoding/json"
	"fmt"   // 新增 fmt 包

	"github.com/gin-gonic/gin"
)

// GetEnterpriseStats 企业组织端统计：返回当前企业用户备案过的证照统计
func GetEnterpriseStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未登录"})
		return
	}
	userType, _ := c.Get("userType")
	if userType != "企业组织" {
		c.JSON(403, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	// 获取所有证照
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

	// 筛选出当前企业用户备案过的证照（driver_input.dr_transport == userID）
	var certs []map[string]interface{}
	for _, cert := range allCerts {
		driver, ok := cert["driver_input"].(map[string]interface{})
		if !ok {
			continue
		}
		operator, _ := driver["dr_transport"].(string)
		if operator == userID.(string) {
			certs = append(certs, cert)
		}
	}

	// 统计阶段分布和种类分布
	stageDist := map[string]int{"已上链": 0, "已审核": 0, "已备案": 0, "已核验": 0}
	certTypeDist := make(map[string]int)

	for _, cert := range certs {
		farmer, _ := cert["farmer_input"].(map[string]interface{})
		govt, _ := cert["factory_input"].(map[string]interface{})
		ent, _ := cert["driver_input"].(map[string]interface{})
		tech, _ := cert["shop_input"].(map[string]interface{})

		stage := "已上链"
		if tech != nil && tech["sh_txid"] != "" {
			stage = "已核验"
		} else if ent != nil && ent["dr_txid"] != "" {
			stage = "已备案"
		} else if govt != nil && govt["fac_txid"] != "" {
			stage = "已审核"
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
			"totalCerts":   len(certs),
			"stageDist":    stageDist,
			"certTypeDist": certTypeDist,
		},
	})
}

// ExportAllEnterpriseHistory 导出当前企业所有备案证照的历史记录为CSV
func ExportAllEnterpriseHistory(c *gin.Context) {
	userID, _ := c.Get("userID")
	userType, _ := c.Get("userType")
	if userType != "企业组织" {
		c.JSON(403, gin.H{"code": 403, "message": "权限不足"})
		return
	}

	// 获取所有证照
	allRes, err := pkg.ChaincodeQuery("GetAllFruitInfo", "")
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "查询证照失败"})
		return
	}
	var allCerts []map[string]interface{}
	json.Unmarshal([]byte(allRes), &allCerts)

	// 筛选当前企业备案的证照
	var certs []map[string]interface{}
	for _, cert := range allCerts {
		driver, ok := cert["driver_input"].(map[string]interface{})
		if !ok {
			continue
		}
		operator, _ := driver["dr_transport"].(string)
		if operator == userID.(string) {
			certs = append(certs, cert)
		}
	}

	// 收集所有历史记录
	var allHistory []map[string]interface{}
	for _, cert := range certs {
		certId, _ := cert["traceability_code"].(string)
		historyRes, _ := pkg.ChaincodeQuery("QueryCertificateHistory", certId)
		var history []map[string]interface{}
		json.Unmarshal([]byte(historyRes), &history)
		for _, h := range history {
			h["certId"] = certId
			allHistory = append(allHistory, h)
		}
	}

	// 生成CSV
	csvData := "证照溯源码,时间,事件类型,操作人,角色,证件编号,状态变更前,状态变更后,交易哈希,CID\n"
	for _, h := range allHistory {
		csvData += fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
			h["certId"], h["time"], h["eventType"], h["operator"], h["role"],
			h["certNumber"], h["statusBefore"], h["statusAfter"], h["txid"], h["cid"])
	}
	// 添加 UTF-8 BOM 避免中文乱码
	bom := []byte{0xEF, 0xBB, 0xBF}
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment;filename=enterprise_audit_all.csv")
	c.Data(200, "text/csv; charset=utf-8", append(bom, []byte(csvData)...))
}
// EnterpriseVerifyEvidence 企业端证据核验（验证 txid 和 cid 是否匹配）
func EnterpriseVerifyEvidence(c *gin.Context) {
	// 只检查登录，不限制角色
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未登录"})
		return
	}
	_ = userID // 可用于日志

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

	// 查询证照历史
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
