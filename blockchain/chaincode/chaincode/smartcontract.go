package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// ==================== 用户注册 ====================
func (s *SmartContract) RegisterUser(ctx contractapi.TransactionContextInterface, userID string, userType string, realInfoHash string) error {
	user := User{
		UserID:       userID,
		UserType:     userType,
		RealInfoHash: realInfoHash,
		CertList:     []*Certificate{},
	}
	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(userID, userAsBytes)
}

// ==================== 个人用户 — 证照上链 ====================
func (s *SmartContract) Uplink(ctx contractapi.TransactionContextInterface, userID string, traceability_code string, arg1 string, arg2 string, arg3 string, arg4 string, arg5 string, cert_image string) (string, error) {
	exists, err := ctx.GetStub().GetState(traceability_code)
	if err != nil {
		return "", fmt.Errorf("查询溯源码失败: %v", err)
	}
	if exists != nil {
		return "", fmt.Errorf("溯源码 %s 已存在", traceability_code)
	}

	cert := Certificate{
		TraceabilityCode: traceability_code,
		PersonalInput: PersonalInput{
			CertType:   arg1,
			CertNumber: arg2,
			Gender:     arg3,
			Phone:      arg4,
			HolderName: arg5,
			CertImage:  cert_image,
		},
	}

	txtime, _ := ctx.GetStub().GetTxTimestamp()
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	txTime := time.Unix(txtime.Seconds, 0).In(timeLocation).Format("2006-01-02 15:04:05")
	txid := ctx.GetStub().GetTxID()

	cert.PersonalInput.TxId = txid
	cert.PersonalInput.Timestamp = txTime

	certAsBytes, _ := json.Marshal(cert)
	if err := ctx.GetStub().PutState(traceability_code, certAsBytes); err != nil {
		return "", fmt.Errorf("写入区块链失败: %v", err)
	}

	if err := s.AddCertificate(ctx, userID, &cert); err != nil {
		return "", fmt.Errorf("添加到用户列表失败: %v", err)
	}

	// 记录上链历史（使用持证人姓名和证件编号）
	recordHistory(ctx.GetStub(),
		cert.TraceabilityCode,
		"证照上链",
		cert.PersonalInput.HolderName,
		"个人用户",
		cert.PersonalInput.CertNumber,
		"无",
		"待审核")

	return txid, nil
}

// ==================== 政务部门 — 审核证照 ====================
func (s *SmartContract) GovtAudit(ctx contractapi.TransactionContextInterface, userID string, traceability_code string, deptName string, deptCode string, auditorName string, auditResult string) (string, error) {
	certBytes, err := ctx.GetStub().GetState(traceability_code)
	if err != nil || certBytes == nil {
		return "", fmt.Errorf("证照 %s 不存在", traceability_code)
	}

	var cert Certificate
	json.Unmarshal(certBytes, &cert)

	txtime, _ := ctx.GetStub().GetTxTimestamp()
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	txTime := time.Unix(txtime.Seconds, 0).In(timeLocation).Format("2006-01-02 15:04:05")
	txid := ctx.GetStub().GetTxID()

	cert.GovtInput = GovtInput{
		DeptName:    deptName,
		DeptCode:    deptCode,
		AuditTime:   txTime,
		AuditorName: auditorName,
		AuditResult: auditResult,
		TxId:        txid,
		Timestamp:   txTime,
	}

	certAsBytes, _ := json.Marshal(cert)
	if err := ctx.GetStub().PutState(traceability_code, certAsBytes); err != nil {
		return "", fmt.Errorf("更新证照失败: %v", err)
	}

	// 记录审核历史
	holderName := cert.PersonalInput.HolderName
	certNumber := cert.PersonalInput.CertNumber
	newStatus := "审核通过"
	if auditResult == "驳回" {
		newStatus = "审核驳回"
	}
	recordHistory(ctx.GetStub(),
		cert.TraceabilityCode,
		"政务审核",
		holderName,
		"政务用户",
		certNumber,
		"待审核",
		newStatus)

	return txid, nil
}

