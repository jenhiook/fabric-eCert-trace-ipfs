# 基于区块链的电子证照存储与溯源系统

本项目基于 Hyperledger Fabric + IPFS + Go + Vue 实现，完成电子证照从创建、审核、备案、核验到全生命周期溯源的全流程管理。支持个人、政务、企业、技术四个角色，并提供供应商关联管理、合规事件存证、可视化分析及审计报表导出等功能。

## 技术栈
- 区块链：Hyperledger Fabric v2.5.6
- 链码：Go 1.22
- 分布式存储：IPFS Kubo 0.28.0
- 后端：Go + Gin
- 前端：Vue 2 + Element UI + ECharts
- 数据库：MySQL 8.0

## 环境要求
- Ubuntu 22.04 (或 WSL2)
- Docker & Docker Compose
- Go 1.22+
- Node.js 16+
- IPFS 已安装 (`ipfs` 命令可用)

## 快速启动（一键脚本）
```bash
# 确保 Docker 已启动
sudo systemctl start docker

# 执行一键启动脚本（会自动启动 MySQL、IPFS、Fabric 网络、后端、前端）
bash start-eCert.sh

启动后访问：http://localhost:9528

预置测试账号
角色	       账号	    密码
个人用户	   wrf	    123456
政务部门	   gov001	123456
企业组织	   ent001	123456
技术支撑实体	tech001	 123456
主要功能
个人用户：证照上链、溯源时间轴、独立时间轴（CSV导出）、个人可视化分析

政务部门：证照审核、政府溯源（分页/时间区间/导出）、政府可视化分析、对外核验接口

企业组织：证照备案、供应商关联管理（链上存证）、合规事件存证（IPFS附件）、企业可视化分析

技术实体：证照核验、核验工具（证据真实性验证）、技术可视化分析

目录结构
text
.
├── application/       # 后端与前端代码
├── blockchain/        # 链码与Fabric网络配置
├── start-eCert.sh     # 一键启动脚本
└── README.md
注意事项
首次启动需拉取Docker镜像，可能较慢。

确保 IPFS 守护进程正常运行。

如遇权限问题，执行 chmod +x start-eCert.sh。

许可证
Apache-2.0