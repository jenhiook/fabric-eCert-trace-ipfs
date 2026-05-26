package router

import (
	con "backend/controller"
	"backend/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 公开路由
	r.POST("/register", con.Register)
	r.POST("/login", con.Login)
	r.POST("/getFruitInfo", con.GetFruitInfo)

	// 需要鉴权的路由组（普通用户）
	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.POST("/logout", con.Logout)
		auth.POST("/getInfo", con.GetInfo)
		auth.POST("/uplink", con.Uplink)
		auth.POST("/getFruitList", con.GetFruitList)
		auth.POST("/getAllFruitInfo", con.GetAllFruitInfo)
		auth.POST("/getFruitHistory", con.GetFruitHistory)
		auth.POST("/getStats", con.GetStats)
		auth.POST("/getDashboard", con.GetDashboard)
		auth.POST("/govtAudit", con.GovtAudit)
		auth.POST("/enterpriseUse", con.EnterpriseUse)
		auth.POST("/techVerify", con.TechVerify)
		auth.GET("/cert/:certId/history", con.GetHistory)
		auth.GET("/cert/:certId/export", con.ExportHistory)
		auth.POST("/userStats", con.GetUserStats)
		// 企业端扩展功能（链上存证）
		auth.POST("/enterprise/supplier/add", con.AddSupplierLink)
		auth.GET("/enterprise/supplier/list", con.GetSupplierLinks)
		auth.POST("/enterprise/event/add", con.AddComplianceEvent)
		auth.GET("/enterprise/event/list", con.GetComplianceEvents)
		
		// 企业统计接口（需要登录）
		auth.POST("/enterpriseStats", con.GetEnterpriseStats)
		auth.GET("/enterprise/exportAll", con.ExportAllEnterpriseHistory)
		auth.POST("/enterprise/verify", con.EnterpriseVerifyEvidence)
		auth.POST("/upload/ipfs", con.UploadToIPFS)
		auth.POST("/enterprise/verify/chaindata", con.VerifyChainData)
		
		auth.POST("/techStats", con.GetTechStats)
		auth.POST("/tech/verify", con.TechVerifyEvidence)
	}

	// ==================== 政府端路由组 ====================
	gov := r.Group("/gov")
	gov.Use(middleware.JWTAuthMiddleware(), middleware.RoleCheck("政务部门"))
	{
		gov.GET("/cert/:certId/history", con.GovGetHistory)
		gov.GET("/cert/:certId/auditlog", con.GovGetAuditLog)
		gov.GET("/cert/:certId/evidence", con.GovExportEvidencePack)
		gov.GET("/reports/audit", con.GovAuditReport)
		//gov.POST("/verify/evidence", con.GovVerifyEvidence)
		gov.GET("/stats", con.GovGetStats)
	}

	return r
}