// ==================== 企业组织 — 使用证照 ====================
func (s *SmartContract) EnterpriseUse(ctx contractapi.TransactionContextInterface, userID string, traceability_code string, companyName string, companyCode string, usePurpose string, operator string) (string, error) {
	certBytes, err := ctx.GetStub().GetState(traceability_code)
	if err != nil || certBytes == nil {
		return "", fmt.Errorf("证照 %s 不存在", traceability_code)
	}

	var cert Certificate
	json.Unmarshal(certBytes, &cert)

	txtime, _ := ctx.GetStub().GetTxTimestamp()
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	txTime := time.Unix(txtime.Seconds, 0).In(timeLocation).Format("2006-01-02 15:04:05")
	txid := ctx.GetStub().GetTxID()

	cert.EnterpriseInput = EnterpriseInput{
		CompanyName: companyName,
		CompanyCode: companyCode,
		UsePurpose:  usePurpose,
		UseTime:     txTime,
		Operator:    operator,
		TxId:        txid,
		Timestamp:   txTime,
	}

	certAsBytes, _ := json.Marshal(cert)
	if err := ctx.GetStub().PutState(traceability_code, certAsBytes); err != nil {
		return "", fmt.Errorf("更新证照失败: %v", err)
	}

	// 记录企业备案历史
	holderName := cert.PersonalInput.HolderName
	certNumber := cert.PersonalInput.CertNumber
	recordHistory(ctx.GetStub(),
		cert.TraceabilityCode,
		"企业备案",
		holderName,
		"企业用户",
		certNumber,
		"已审核",
		"已备案")

	return txid, nil
}

// ==================== 技术支撑实体 — 核验证照 ====================
func (s *SmartContract) TechVerify(ctx contractapi.TransactionContextInterface, userID string, traceability_code string, entityName string, serviceType string, securityLevel string, verifyResult string) (string, error) {
	certBytes, err := ctx.GetStub().GetState(traceability_code)
	if err != nil || certBytes == nil {
		return "", fmt.Errorf("证照 %s 不存在", traceability_code)
	}

	var cert Certificate
	json.Unmarshal(certBytes, &cert)

	txtime, _ := ctx.GetStub().GetTxTimestamp()
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	txTime := time.Unix(txtime.Seconds, 0).In(timeLocation).Format("2006-01-02 15:04:05")
	txid := ctx.GetStub().GetTxID()

	cert.TechSupportInput = TechSupportInput{
		EntityName:    entityName,
		ServiceType:   serviceType,
		SecurityLevel: securityLevel,
		VerifyResult:  verifyResult,
		ContactPhone:  "",
		TxId:          txid,
		Timestamp:     txTime,
	}

	certAsBytes, _ := json.Marshal(cert)
	if err := ctx.GetStub().PutState(traceability_code, certAsBytes); err != nil {
		return "", fmt.Errorf("更新证照失败: %v", err)
	}

	// 记录技术核验历史
	holderName := cert.PersonalInput.HolderName
	certNumber := cert.PersonalInput.CertNumber
	newStatus := "已核验"
	if verifyResult != "有效" {
		newStatus = "核验无效"
	}
	recordHistory(ctx.GetStub(),
		cert.TraceabilityCode,
		"技术核验",
		holderName,
		"技术支撑实体",
		certNumber,
		"已备案",
		newStatus)

	return txid, nil
}

// ==================== 添加证照到用户列表 ====================
func (s *SmartContract) AddCertificate(ctx contractapi.TransactionContextInterface, userID string, cert *Certificate) error {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil || userBytes == nil {
		return fmt.Errorf("用户 %s 不存在", userID)
	}
	var user User
	json.Unmarshal(userBytes, &user)
	user.CertList = append(user.CertList, cert)
	userAsBytes, _ := json.Marshal(user)
	return ctx.GetStub().PutState(userID, userAsBytes)
}

// ==================== 获取用户类型 ====================
func (s *SmartContract) GetUserType(ctx contractapi.TransactionContextInterface, userID string) (string, error) {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil || userBytes == nil {
		return "", fmt.Errorf("用户 %s 不存在", userID)
	}
	var user User
	json.Unmarshal(userBytes, &user)
	return user.UserType, nil
}

// ==================== 获取用户信息 ====================
func (s *SmartContract) GetUserInfo(ctx contractapi.TransactionContextInterface, userID string) (*User, error) {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil || userBytes == nil {
		return &User{}, fmt.Errorf("用户 %s 不存在", userID)
	}
	var user User
	json.Unmarshal(userBytes, &user)
	return &user, nil
}

