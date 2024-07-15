package chaincode

/*
定义用户结构体
用户ID
用户类型
实名认证信息哈希,包括用户注册的姓名、身份证号、手机号、注册平台同意协议签名的哈希
电子证照列表
*/
type User struct {
	UserID       string   `json:"userID"`
	UserType     string   `json:"userType"`
	RealInfoHash string   `json:"realInfoHash"`
	FruitList    []*Fruit `json:"fruitList"`
}

/*
定义电子证照结构体
溯源码
个人用户输入
政务部门输入
企业组织输入
技术支撑实体输入
*/
type Fruit struct {
	Traceability_code string        `json:"traceability_code"`
	Farmer_input      Farmer_input  `json:"farmer_input"`
	Factory_input     Factory_input `json:"factory_input"`
	Driver_input      Driver_input  `json:"driver_input"`
	Shop_input        Shop_input    `json:"shop_input"`
}

// HistoryQueryResult structure used for handling result of history query
type HistoryQueryResult struct {
	Record    *Fruit `json:"record"`
	TxId      string `json:"txId"`
	Timestamp string `json:"timestamp"`
	IsDelete  bool   `json:"isDelete"`
}

/*
个人用户
电子证照的溯源码，一物一码，主打高端市场（自动生成）
电子证照类型、身份证号、性别、联系电话、个人用户名称
*/
type Farmer_input struct {
	Fa_fruitName   string `json:"fa_fruitName"`
	Fa_origin      string `json:"fa_origin"`
	Fa_plantTime   string `json:"fa_plantTime"`
	Fa_pickingTime string `json:"fa_pickingTime"`
	Fa_farmerName  string `json:"fa_farmerName"`
	Fa_Txid        string `json:"fa_txid"`
	Fa_Timestamp   string `json:"fa_timestamp"`
}

/*
政务部门
部门名称、部门代码、地址、负责人、联系电话
*/
type Factory_input struct {
	Fac_productName     string `json:"fac_productName"`
	Fac_productionbatch string `json:"fac_productionbatch"`
	Fac_productionTime  string `json:"fac_productionTime"`
	Fac_factoryName     string `json:"fac_factoryName"`
	Fac_contactNumber   string `json:"fac_contactNumber"`
	Fac_Txid            string `json:"fac_txid"`
	Fac_Timestamp       string `json:"fac_timestamp"`
}

/*
企业组织
企业名称、企业代码、成立日期、地址、联系方式
*/
type Driver_input struct {
	Dr_name      string `json:"dr_name"`
	Dr_age       string `json:"dr_age"`
	Dr_phone     string `json:"dr_phone"`
	Dr_carNumber string `json:"dr_carNumber"`
	Dr_transport string `json:"dr_transport"`
	Dr_Txid      string `json:"dr_txid"`
	Dr_Timestamp string `json:"dr_timestamp"`
}

/*
技术支撑实体
实体名称、服务类型、安全认证等级、地址、联系方式
*/
type Shop_input struct {
	Sh_storeTime   string `json:"sh_storeTime"`
	Sh_sellTime    string `json:"sh_sellTime"`
	Sh_shopName    string `json:"sh_shopName"`
	Sh_shopAddress string `json:"sh_shopAddress"`
	Sh_shopPhone   string `json:"sh_shopPhone"`
	Sh_Txid        string `json:"sh_txid"`
	Sh_Timestamp   string `json:"sh_timestamp"`
}
