package controller

import (
	"backend/pkg"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 1. 证照上链
func Uplink(c *gin.Context) {
	farmer_traceability_code := pkg.GenerateID()[1:]
	args := buildArgs(c, farmer_traceability_code)
	if args == nil {
		return
	}
	res, err := pkg.ChaincodeInvoke("Uplink", args)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "uplink failed " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code":              200,
		"message":           "uplink success",
		"txid":              res,
		"traceability_code": args[1],
	})
}

// 2. 获取单个证照信息
func GetFruitInfo(c *gin.Context) {
	traceCode := c.Query("traceabilityCode")
	if traceCode == "" {
		traceCode = c.PostForm("traceability_code")
	}

	res, err := pkg.ChaincodeQuery("GetFruitInfo", traceCode)
	if err == nil {
		_ = pkg.InsertTraceLog(traceCode, "")
	}
	if err != nil {
		c.JSON(200, gin.H{
			"message": "查询失败：" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}

// 3. 获取个人证照列表
func GetFruitList(c *gin.Context) {
	userID, _ := c.Get("userID")
	res, err := pkg.ChaincodeQuery("GetFruitList", userID.(string))
	if err != nil {
		c.JSON(200, gin.H{"code": 500, "message": "查询失败：" + err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "message": "query success", "data": res})
}

// 4. 获取所有证照信息
func GetAllFruitInfo(c *gin.Context) {
	res, err := pkg.ChaincodeQuery("GetAllFruitInfo", "")
	fmt.Print("res", res)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed " + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}

// 5. 获取上链历史（Fabric原生历史）
func GetFruitHistory(c *gin.Context) {
	traceCode := c.Query("traceabilityCode")
	if traceCode == "" {
		traceCode = c.PostForm("traceability_code")
	}

	res, err := pkg.ChaincodeQuery("GetFruitHistory", traceCode)
	if err == nil {
		_ = pkg.InsertTraceLog(traceCode, "")
	}

	if err != nil {
		c.JSON(200, gin.H{
			"message": "query failed " + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "query success",
		"data":    res,
	})
}

func buildArgs(c *gin.Context, farmer_traceability_code string) []string {
	var args []string
	userID, _ := c.Get("userID")
	userType, _ := pkg.ChaincodeQuery("GetUserType", userID.(string))
	args = append(args, userID.(string))

	if userType == "个人用户" {
		args = append(args, farmer_traceability_code)
	} else {
		args = append(args, c.PostForm("traceability_code"))
	}

	args = append(args, c.PostForm("arg1"))
	args = append(args, c.PostForm("arg2"))
	args = append(args, c.PostForm("arg3"))
	args = append(args, c.PostForm("arg4"))
	args = append(args, c.PostForm("arg5"))

	var certImageURL string
	file, err := c.FormFile("cert_image")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			c.JSON(200, gin.H{"code": 500, "message": "打开图片失败"})
			return nil
		}
		defer src.Close()

		cid, err := pkg.UploadToIPFS(src, file)
		if err != nil {
			c.JSON(200, gin.H{"code": 500, "message": "IPFS上传失败: " + err.Error()})
			return nil
		}
		certImageURL = "ipfs://" + cid
	} else {
		certImageURL = ""
	}

	args = append(args, certImageURL)
	fmt.Println("最终传给链码的 args：", args)
	return args
}

// 6. 统计数据
func GetStats(c *gin.Context) {
	res, err := pkg.ChaincodeQuery("GetStats", "")
	if err != nil {
		c.JSON(200, gin.H{"code": 500, "message": "查询统计失败：" + err.Error()})
		return
	}

	height, txCount, bcErr := pkg.GetBlockchainInfo()
	blockHeight := "--"
	txcountStr := "--"
	if bcErr == nil {
		blockHeight = fmt.Sprintf("%d", height)
		txcountStr = fmt.Sprintf("%d", txCount)
	}
	traceCount, _ := pkg.GetTraceCount()
	c.JSON(200, gin.H{
		"code":        200,
		"message":     "query success",
		"data":        res,
		"blockHeight": blockHeight,
		"txCount":     txcountStr,
		"nodeCount":   4,
		"traceCount":  fmt.Sprintf("%d", traceCount),
	})
}

// 7. 仪表盘数据
func GetDashboard(c *gin.Context) {
	res, err := pkg.ChaincodeQuery("GetChainStats", "")
	if err != nil {
		c.JSON(200, gin.H{"code": 500, "message": "查询失败: " + err.Error()})
		return
	}

	var fullData map[string]interface{}
	json.Unmarshal([]byte(res), &fullData)

	c.JSON(200, gin.H{
		"code":       200,
		"fullData":   res,
		"nodeCount":  4,
		"traceCount": 0,
	})
}

// 8. 政务审核 —— 状态流转：待审核 → 通过/驳回
func GovtAudit(c *gin.Context) {
	userID, _ := c.Get("userID")
	traceCode := c.PostForm("traceability_code")
	auditRes := c.PostForm("auditResult")

	txid, err := pkg.ChaincodeInvoke("GovtAudit", []string{
		userID.(string),
		traceCode,
		c.PostForm("deptName"),
		c.PostForm("deptCode"),
		c.PostForm("auditorName"),
		auditRes,
	})
	if err != nil {
		c.JSON(200, gin.H{"code": 500, "message": "审核失败：" + err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "message": "审核成功", "txid": txid})
}

// 9. 企业备案 —— 状态流转：已审核 → 已备案
func EnterpriseUse(c *gin.Context) {
	userID, _ := c.Get("userID")
	traceCode := c.PostForm("traceability_code")

	txid, err := pkg.ChaincodeInvoke("EnterpriseUse", []string{
		userID.(string),
		traceCode,
		c.PostForm("companyName"),
		c.PostForm("companyCode"),
		c.PostForm("usePurpose"),
		c.PostForm("operator"),
	})
	if err != nil {
		c.JSON(200, gin.H{"code": 500, "message": "备案失败：" + err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "message": "使用备案成功", "txid": txid})
}

// 10. 技术核验 —— 状态流转：已备案 → 已核验/不通过
func TechVerify(c *gin.Context) {
	userID, _ := c.Get("userID")
	traceCode := c.PostForm("traceability_code")
	verifyRes := c.PostForm("verifyResult")

	txid, err := pkg.ChaincodeInvoke("TechVerify", []string{
		userID.(string),
		traceCode,
		c.PostForm("entityName"),
		c.PostForm("serviceType"),
		c.PostForm("securityLevel"),
		verifyRes,
	})
	if err != nil {
		c.JSON(200, gin.H{"code": 500, "message": "核验失败：" + err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "message": "核验成功", "txid": txid})
}

// 11. 获取完整溯源历史（前端时间轴使用）
func GetHistory(c *gin.Context) {
	certId := c.Param("certId")
	if certId == "" {
		certId = c.Query("certId")
	}
	res, err := pkg.ChaincodeQuery("QueryCertificateHistory", certId)
	if err != nil {
		c.JSON(200, gin.H{"code": 200, "data": []interface{}{}})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": json.RawMessage(res),
	})
}

// GetUserStats 获取当前登录用户自己上传的证照统计（用于可视化分析页面）
func GetUserStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"code": 401, "message": "未登录"})
		return
	}

	res, err := pkg.ChaincodeQuery("GetFruitList", userID.(string))
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "查询证照列表失败"})
		return
	}

	var certs []map[string]interface{}
	if err := json.Unmarshal([]byte(res), &certs); err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "解析数据失败"})
		return
	}

	stageDist := map[string]int{"已上链": 0, "已审核": 0, "已备案": 0, "已核验": 0}
	certTypeDist := make(map[string]int)

	for _, cert := range certs {
		farmer, _ := cert["farmer_input"].(map[string]interface{})
		govt, _ := cert["factory_input"].(map[string]interface{})
		ent, _ := cert["driver_input"].(map[string]interface{})
		tech, _ := cert["shop_input"].(map[string]interface{})

		if tech != nil && tech["sh_txid"] != "" {
			stageDist["已核验"]++
		} else if ent != nil && ent["dr_txid"] != "" {
			stageDist["已备案"]++
		} else if govt != nil && govt["fac_txid"] != "" {
			stageDist["已审核"]++
		} else if farmer != nil && farmer["fa_txid"] != "" {
			stageDist["已上链"]++
		}

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

// ExportHistory 导出证照历史为 CSV 文件
func ExportHistory(c *gin.Context) {
	certId := c.Param("certId")
	if certId == "" {
		certId = c.Query("certId")
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

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Write([]string{"时间", "事件类型", "操作人", "角色", "证件编号", "状态变更前", "状态变更后", "交易哈希", "CID"})
	for _, h := range history {
		writer.Write([]string{
			fmt.Sprintf("%v", h["time"]),
			fmt.Sprintf("%v", h["eventType"]),
			fmt.Sprintf("%v", h["operator"]),
			fmt.Sprintf("%v", h["role"]),
			fmt.Sprintf("%v", h["certNumber"]),
			fmt.Sprintf("%v", h["statusBefore"]),
			fmt.Sprintf("%v", h["statusAfter"]),
			fmt.Sprintf("%v", h["txid"]),
			fmt.Sprintf("%v", h["cid"]),
		})
	}
	writer.Flush()
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=history_%s.csv", certId))
	c.String(200, buf.String())
}
// UploadToIPFS 上传文件到IPFS，返回CID
func UploadToIPFS(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "未提供文件"})
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "打开文件失败"})
		return
	}
	defer src.Close()
	cid, err := pkg.UploadToIPFS(src, file)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "上传IPFS失败: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": 200, "cid": cid})
}