// ==================== 查询证照信息 ====================
func (s *SmartContract) GetFruitInfo(ctx contractapi.TransactionContextInterface, traceability_code string) (*Certificate, error) {
	certBytes, err := ctx.GetStub().GetState(traceability_code)
	if err != nil || certBytes == nil {
		return &Certificate{}, fmt.Errorf("证照 %s 不存在", traceability_code)
	}
	var cert Certificate
	json.Unmarshal(certBytes, &cert)
	return &cert, nil
}

// ==================== 获取用户证照列表 ====================
func (s *SmartContract) GetFruitList(ctx contractapi.TransactionContextInterface, userID string) ([]*Certificate, error) {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil || userBytes == nil {
		return nil, fmt.Errorf("用户 %s 不存在", userID)
	}
	var user User
	json.Unmarshal(userBytes, &user)

	var freshCerts []*Certificate
	for _, c := range user.CertList {
		if c == nil || c.TraceabilityCode == "" {
			continue
		}
		certBytes, err := ctx.GetStub().GetState(c.TraceabilityCode)
		if err != nil || certBytes == nil {
			continue
		}
		var cert Certificate
		json.Unmarshal(certBytes, &cert)
		if cert.TraceabilityCode != "" {
			freshCerts = append(freshCerts, &cert)
		}
	}
	return freshCerts, nil
}

// ==================== 获取所有证照 ====================
func (s *SmartContract) GetAllFruitInfo(ctx contractapi.TransactionContextInterface) ([]Certificate, error) {
	iter, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var certs []Certificate
	for iter.HasNext() {
		res, err := iter.Next()
		if err != nil {
			return nil, err
		}
		var cert Certificate
		json.Unmarshal(res.Value, &cert)
		if cert.TraceabilityCode != "" {
			certs = append(certs, cert)
		}
	}
	return certs, nil
}

// ==================== 获取证照历史（Fabric自带历史） ====================
func (s *SmartContract) GetFruitHistory(ctx contractapi.TransactionContextInterface, traceability_code string) ([]HistoryQueryResult, error) {
	log.Printf("GetCertHistory: %v", traceability_code)
	iter, err := ctx.GetStub().GetHistoryForKey(traceability_code)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var records []HistoryQueryResult
	for iter.HasNext() {
		res, err := iter.Next()
		if err != nil {
			return nil, err
		}
		var cert Certificate
		if len(res.Value) > 0 {
			json.Unmarshal(res.Value, &cert)
		} else {
			cert = Certificate{TraceabilityCode: traceability_code}
		}

		timestamp, _ := ptypes.Timestamp(res.Timestamp)
		loc, _ := time.LoadLocation("Asia/Shanghai")
		formattedTime := timestamp.In(loc).Format("2006-01-02 15:04:05")

		records = append(records, HistoryQueryResult{
			TxId:      res.TxId,
			Timestamp: formattedTime,
			Record:    &cert,
			IsDelete:  res.IsDelete,
		})
	}
	return records, nil
}

// ==================== 仪表盘统计 ====================
func (s *SmartContract) GetStats(ctx contractapi.TransactionContextInterface) (*DashboardStats, error) {
	iter, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	stats := &DashboardStats{
		CertTypeDist: make(map[string]int),
		StageDist: map[string]int{
			"已上链": 0, "已审核": 0, "已备案": 0, "已核验": 0,
		},
		RecentActivity: []string{},
	}

	for iter.HasNext() {
		res, err := iter.Next()
		if err != nil {
			continue
		}
		var cert Certificate
		json.Unmarshal(res.Value, &cert)
		if cert.TraceabilityCode == "" {
			continue
		}

		stats.TotalCerts++

		certType := cert.PersonalInput.CertType
		if certType == "" {
			certType = "未分类"
		}
		stats.CertTypeDist[certType]++

		if cert.TechSupportInput.TxId != "" {
			stats.StageDist["已核验"]++
		} else if cert.EnterpriseInput.TxId != "" {
			stats.StageDist["已备案"]++
		} else if cert.GovtInput.TxId != "" {
			stats.StageDist["已审核"]++
		} else if cert.PersonalInput.TxId != "" {
			stats.StageDist["已上链"]++
		}
	}

	return stats, nil
}

