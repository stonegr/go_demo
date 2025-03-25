# Go Blog

一个使用 Go 语言开发的现代化博客系统，采用 Gin 框架和 GORM ORM。

## 技术栈

- **后端框架**: Gin v1.9.1
- **数据库**: MySQL 8.0
- **ORM**: GORM v1.25.7
- **配置管理**: Viper v1.18.2
- **日志系统**: Logrus v1.9.3
- **开发语言**: Go 1.20+

## 项目结构

```
.
├── config/         # 配置文件目录
├── controllers/    # 控制器层，处理请求逻辑
├── models/        # 数据模型层
├── routes/        # 路由配置
├── templates/     # 模板文件
├── utils/         # 工具函数
├── logs/          # 日志文件
├── main.go        # 程序入口
├── go.mod         # Go 模块文件
├── go.sum         # Go 依赖版本锁定文件
├── Dockerfile     # Docker 构建文件
└── docker-compose.yml  # Docker 编排文件
```

## 快速开始

### 环境要求

- Go 1.20 或更高版本
- Docker 和 Docker Compose（可选，用于容器化部署）
- MySQL 8.0（如果使用本地数据库）

### 本地开发

1. 克隆项目
```bash
git clone [项目地址]
cd go_blog
```

2. 安装依赖
```bash
go mod download
```

3. 配置数据库
- 创建 MySQL 数据库
- 修改 `config/config.yaml.sample` 中的数据库配置

4.环境变量
- BLOG_ENV 值为DEV时会读取`config/config-dev.yaml`作为配置文件
- OPENAI_ENV 值为DEV时会使用代码中的mock数据
```
    "BLOG_DATABASE_HOST": "localhost",
    "BLOG_DATABASE_PORT": "3306",
    "BLOG_DATABASE_USER": "root",
    "BLOG_DATABASE_PASSWORD": "123456",
    "BLOG_DATABASE_NAME": "go_blog"
```
如上用`BLOG`开头和`_`连接的环境变量也可覆盖配置文件中的数值

5. 运行项目
```bash
go run main.go
```

### Docker 部署

1. 构建并启动服务
```bash
docker-compose up -d
```

2. 查看服务状态
```bash
docker-compose ps
```

3. 查看日志
```bash
docker-compose logs -f app
```

## 开发指南

### 目录说明

- `controllers/`: 包含所有控制器，处理 HTTP 请求和响应
- `models/`: 定义数据模型和数据库结构
- `routes/`: 配置 API 路由和中间件
- `config/`: 存放配置文件
- `utils/`: 通用工具函数和辅助方法
- `templates/`: 视图模板文件
- `logs/`: 应用日志文件

### 开发流程

1. 创建新的功能分支
```bash
git checkout -b feature/your-feature-name
```

2. 开发新功能
- 在 `models/` 中定义数据模型
- 在 `controllers/` 中实现业务逻辑
- 在 `routes/` 中配置路由
- 在 `templates/` 中添加视图模板

3. 提交代码
```bash
git add .
git commit -m "feat: your feature description"
git push origin feature/your-feature-name
```

### 代码规范

- 遵循 Go 标准代码规范
- 使用有意义的变量和函数命名
- 添加必要的注释和文档
- 确保代码经过测试

## 配置说明

主要配置文件位于 `config/config.yaml`，包含：

- 数据库配置
- 服务器配置
- 日志配置
- 其他系统配置

## 日志系统

- 使用 Logrus 进行日志管理
- 日志文件位于 `logs/` 目录
- 支持不同级别的日志记录
- 包含时间戳和上下文信息

## 部署说明

### 生产环境部署

1. 修改配置文件
- 更新数据库连接信息
- 配置正确的环境变量
- 设置适当的日志级别

2. 构建 Docker 镜像
```bash
docker build -t go-blog .
```

3. 使用 Docker Compose 部署
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### 监控和维护

- 定期检查日志文件
- 监控数据库性能
- 备份重要数据
- 更新依赖包

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License

Copyright (c) 2024 [mengxuan]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.