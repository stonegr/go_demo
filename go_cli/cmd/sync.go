package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yourusername/blog-cli/models"
	"github.com/yourusername/blog-cli/utils"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dirPath    string
	configPath string
	// 数据库配置参数
	dbHost      string
	dbPort      int
	dbUser      string
	dbPassword  string
	dbName      string
	workerCount int
)

type articleTask struct {
	path     string
	relPath  string
	fileInfo os.FileInfo
}

type syncStats struct {
	created int32
	updated int32
	skipped int32
	deleted int32
	scanned int32
	errored int32
}

// 添加新的结构体用于收集数据库操作
type articleOperation struct {
	article *models.Article
	isNew   bool
}

// 添加新的结构体用于缓存文章信息
type ArticleCache struct {
	ID        int64
	MD5Check  string
	CreatedAt time.Time
	Existing  bool
}

// 添加新的结构体来保证并发安全
type SafeArticleCache struct {
	sync.RWMutex
	cache map[int64]ArticleCache
}

// 添加安全的操作方法
func (sc *SafeArticleCache) Get(id int64) (ArticleCache, bool) {
	sc.RLock()
	defer sc.RUnlock()
	val, ok := sc.cache[id]
	return val, ok
}

func (sc *SafeArticleCache) Set(id int64, value ArticleCache) {
	sc.Lock()
	defer sc.Unlock()
	sc.cache[id] = value
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync blog posts to database",
	Long:  `Sync markdown blog posts from the specified directory to the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取日志实例
		log := utils.GetLogger()

		// 初始化计时器
		startTime := time.Now()

		var cfg *utils.Config = &utils.Config{}
		if _, err := os.Stat(configPath); err == nil {
			// Load configuration
			var err error
			cfg, err = utils.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %v", err)
			}
		} else {
			panic("please run `blog_cli config` to generate config.yaml")
		}

		// 使用 Changed() 方法检查命令行参数是否被设置
		if !cmd.Flags().Changed("db-host") {
			dbHost = cfg.Database.Host
		}
		if !cmd.Flags().Changed("db-port") {
			dbPort = cfg.Database.Port
		}
		if !cmd.Flags().Changed("db-user") {
			dbUser = cfg.Database.User
		}
		if !cmd.Flags().Changed("db-password") {
			dbPassword = cfg.Database.Password
		}
		if !cmd.Flags().Changed("db-name") {
			dbName = cfg.Database.DBName
		}

		// Connect to database
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser,
			dbPassword,
			dbHost,
			dbPort,
			dbName,
		)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.New(
				log,
				logger.Config{
					SlowThreshold:             time.Second,
					LogLevel:                  logger.Error,
					IgnoreRecordNotFoundError: true,
					Colorful:                  false,
				},
			),
		})
		if err != nil {
			return fmt.Errorf("failed to connect to database: %v", err)
		}

		// Auto migrate the schema
		if err := db.AutoMigrate(&models.Article{}); err != nil {
			return fmt.Errorf("failed to migrate database: %v", err)
		}

		// 使用 Changed() 方法检查目录参数
		if !cmd.Flags().Changed("dir") {
			dirPath = cfg.Scan.Dir
		}
		if dirPath == "" {
			dirPath, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %v", err)
			}
		}

		// 使用 Changed() 方法检查 workers 参数
		if !cmd.Flags().Changed("workers") {
			workerCount = cfg.Scan.Workers
		}
		// 如果workerCount为0，则使用cpu核心数
		if workerCount == 0 {
			workerCount = runtime.NumCPU()
		}

		// 预加载所有文章信息到内存
		safeCache := &SafeArticleCache{
			cache: make(map[int64]ArticleCache),
		}
		var existingArticles []models.Article
		if err := db.Select("id, md5_check, created_at").Find(&existingArticles).Error; err != nil {
			return fmt.Errorf("failed to load existing articles: %v", err)
		}
		for _, article := range existingArticles {
			safeCache.Set(article.ID, ArticleCache{
				ID:        article.ID,
				MD5Check:  article.MD5Check,
				CreatedAt: article.CreatedAt,
				Existing:  false,
			})
		}

		// 创建任务channel和结果channel
		tasks := make(chan articleTask, 100)
		results := make(chan error, 100)
		var stats syncStats

		// 创建操作收集通道
		opChan := make(chan articleOperation, 100)

		// 启动批量处理 goroutine
		var processingWg sync.WaitGroup
		processingWg.Add(1)
		go func() {
			defer processingWg.Done()

			const batchSize = 50
			newArticles := make([]*models.Article, 0, batchSize)
			updateArticles := make([]*models.Article, 0, batchSize)

			// 定时提交批次
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			commitBatch := func() {
				if len(newArticles) > 0 {
					if err := db.Create(&newArticles).Error; err != nil {
						log.Error("Error batch creating articles: ", err)
					}
					newArticles = newArticles[:0]
				}
				if len(updateArticles) > 0 {
					for _, article := range updateArticles {
						// 使用事务批量更新
						if err := db.Save(article).Error; err != nil {
							log.Error("Error batch updating articles: ", err)
						}
					}
					updateArticles = updateArticles[:0]
				}
			}

			for {
				select {
				case op, ok := <-opChan:
					if !ok {
						// 通道关闭，提交剩余的批次
						commitBatch()
						return
					}

					if op.isNew {
						newArticles = append(newArticles, op.article)
					} else {
						updateArticles = append(updateArticles, op.article)
					}

					// 当达到批次大小时提交
					if len(newArticles) >= batchSize || len(updateArticles) >= batchSize {
						commitBatch()
					}

				case <-ticker.C:
					// 定时提交，避免数据积压
					commitBatch()
				}
			}
		}()

		// 启动worker pool
		var wg sync.WaitGroup
		for i := 0; i < workerCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for task := range tasks {
					if err := processArticle(db, task, &stats, safeCache, opChan, log); err != nil {
						results <- fmt.Errorf("error processing %s: %v", task.relPath, err)
					}
				}
			}()
		}

		// 遍历目录并发送任务
		go func() {
			err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() || !strings.HasSuffix(path, ".md") {
					return nil
				}

				atomic.AddInt32(&stats.scanned, 1)
				relPath, _ := filepath.Rel(dirPath, path)
				tasks <- articleTask{path: path, relPath: relPath, fileInfo: info}
				return nil
			})

			if err != nil {
				results <- fmt.Errorf("failed to walk directory: %v", err)
			}
			close(tasks)
		}()

		// 等待所有worker完成
		go func() {
			wg.Wait()
			close(results)
		}()

		// 等待所有 worker 完成后关闭操作通道
		wg.Wait()
		close(opChan)

		// 等待批处理完成
		processingWg.Wait()

		// 处理错误
		for err := range results {
			if err != nil {
				atomic.AddInt32(&stats.errored, 1)
				log.Error("Error: ", err)
			}
		}

		// 删除不存在的文章
		deleteNonExistentArticles(db, safeCache, &stats, log)

		// 计算运行时间
		duration := time.Since(startTime)

		// 打印统计信息
		log.Info("Sync completed successfully!")
		log.Info("Statistics:")
		log.Infof("- Files scanned: %d", stats.scanned)
		log.Infof("- New articles: %d", stats.created)
		log.Infof("- Updated articles: %d", stats.updated)
		log.Infof("- Deleted articles: %d", stats.deleted)
		log.Infof("- Skipped (no changes): %d", stats.skipped)
		log.Infof("- Errors: %d", stats.errored)
		log.Infof("- Total execution time: %v", duration)

		return nil
	},
}

// 修改 processArticle 函数签名，添加日志参数
func processArticle(db *gorm.DB, task articleTask, stats *syncStats, articleCache *SafeArticleCache, opChan chan<- articleOperation, log *logrus.Logger) error {
	content, err := os.ReadFile(task.path)
	if err != nil {
		atomic.AddInt32(&stats.errored, 1)
		log.Error("Error reading file ", task.path, ": ", err)
		return fmt.Errorf("failed to read file: %v", err)
	}

	hash := md5.Sum(content)
	md5Hash := hex.EncodeToString(hash[:])

	article, err := parseArticle(content)
	if err != nil {
		atomic.AddInt32(&stats.errored, 1)
		log.Error("Error parsing article ", task.path, ": ", err)
		return fmt.Errorf("failed to parse article: %v", err)
	}

	// 修改文章缓存的处理逻辑
	var isNew bool
	if existingArticle, ok := articleCache.Get(article.ID); ok {
		if existingArticle.MD5Check == md5Hash {
			atomic.AddInt32(&stats.skipped, 1)
			articleCache.Set(article.ID, ArticleCache{
				ID:        article.ID,
				MD5Check:  md5Hash,
				CreatedAt: existingArticle.CreatedAt,
				Existing:  true,
			})
			return nil
		}
		article.CreatedAt = existingArticle.CreatedAt
		article.UpdatedAt = time.Now()
		atomic.AddInt32(&stats.updated, 1)
		log.Info("Updated: ", task.relPath)

		articleCache.Set(article.ID, ArticleCache{
			ID:        article.ID,
			MD5Check:  md5Hash,
			CreatedAt: existingArticle.CreatedAt,
			Existing:  true,
		})
		isNew = false
	} else {
		article.CreatedAt = time.Now()
		article.UpdatedAt = time.Now()
		atomic.AddInt32(&stats.created, 1)
		log.Info("Created: ", task.relPath)

		articleCache.Set(article.ID, ArticleCache{
			ID:        article.ID,
			MD5Check:  md5Hash,
			CreatedAt: article.CreatedAt,
			Existing:  true,
		})
		isNew = true
	}

	article.MD5Check = md5Hash
	// 将操作发送到通道
	opChan <- articleOperation{
		article: article,
		isNew:   isNew,
	}
	return nil
}

func deleteNonExistentArticles(db *gorm.DB, articleCache *SafeArticleCache, stats *syncStats, log *logrus.Logger) {
	var idsToDelete []int64

	// 检查缓存中未标记为存在的文章
	for id := range articleCache.cache {
		if !articleCache.cache[id].Existing {
			idsToDelete = append(idsToDelete, id)
		}
	}

	if len(idsToDelete) > 0 {
		// 获取要删除的文章的标题，用于日志记录
		var articles []models.Article
		if err := db.Where("id IN ?", idsToDelete).Find(&articles).Error; err == nil {
			for _, article := range articles {
				log.Infof("Deleted: ID=%d, Title=\"%s\"", article.ID, article.Title)
			}
		}

		if err := db.Delete(&models.Article{}, idsToDelete).Error; err != nil {
			atomic.AddInt32(&stats.errored, int32(len(idsToDelete)))
			log.Error("Error: Failed to delete articles: ", err)
		} else {
			atomic.AddInt32(&stats.deleted, int32(len(idsToDelete)))
		}
	}
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// 文件扫描目录配置
	syncCmd.Flags().StringVarP(&dirPath, "dir", "d", "", "Directory to scan for markdown files")
	syncCmd.Flags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to configuration file")

	// 数据库配置
	syncCmd.Flags().StringVar(&dbHost, "db-host", "localhost", "Database host")
	syncCmd.Flags().IntVar(&dbPort, "db-port", 3306, "Database port")
	syncCmd.Flags().StringVar(&dbUser, "db-user", "go_blog", "Database user")
	syncCmd.Flags().StringVar(&dbPassword, "db-password", "go_blog", "Database password")
	syncCmd.Flags().StringVar(&dbName, "db-name", "go_blog", "Database name")

	// 添加worker数量配置
	syncCmd.Flags().IntVarP(&workerCount, "workers", "w", 5, "Number of concurrent workers, 0 means use all cores")
}

func parseArticle(content []byte) (*models.Article, error) {
	// Split content into front matter and markdown
	parts := strings.Split(string(content), "---")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid article format: missing front matter")
	}

	// Parse front matter
	var frontMatter struct {
		ID           int64    `yaml:"id"`
		Title        string   `yaml:"title"`
		Tags         []string `yaml:"tags"`
		Cover        string   `yaml:"cover"`
		Excerpt      string   `yaml:"excerpt"`
		Listed       bool     `yaml:"listed"`
		DateModified string   `yaml:"dateModified"`
	}

	// 预处理 front matter，为包含冒号的标题添加引号
	frontMatterStr := parts[1]
	lines := strings.Split(frontMatterStr, "\n")
	for i, line := range lines {
		// 找到第一个冒号后的内容
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			parts2 := strings.TrimSpace(parts[1])
			// 如果标题包含冒号且没有被引号包围，则添加双引号
			if strings.Contains(parts2, ":") && !strings.HasPrefix(parts2, "\"") {
				lines[i] = fmt.Sprintf("%s: \"%s\"", parts[0], parts2)
			}
		}
	}
	frontMatterStr = strings.Join(lines, "\n")

	if err := yaml.Unmarshal([]byte(frontMatterStr), &frontMatter); err != nil {
		return nil, fmt.Errorf("failed to parse front matter: %v", err)
	}

	// Parse date
	dateModified, err := time.Parse("2006-01-02", frontMatter.DateModified)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %v", err)
	}

	// Create article
	article := &models.Article{
		ID:           frontMatter.ID,
		Title:        frontMatter.Title,
		Tags:         strings.Join(frontMatter.Tags, ","),
		Cover:        frontMatter.Cover,
		Excerpt:      frontMatter.Excerpt,
		Listed:       frontMatter.Listed,
		DateModified: dateModified,
		Content:      strings.TrimSpace(parts[2]),
	}

	return article, nil
}