// GetChainStats 返回仪表盘统计（证照维度）
func (s *SmartContract) GetChainStats(ctx contractapi.TransactionContextInterface) (string, error) {
	certs, err := s.GetAllFruitInfo(ctx)
	if err != nil {
		return "", err
	}

	stats := DashboardStats{
		TotalCerts:     len(certs),
		CertTypeDist:   make(map[string]int),
		StageDist:      map[string]int{"已上链": 0, "已审核": 0, "已备案": 0, "已核验": 0},
		RecentActivity: []string{},
	}

	for _, cert := range certs {
		ct := cert.PersonalInput.CertType
		if ct == "" {
			ct = "未分类"
		}
		stats.CertTypeDist[ct]++

		if cert.TechSupportInput.TxId != "" {
			stats.StageDist["已核验"]++
		} else if cert.EnterpriseInput.TxId != "" {
			stats.StageDist["已备案"]++
		} else if cert.GovtInput.TxId != "" {
			stats.StageDist["已审核"]++
		} else if cert.PersonalInput.TxId != "" {
			stats.StageDist["已上链"]++
		}
	}
	statsBytes, _ := json.Marshal(stats)
	return string(statsBytes), nil
}

// ==================== 内部历史记录函数 ====================
func recordHistory(stub shim.ChaincodeStubInterface, certId string, eventType string, operator string, role string, certNumber string, statusBefore string, statusAfter string) error {
	txid := stub.GetTxID()
	ts, _ := stub.GetTxTimestamp()
	timestampUnix := fmt.Sprintf("%d", ts.Seconds) // Unix时间戳（秒）

	history := map[string]interface{}{
		"eventType":    eventType,
		"operator":     operator,
		"role":         role,
		"certNumber":   certNumber,
		"statusBefore": statusBefore,
		"statusAfter":  statusAfter,
		"time":         timestampUnix,
		"txid":         txid,
	}
	b, _ := json.Marshal(history)
	key := "HISTORY_" + certId + "_" + txid
	return stub.PutState(key, b)
}

// ==================== 查询历史记录（供前端调用） ====================
func (s *SmartContract) QueryCertificateHistory(ctx contractapi.TransactionContextInterface, certId string) (string, error) {
	iter, err := ctx.GetStub().GetStateByRange("HISTORY_"+certId+"_", "HISTORY_"+certId+"_~")
	if err != nil {
		return "[]", err
	}
	defer iter.Close()

	var histories []map[string]interface{}
	for iter.HasNext() {
		v, err := iter.Next()
		if err != nil {
			continue
		}
		var record map[string]interface{}
		if err := json.Unmarshal(v.Value, &record); err != nil {
			continue
		}
		histories = append(histories, record)
	}
	// 按时间戳排序（升序）
	sort.Slice(histories, func(i, j int) bool {
		ti, _ := strconv.ParseInt(histories[i]["time"].(string), 10, 64)
		tj, _ := strconv.ParseInt(histories[j]["time"].(string), 10, 64)
		return ti < tj
	})
	resultBytes, _ := json.Marshal(histories)
	return string(resultBytes), nil
}
// ==================== 供应商关联 ====================
// 添加供应商关联（企业针对某证照记录供应商信息）
func (s *SmartContract) AddSupplierLink(ctx contractapi.TransactionContextInterface, traceabilityCode string, enterpriseId string, supplierName string, contactPerson string, contactPhone string, startDate string, endDate string, notes string) (string, error) {
    exists, err := ctx.GetStub().GetState(traceabilityCode)
    if err != nil || exists == nil {
        return "", fmt.Errorf("证照 %s 不存在", traceabilityCode)
    }

    txid := ctx.GetStub().GetTxID()
    ts, _ := ctx.GetStub().GetTxTimestamp()
    timestamp := fmt.Sprintf("%d", ts.Seconds)

    link := map[string]interface{}{
        "traceability_code": traceabilityCode,
        "enterprise_id":     enterpriseId,
        "supplier_name":     supplierName,
        "contact_person":    contactPerson,
        "contact_phone":     contactPhone,
        "cooperation_start": startDate,
        "cooperation_end":   endDate,
        "notes":             notes,
        "timestamp":         timestamp,
        "txid":              txid,
    }
    key := "SUPPLIER_" + traceabilityCode + "_" + enterpriseId + "_" + txid
    linkBytes, _ := json.Marshal(link)
    err = ctx.GetStub().PutState(key, linkBytes)
    if err != nil {
        return "", fmt.Errorf("保存供应商关联失败: %v", err)
    }
    return txid, nil
}

