# Blog CLI

一个用Go语言开发的命令行博客工具，用于管理和发布博客文章。

## 功能特点

- 文章管理：支持创建、编辑、删除和查看博客文章
- 数据库支持：使用MySQL数据库存储文章数据
- 多平台支持：支持Windows、Linux和macOS（包括Intel和Apple Silicon）
- 配置管理：使用YAML配置文件进行系统配置

## 系统要求

- Go 1.21 或更高版本
- MySQL 数据库
- 支持的操作系统：
  - Windows (64-bit)
  - Linux (64-bit)
  - macOS (Intel/Apple Silicon)

## 安装

### 从源码构建

1. 克隆仓库：
```bash
git clone https://github.com/yourusername/blog-cli.git
cd blog-cli
```

2. 安装依赖：
```bash
go mod download
```

3. 构建项目：
```bash
# Windows
go build -o blog_cli.exe

# Linux/macOS
go build -o blog_cli
```

或者使用提供的构建脚本：
```bash
./build.sh
```

构建完成后，可执行文件将位于 `build` 目录中。

### 配置

1. 复制配置文件模板：
```bash
cp config.yaml.example config.yaml
```

2. 编辑 `config.yaml` 文件，配置数据库连接信息：
```yaml
database:
    host: localhost
    port: 3306
    user: go_blog
    password: go_blog
    dbname: go_blog
scan:
    dir: ""  # 文章目录路径
    workers: 0  # 扫描工作线程数
```

3. 创建数据库：
```sql
CREATE DATABASE go_blog;
```

## 使用方法

### 基本命令

```bash
# 查看帮助信息
blog_cli --help

# 创建新文章
blog_cli create "文章标题"

# 编辑文章
blog_cli edit <文章ID>

# 删除文章
blog_cli delete <文章ID>

# 列出所有文章
blog_cli list

# 查看文章详情
blog_cli view <文章ID>
```

## 项目结构

```
blog-cli/
├── articles/     # 文章存储目录
├── build/        # 编译输出目录
├── cmd/          # 命令行工具实现
├── models/       # 数据模型
├── utils/        # 工具函数
├── config.yaml   # 配置文件
├── go.mod        # Go模块文件
├── go.sum        # Go依赖版本锁定文件
└── build.sh      # 构建脚本
```

## 开发

### 目录说明

- `articles/`: 存放博客文章的目录
- `cmd/`: 包含命令行工具的主要实现代码
- `models/`: 定义数据模型和数据库操作
- `utils/`: 包含各种工具函数和辅助方法
- `build/`: 存放编译后的可执行文件
- `log/`: 存放日志文件

### 依赖管理

项目使用Go模块进行依赖管理，主要依赖包括：

- `github.com/spf13/cobra`: 命令行工具框架
- `gopkg.in/yaml.v3`: YAML配置文件解析
- `gorm.io/gorm`: ORM框架
- `gorm.io/driver/mysql`: MySQL数据库驱动

## 贡献

欢迎提交Issue和Pull Request来帮助改进这个项目。

## 许可证

MIT License 