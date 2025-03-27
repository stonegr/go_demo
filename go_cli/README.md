# Blog CLI

一个用于同步博客文章到数据库的命令行工具。该工具可以跟踪Markdown博客文章的变化，并只更新已修改的文件。

## 功能特点

- 支持Markdown文件的同步
- 智能追踪文件变化
- 增量更新，只同步修改过的文件
- 跨平台支持（Windows、Linux、macOS）
- 配置灵活，支持自定义设置

## 项目结构

```
.
├── articles/     # 存放Markdown博客文章
├── build/        # 编译后的二进制文件
├── cmd/          # 命令行工具的核心代码
├── models/       # 数据模型定义
├── utils/        # 工具函数
├── main.go       # 程序入口
└── build.sh      # 构建脚本
```

## 系统要求

- Go 1.16 或更高版本
- 支持的操作系统：
  - Windows (64-bit)
  - Linux (64-bit)
  - macOS (Intel/ARM)

## 安装

### 从源码构建

1. 克隆仓库：
```bash
git clone https://github.com/stonegr/blog-cli.git
cd blog-cli
```

2. 运行构建脚本：
```bash
./build.sh
```

构建完成后，可执行文件将位于 `build` 目录中。

### 直接下载

你可以从 releases 页面下载预编译的二进制文件。

## 配置

编辑 `config.yaml` 文件来配置你的设置：

```yaml
# 数据库配置
database:
  host: localhost
  port: 5432
  user: your_username
  password: your_password
  dbname: your_database

# 同步设置
scan:
    dir: "" # 需要扫描的目录
    workers: 0 # 并发数量
```

## 使用方法

1. 生成配置文件：
```bash
./blog_cli config
```
> 请根据config.yaml格式修改数据库相关信息

2. 生成示例文章：
```bash
./blog_cli generate -c 1000
```
> 上述代码会默认在articles文件夹生成1000个示例markdown文件

3. 同步文章：
```bash
./blog_cli sync
```

4. 查看帮助：
```bash
blog_cli --help
```

## 开发

### 本地开发

1. 安装依赖：
```bash
go mod download
```

2. 本地构建：
```bash
go build
```

### 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 作者

- 你的名字 - [@stonegr](https://github.com/stonegr)

## 致谢

- [Cobra](https://github.com/spf13/cobra) - 用于构建命令行应用
- 其他使用的开源项目 