// 查询证照的所有供应商关联记录
func (s *SmartContract) GetSupplierLinks(ctx contractapi.TransactionContextInterface, traceabilityCode string) (string, error) {
    iter, err := ctx.GetStub().GetStateByRange("SUPPLIER_"+traceabilityCode+"_", "SUPPLIER_"+traceabilityCode+"_~")
    if err != nil {
        return "[]", err
    }
    defer iter.Close()
    var links []map[string]interface{}
    for iter.HasNext() {
        res, err := iter.Next()
        if err != nil {
            continue
        }
        var link map[string]interface{}
        if err := json.Unmarshal(res.Value, &link); err != nil {
            continue
        }
        links = append(links, link)
    }
    result, _ := json.Marshal(links)
    return string(result), nil
}

// ==================== 合规事件 ====================
// 添加合规事件（企业针对证照记录合规事件，如年检、整改等）
func (s *SmartContract) AddComplianceEvent(ctx contractapi.TransactionContextInterface, traceabilityCode string, enterpriseId string, eventType string, eventTime string, description string, attachmentCid string, attachmentFingerprint string) (string, error) {
    exists, err := ctx.GetStub().GetState(traceabilityCode)
    if err != nil || exists == nil {
        return "", fmt.Errorf("证照 %s 不存在", traceabilityCode)
    }

    txid := ctx.GetStub().GetTxID()
    ts, _ := ctx.GetStub().GetTxTimestamp()
    timestamp := fmt.Sprintf("%d", ts.Seconds)

    event := map[string]interface{}{
        "traceability_code":   traceabilityCode,
        "enterprise_id":       enterpriseId,
        "event_type":          eventType,
        "event_time":          eventTime,
        "description":         description,
        "attachment_cid":      attachmentCid,
        "attachment_fingerprint": attachmentFingerprint,
        "record_timestamp":    timestamp,
        "txid":                txid,
    }
    key := "COMPLIANCE_" + traceabilityCode + "_" + enterpriseId + "_" + txid
    eventBytes, _ := json.Marshal(event)
    err = ctx.GetStub().PutState(key, eventBytes)
    if err != nil {
        return "", fmt.Errorf("保存合规事件失败: %v", err)
    }
    return txid, nil
}

// 查询证照的所有合规事件
func (s *SmartContract) GetComplianceEvents(ctx contractapi.TransactionContextInterface, traceabilityCode string) (string, error) {
    iter, err := ctx.GetStub().GetStateByRange("COMPLIANCE_"+traceabilityCode+"_", "COMPLIANCE_"+traceabilityCode+"_~")
    if err != nil {
        return "[]", err
    }
    defer iter.Close()
    var events []map[string]interface{}
    for iter.HasNext() {
        res, err := iter.Next()
        if err != nil {
            continue
        }
        var event map[string]interface{}
        if err := json.Unmarshal(res.Value, &event); err != nil {
            continue
        }
        events = append(events, event)
    }
    result, _ := json.Marshal(events)
    return string(result), nil
}
// GetDataByTxId 根据交易哈希查询存储的数据（用于验证）
func (s *SmartContract) GetDataByTxId(ctx contractapi.TransactionContextInterface, txId string) (string, error) {
    // 遍历所有以 "SUPPLIER_" 开头的 key（可根据需要扩展其他前缀）
    iter, err := ctx.GetStub().GetStateByRange("SUPPLIER_", "SUPPLIER_~")
    if err != nil {
        return "", err
    }
    defer iter.Close()
    for iter.HasNext() {
        res, err := iter.Next()
        if err != nil {
            continue
        }
        var data map[string]interface{}
        if err := json.Unmarshal(res.Value, &data); err != nil {
            continue
        }
        if data["txid"] == txId {
            // 找到匹配的记录，返回原始数据
            result, _ := json.Marshal(data)
            return string(result), nil
        }
    }
    return "", fmt.Errorf("未找到 txid %s 对应的记录", txId)
}

