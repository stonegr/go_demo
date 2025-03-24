# Go Spider

这是一个使用 Colly 框架实现的网页爬虫程序，可以爬取指定网站的文章内容并保存到 MySQL 数据库中。

## 功能特点

- 支持按分类爬取文章
- 自动提取文章标题、摘要、分类、发布时间
- 保存文章正文 HTML 内容
- 数据持久化到 MySQL 数据库
- 内置请求延迟，避免对目标网站造成压力

## 环境要求

- Go 1.20 或更高版本
- MySQL 数据库

## 安装依赖

```bash
go mod tidy
```

## 配置说明

在使用之前，需要修改 `main.go` 中的数据库配置：

```go
const (
    dbUser     = "root"         // 修改为你的数据库用户名
    dbPassword = "your_password" // 修改为你的数据库密码
    dbName     = "spider_db"    // 修改为你的数据库名
)
```

同时，需要修改爬虫的目标网站和分类：

```go
baseURL := "https://example.com" // 修改为要爬取的网站URL
categories := []string{"category1", "category2"} // 修改为要爬取的分类
```

## 数据库结构

程序会自动创建以下数据表：

```sql
CREATE TABLE articles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    summary TEXT,
    category VARCHAR(50),
    publish_time VARCHAR(50),
    content LONGTEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
```

## 使用方法

1. 确保 MySQL 数据库已经启动
2. 修改配置信息
3. 运行程序：

```bash
go run main.go
```

## 注意事项

- 请确保遵守目标网站的爬虫规则
- 建议适当调整请求延迟时间，避免对目标网站造成压力
- 如果遇到反爬虫机制，可能需要调整 User-Agent 或添加其他请求头 