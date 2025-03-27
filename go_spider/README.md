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
	// 数据库配置
	cfg.DB.Host = "127.0.0.1"
	cfg.DB.Port = "33306"
	cfg.DB.User = "go_blog"
	cfg.DB.Password = "go_blog"
	cfg.DB.Name = "go_blog"
	cfg.DB.Table = "posts"

	// 爬虫配置
	cfg.Scraper.BaseURL = "https://www.30secondsofcode.org/%s/p/%d/"
	cfg.Scraper.Categories = []string{"js", "css", "html", "react", "node", "git", "python"}  //抓取的分类名称
	cfg.Scraper.CatchPages = 1                                                                // 分类下抓取的页数
	cfg.Scraper.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
```

## 数据库结构

在运行前需先创建以下数据表：

```sql
CREATE TABLE `post`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '文章ID',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NOT NULL COMMENT '文章标题',
  `summary` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT NULL COMMENT '文章摘要',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL COMMENT '文章内容',
  `category` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT NULL COMMENT '文章分类',
  `publish_time` date NOT NULL COMMENT '发布时间',
  `image_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NULL DEFAULT NULL COMMENT '文章配图URL',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_posts_deleted_at`(`deleted_at`) USING BTREE,
  INDEX `idx_posts_category`(`category`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 140 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_german2_ci ROW_FORMAT = Dynamic;
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