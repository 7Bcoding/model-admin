# LLM-OPS 管理平台

LLM-OPS 是一个为LLM API设计的运维与管理平台，可以对模型部署、监控和运维进行统一管理。

## 项目特性

- **模型管理**：管理各种LLM模型，后端endpoint配置/查看，观测面板的链接跳转等
- **Kubernetes 连接**：支持nitor / nebula 的主要cr资源的查看，以及主要属性的设置功能
- **用户认证**：内置用户认证系统，支持飞书登录
- **审计日志**：完整的操作审计记录

## 架构设计

系统由以下主要组件构成：

- **前端**：基于现代Web技术的管理界面
- **后端API**：Go语言实现的RESTful API服务
- **数据库**：存储用户、审计日志等系统数据
- **Kubernetes集成**：通过自定义资源与集群交互

## 技术栈

- **后端**：Go、Gorilla Mux、GORM
- **前端**：现代JavaScript/TypeScript框架
- **容器化**：Docker、Kubernetes
- **数据库**：MySQL

## 快速开始

### 前置条件
- Kubernetes 集群 (用于生产环境)
- MySQL 数据库

### 本地开发环境

1. 克隆仓库

```bash
git clone <仓库URL>
cd llm-admin
```

2. 配置环境

准备本地Mysql数据库:

```
参考sql目录中的建表语句初始化本地DB
手动将个人账户加入到user表中，并设置role为admin:
INSERT INTO llm_admin.users (username, account_name, password_hash, role) 
VALUES ('{name}', '{name}', 'hashed_password_123', 'admin');
```

创建配置文件:

```bash
cp config.yaml.template config/config.yaml
# 编辑配置文件，填入数据库信息
```

3. 启动后端开发环境

```bash
./scripts/start-dev.sh
```

4. 启动前端开发环境

```bash
cd frontend
./scripts/start-dev.sh
```


5. 构建后端镜像

``` 
bash build.sh -p -v {tag}
``` 

6. 构建前端镜像

```
cd frontend
bash build.sh -p -v {tag}
``` 



