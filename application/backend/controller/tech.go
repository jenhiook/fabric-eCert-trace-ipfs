package controller

import (
	"backend/pkg"
	"encoding/json"
	"log"   // 新增导入
	"github.com/gin-gonic/gin"
)

// GetTechStats 技术实体端统计：返回当前技术用户核验过的证照统计
func GetTechStats(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未登录"})
		return
	}
	userType, _ := c.Get("userType")
	if userType != "技术支撑实体" {
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

	// 筛选出当前技术用户核验过的证照（存在 shop_input 且 sh_txid 不为空）
	var certs []map[string]interface{}
	for _, cert := range allCerts {
		tech, ok := cert["shop_input"].(map[string]interface{})
		if ok && tech != nil {
			if txid, ok := tech["sh_txid"]; ok && txid != nil && txid != "" {
				certs = append(certs, cert)
			}
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

// TechVerifyEvidence 技术实体端证据核验接口
func TechVerifyEvidence(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未登录"})
		return
	}
	var req struct {
		CertId string `json:"certId"`
		TxId   string `json:"txId"`
		Cid    string `json:"cid"`
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
	if err := json.Unmarshal([]byte(res), &history); err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "解析历史数据失败"})
		return
	}

	// 调试日志
	for i, h := range history {
		if i < 3 {
			log.Printf("历史记录 %d: txid=%v, cid=%v", i, h["txid"], h["cid"])
		}
	}

	valid := false
	for _, h := range history {
		chainTxid, _ := h["txid"].(string)
		chainCid, _ := h["cid"].(string)

		if chainTxid != req.TxId {
			continue
		}

		// 去除 ipfs:// 前缀后比较
		inputCid := req.Cid
		if len(inputCid) > 7 && inputCid[:7] == "ipfs://" {
			inputCid = inputCid[7:]
		}
		if len(chainCid) > 7 && chainCid[:7] == "ipfs://" {
			chainCid = chainCid[7:]
		}

		if inputCid == chainCid {
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
