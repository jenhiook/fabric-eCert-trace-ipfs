#  基于区块链的电子证照存储与溯源系统 FabricV2.5溯源系统课堂活动 fabric-eCert-trace-ipfs

B站@[ammlhguhj](https://space.bilibili.com/9911299)

[项目开源链接](https://gitee.com/surgar2022/fabric-eCert-trace-ipfs)

2024/07/15
## 介绍

本项目基于Hyperledger Fabric V2.5，实现了一个电子证照存储与溯源系统。在本区块链系统中，有5个内置的角色：个人用户、政务部门、企业组织、技术支撑实体、其他相关实体。其中个人用户、政务部门、企业组织、技术支撑实体可以将信息上链，其他相关实体有信息溯源权限。

## 特别说明

**Fabric V2.5通用溯源项目讲解与二次开发课堂活动作品**

请注意，本指南中的内容仅用于演示目的。群主大大强调，所有付费内容严禁分享、公开与倒卖！如果您对Fabric V2.5通用溯源项目的深入讲解和二次开发感兴趣，并希望获取完整的教程和支持，请尊重知识产权，通过正规渠道自行购买相关课程或服务。

原始项目地址：[fabric-trace](https://github.com/TrueTechLabs/fabric-trace)

## 项目介绍

**采用的技术栈**：Fabric V2.5、Gin、Vue.js、Mysql

**开发工具**：1Panel, CentOS, MySQL, VSCode, Docker

**技术架构**：Hyperledger Fabric V2.5


## 项目背景
**指导方针**：国务院办公厅推动政务服务标准化、规范化、便利化。

**政策背景**：加快实现政务服务线上线下一体化。

**电子证照应用需求**：现存问题包括数据共享不及时、监管缺失、电子证照真实性验证困难。

**传统存储方式**：中心化存储、数据易篡改、存储量大。

**解决方案**：利用区块链技术实现数据溯源与存储。

**目标**：实现跨部门、跨区域电子证照互认共享。


## 项目工作与技术实现

## 官方环境要求
- [Hyperledger Fabric V2.5 环境要求](https://hyperledger-fabric.readthedocs.io/en/release-2.5/prereqs.html)
- [Hyperledger Fabric V2.5 安装教程](https://hyperledger-fabric.readthedocs.io/en/release-2.5/install.html)
## 系统环境
- 操作系统：CentOS 7.9 64 bit

- CPU：2核

- 内存：2G

- 硬盘：40G SSD
![输入图片说明](pic/%E5%9B%BE1%20%E7%B3%BB%E7%BB%9F%E7%8E%AF%E5%A2%83.jpg)


## 安装步骤
### 安装 Linux 控制面板
使用1Panel进行环境的安装与管理，它会自动安装Docker，简化安装步骤。
```sh
curl -sSL https://resource.fit2cloud.com/1panel/package/quick_start.sh -o quick_start.sh && sh quick_start.sh
```

测试 Docker 是否安装成功：
```sh
docker -version
```
### 配置防火墙

```sh
#安装 firewalld
sudo yum install firewalld
#启动 firewalld
sudo systemctl start firewalld
```

### 安装 MySQL
在 1Panel 应用商店中安装 MySQL 8.2.0，并设置端口号为 3337，开启外部端口访问。
![输入图片说明](pic/%E5%9B%BE2%20mysql%E6%95%B0%E6%8D%AE%E5%BA%93%E5%92%8CphpMyadmin%E6%95%B0%E6%8D%AE%E5%BA%93%E7%AE%A1%E7%90%86%E7%B3%BB%E7%BB%9F.jpg)

![输入图片说明](pic/%E5%9B%BE3%20%E6%95%B0%E6%8D%AE%E5%BA%93%E7%9B%B8%E5%85%B3%E4%BF%A1%E6%81%AF.jpg)
确保配置mysql为：
```sh
#127.0.0.1
#3337
#名称 sugar20240713
#用户名 sugar20240713
#密码 sugar20240713
```

### 克隆项目
```sh
git clone https://gitee.com/surgar2022/fabric-eCert-trace-ipfs
cd fabric-eCert-trace-ipfs
```
### 安装依赖
#### 安装 Go
```sh
sudo yum install golang
go version
```
### 安装Node.js环境（使用NVM）

#### 安装 NVM 和 Node.js

```sh
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
source ~/.bashrc
nvm install 16
nvm use 16
```
测试 Node.js 和 npm：
```sh
node -v
npm -v
```
#### 安装 jq
```sh
sudo yum install jq
jq --version
```
### 防火墙设置
在 1Panel 中放行以下 TCP 端口：8080, 9090, 9528, 3306，3337。

1Panel防火墙放行端口

3337： (mysql)	

9528：web前端，

9090 web后端，

8080 区块链浏览器

![输入图片说明](pic/%E5%9B%BE4%20%E9%98%B2%E7%81%AB%E5%A2%99%E8%AE%BE%E7%BD%AE.jpg)

### 启动区块链网络
使用cURL命令下载install-fabric.sh脚本，并为其赋予执行权限。

下载 Fabric 示例、Docker 镜像和二进制文件


```sh
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
./install-fabric.sh -h
./install-fabric.sh --fabric-version 2.5.9 docker samples binary
```
执行上述命令后，Fabric示例、Docker镜像和二进制文件将被下载并安装到您的系统中。

在 `fabric-trace/blockchain/network` 目录下执行：
```sh
# 启动区块链网络
./start.sh
```
![输入图片说明](pic/%E5%9B%BE5%20%E5%90%AF%E5%8A%A8%E5%8C%BA%E5%9D%97%E9%93%BE%E7%BD%91%E7%BB%9C.jpg)
### 启动后端

#### 修改后端mysql端口

在 `fabric-trace/application/backend` 找到config.yaml修改：

![输入图片说明](pic/%E5%9B%BE6%20%E4%BF%AE%E6%94%B9%E5%90%8E%E7%AB%AF%E6%95%B0%E6%8D%AE%E5%BA%93%E8%AE%BE%E7%BD%AE%E6%96%87%E4%BB%B6.jpg)

在 `fabric-trace/application/backend` 目录下执行：
```sh
go run main.go
```
修改后端 IP，替换以下文件中的 IP 地址为您的云服务 IP：

```sh
find . -type f -exec grep -l "107.175.32.34" {} \;
./application/replaceip.sh
./application/web/.env.staging
./application/web/.env.development
./application/web/.env.production
./application/web/src/router/index.js

```

### 启动前端
新开一个 SSH 窗口，在 `fabric-trace/application/web` 目录下执行：
```sh
# 安装依赖（仅在首次运行时执行）
npm install
# 启动前端
npm run dev

```
![输入图片说明](pic/%E5%9B%BE8%20%E5%90%AF%E5%8A%A8%E5%89%8D%E7%AB%AF.jpg)

### 访问系统

在浏览器中打开：`http://云服务器IP:9528`，即可看到前端页面。
#### 注册登录
在访问系统后，您可以进行注册：
1. 填写注册表单。
2. 设置您的用户名和密码。
3. 提交注册信息。
![输入图片说明](pic/%E5%9B%BE10%20%E6%B3%A8%E5%86%8C%E7%95%8C%E9%9D%A2.jpg)
![输入图片说明](pic/%E5%9B%BE11%20%E8%A1%A8%E5%8D%95%E7%95%8C%E9%9D%A2.jpg)
登录后，您可以进行以下操作：
1. 填写溯源信息表单：输入电子证照的相关信息，如类型、持有者信息等。
2. 上传图片文件：附加相关的图片文件以证明电子证照的真实性。
3. 提交信息：完成表单后，提交信息以存储到区块链上。
![输入图片说明](pic/%E5%9B%BE12%20%E5%8F%B3%E4%B8%8A%E8%A7%92%E6%8E%A8%E5%B9%BF%E7%95%8C%E9%9D%A2.jpg)

![输入图片说明](pic/%E5%9B%BE13%20%E8%8E%B7%E5%8F%96%E6%89%80%E6%9C%89%E7%94%B5%E5%AD%90%E8%AF%81%E7%85%A7%E4%BF%A1%E6%81%AF%E7%95%8C%E9%9D%A2.jpg)
#### 区块链浏览器
系统还提供了区块链浏览器功能，您可以通过以下方式使用：
1. 查询交易：查看区块链上的交易记录。
2. 验证信息：通过交易哈希或区块编号验证电子证照的溯源信息。
请确保在演示过程中所有服务都正常运行，以便用户能够顺畅地体验系统的各项功能。
![输入图片说明](pic/%E5%9B%BE14%20%E5%8C%BA%E5%9D%97%E9%93%BE%E6%B5%8F%E8%A7%88%E5%99%A8%E7%95%8C%E9%9D%A2.jpg)

![输入图片说明](pic/%E5%9B%BE15%20%E5%8C%BA%E5%9D%97%E9%93%BE%E7%BD%91%E7%BB%9C%E7%95%8C%E9%9D%A2.jpg)
## 区块链业务网络

![输入图片说明](pic/%E5%9B%BE15%20%E5%8C%BA%E5%9D%97%E9%93%BE%E7%BD%91%E7%BB%9C%E7%95%8C%E9%9D%A2.jpg)

## 配置nginx

- 购买域名，配置 `nginx.conf`。
  示例 `nginx.conf`：
```nginx
events {
  worker_connections  1024;
}
http {
  server {
    listen 80;
    server_name www.dovis.icu;
    location / {
      proxy_pass http://127.0.0.1:9528;
      proxy_set_header Host $host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
  }
}
```
## 总结
本教程详细介绍了从服务器环境准备到项目部署的整个流程，为用户提供了在 CentOS 系统上部署基于 Hyperledger Fabric 的农产品溯源系统的步骤。教程中涉及了多个技术栈的安装和配置，包括 Docker、MySQL、Go、Node.js、NVM 和 `jq`，以及必要的防火墙设置。通过这个教程，用户可以了解到如何在海外服务器上搭建一个完整的区块链溯源系统，并且为后续的系统优化和扩展提供了方向。