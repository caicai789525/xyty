# 芗韵同音 (xyty) 后端项目

## 项目简介
芗韵同音是一个基于 Go 语言开发的后端服务，提供用户认证、个人中心、视频管理和情景记录等功能。

## 技术栈
- **语言**: Go
- **Web 框架**: Gin
- **数据库**: MySQL
- **云存储**: 七牛云
- **API 文档**: Swagger

## 目录结构
```
xyty-backend/
├── config/          # 配置文件
├── dao/             # 数据访问层
├── docs/            # API 文档
├── handler/         # 路由处理器
├── model/           # 数据模型
├── pkg/             # 公共包
├── router/          # 路由配置
├── services/        # 业务逻辑层
├── main.go          # 入口文件
├── go.mod           # Go 模块文件
└── go.sum           # 依赖校验文件
```

## 快速开始

### 前置条件
- Go 1.18+
- MySQL
- 七牛云账号（用于视频存储）

### 安装与运行
1. 克隆项目
   ```bash
   git clone <项目地址>
   cd xyty-backend
   ```

2. 安装依赖
   ```bash
   go mod tidy
   ```

3. 配置文件
   编辑 `config/config.yaml` 文件，填写数据库连接信息和七牛云配置。

4. 启动服务
   ```bash
   go run main.go
   ```
   服务将在 `http://localhost:8080` 上运行。

## API 文档
启动服务后，可访问以下地址查看 API 文档：
- Swagger UI: `http://localhost:8080/swagger/index.html`
- API 文档 JSON: `http://localhost:8080/swagger/doc.json`

## 功能模块

### 1. 用户认证
- 注册 (`POST /api/v1/signup`)
- 密码登录 (`POST /api/v1/pwd_login`)
- 发送验证码 (`POST /api/v1/send_mail`)
- 验证码登录 (`POST /api/v1/code_login`)

### 2. 个人中心
- 获取个人资料 (`GET /api/v1/user/profile`)
- 更新头像 (`POST /api/v1/user/avatar`)
- 修改密码 (`PUT /api/v1/user/password`)

### 3. 视频管理
- 上传图片 (`POST /upload`)
- 获取视频记录 (`GET /api/v1/user/video-records`)
- 添加视频记录 (`POST /api/v1/user/video-records`)
- 获取七牛云视频列表 (`GET /api/v1/user/qiniu-videos`)
- 获取所有视频 (`GET /api/v1/user/videos`)

### 4. 情景记录
- 获取情景记录 (`GET /api/v1/user/scenario-records`)
- 添加情景记录 (`POST /api/v1/user/scenario-records`)

## 部署说明

### 生产环境部署
1. 构建可执行文件
   ```bash
   go build -o oyhx-backend main.go
   ```

2. 配置环境变量
   根据生产环境需求配置数据库连接和七牛云信息。

3. 启动服务
   ```bash
   ./oyhx-backend
   ```

## 贡献指南
1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 开启 Pull Request

## 许可证
本项目采用 MIT 许可证。

## 联系我们
- 联系人: KitZhangYs
- 邮箱: SJMbaiyang@163.com