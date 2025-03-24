# Go Blog

一个使用 Go 语言开发的现代化博客系统，基于 Gin 框架和 GORM ORM。

## 功能特性

- 基于 Gin 框架的高性能 Web 服务器
- 使用 GORM 进行数据库操作
- 支持 MySQL 数据库
- 使用 Viper 进行配置管理
- 模块化的项目结构
- RESTful API 设计

## 技术栈

- Go 1.20+
- Gin Web 框架
- GORM ORM
- MySQL 数据库
- Viper 配置管理

## 项目结构

```
go_blog/
├── config/         # 配置文件和配置管理
├── controllers/    # 控制器层，处理请求逻辑
├── models/         # 数据模型层
├── routes/         # 路由配置
├── templates/      # 模板文件
├── main.go         # 应用程序入口
└── go.mod          # Go 模块定义文件
```

## 安装说明

1. 确保已安装 Go 1.20 或更高版本
2. 克隆项目到本地：
   ```bash
   git clone [项目地址]
   cd go_blog
   ```
3. 安装依赖：
   ```bash
   go mod download
   ```
4. 配置数据库：
   - 在 `config` 目录下配置数据库连接信息
   - 确保 MySQL 服务已启动

## 运行项目

1. 启动服务器：
   ```bash
   go run main.go
   ```
2. 服务器将在配置的地址和端口上启动（默认为 localhost:8080）

## 开发说明

- 项目使用 Go modules 进行依赖管理
- 遵循 MVC 架构模式
- 使用 GORM 进行数据库操作
- 使用 Gin 框架处理 HTTP 请求

## 贡献指南

欢迎提交 Issue 和 Pull Request 来帮助改进项目。

## 许可证

[MIT License](LICENSE) 