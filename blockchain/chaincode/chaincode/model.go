package chaincode

// ==================== 用户结构体 ====================
type User struct {
	UserID       string         `json:"userID"`
	UserType     string         `json:"userType"`
	RealInfoHash string         `json:"realInfoHash"`
	CertList     []*Certificate `json:"fruitList"` // JSON标签保持兼容
}

// ==================== 电子证照主结构体 ====================
type Certificate struct {
	TraceabilityCode string           `json:"traceability_code"`
	PersonalInput    PersonalInput    `json:"farmer_input"`    // 个人用户
	GovtInput        GovtInput        `json:"factory_input"`   // 政务部门
	EnterpriseInput  EnterpriseInput  `json:"driver_input"`    // 企业组织
	TechSupportInput TechSupportInput `json:"shop_input"`      // 技术支撑实体
}

// ==================== 历史查询结果 ====================
type HistoryQueryResult struct {
	Record    *Certificate `json:"record"`
	TxId      string       `json:"txId"`
	Timestamp string       `json:"timestamp"`
	IsDelete  bool         `json:"isDelete"`
}

// ==================== 个人用户录入 ====================
type PersonalInput struct {
	CertType   string `json:"fa_fruitName"`   // 证照类型
	CertNumber string `json:"fa_origin"`      // 证照编号/身份证号
	Gender     string `json:"fa_plantTime"`   // 性别
	Phone      string `json:"fa_pickingTime"` // 联系电话
	HolderName string `json:"fa_farmerName"`  // 持证人姓名
	TxId       string `json:"fa_txid"`
	Timestamp  string `json:"fa_timestamp"`
	CertImage  string `json:"fa_certImage"` // 证照图片Base64
}

// ==================== 政务部门审核 ====================
type GovtInput struct {
	DeptName    string `json:"fac_productName"`    // 部门名称
	DeptCode    string `json:"fac_productionbatch"` // 部门代码
	AuditTime   string `json:"fac_productionTime"`  // 审核时间
	AuditorName string `json:"fac_factoryName"`     // 审核人
	AuditResult string `json:"fac_contactNumber"`   // 审核结果
	TxId        string `json:"fac_txid"`
	Timestamp   string `json:"fac_timestamp"`
}

// ==================== 企业组织使用记录 ====================
type EnterpriseInput struct {
	CompanyName string `json:"dr_name"`      // 企业名称
	CompanyCode string `json:"dr_age"`       // 统一社会信用代码
	UsePurpose  string `json:"dr_phone"`     // 使用目的
	UseTime     string `json:"dr_carNumber"` // 使用时间
	Operator    string `json:"dr_transport"` // 经办人
	TxId        string `json:"dr_txid"`
	Timestamp   string `json:"dr_timestamp"`
}

// ==================== 技术支撑实体核验 ====================
type TechSupportInput struct {
	EntityName    string `json:"sh_storeTime"`   // 实体名称
	ServiceType   string `json:"sh_sellTime"`    // 服务类型
	SecurityLevel string `json:"sh_shopName"`    // 安全认证等级
	VerifyResult  string `json:"sh_shopAddress"` // 核验结果
	ContactPhone  string `json:"sh_shopPhone"`   // 联系方式
	TxId          string `json:"sh_txid"`
	Timestamp     string `json:"sh_timestamp"`
}

// ==================== 仪表盘统计数据 ====================
type DashboardStats struct {
	TotalCerts      int            `json:"totalCerts"`
	CertTypeDist    map[string]int `json:"certTypeDist"`    // 证照种类分布
	StageDist       map[string]int `json:"stageDist"`       // 阶段分布
	RecentActivity  []string       `json:"recentActivity"`
}
