package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly/v2"
	"golang.org/x/time/rate"
)

// Config 存储应用配置
type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		Table    string
	}
	Scraper struct {
		BaseURL    string
		Categories []string
		CatchPages int
		UserAgent  string
	}
}

// Article 存储文章信息
type Article struct {
	Title       string
	Summary     string
	Category    string
	Content     string
	ImageURL    string
	PublishTime string
}

// Scraper 爬虫结构体
type Scraper struct {
	config    *Config
	db        *sql.DB
	collector *colly.Collector
	limiter   *rate.Limiter
}

// 初始化配置
func newConfig() *Config {
	cfg := &Config{}

	// 数据库配置
	cfg.DB.Host = "127.0.0.1"
	cfg.DB.Port = "33306"
	cfg.DB.User = "go_blog"
	cfg.DB.Password = "go_blog"
	cfg.DB.Name = "go_blog"
	cfg.DB.Table = "posts"

	// 爬虫配置
	cfg.Scraper.BaseURL = "https://www.30secondsofcode.org/%s/p/%d/"
	cfg.Scraper.Categories = []string{"js", "css", "html", "react", "node", "git", "python"}
	cfg.Scraper.CatchPages = 1
	cfg.Scraper.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	return cfg
}

// 初始化数据库连接
func initDB(cfg *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

// NewScraper 创建新的爬虫实例
func NewScraper(cfg *Config, db *sql.DB) *Scraper {
	c := colly.NewCollector(
		colly.UserAgent(cfg.Scraper.UserAgent),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,
	})

	return &Scraper{
		config:    cfg,
		db:        db,
		collector: c,
		limiter:   rate.NewLimiter(rate.Every(200*time.Millisecond), 1),
	}
}

// saveArticle 保存文章到数据库
func (s *Scraper) saveArticle(article Article) error {
	insertSQL := `
		INSERT INTO ` + s.config.DB.Table + ` (title, summary, category, publish_time, content, image_url)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(insertSQL,
		article.Title,
		article.Summary,
		article.Category,
		article.PublishTime,
		article.Content,
		article.ImageURL)

	if err != nil {
		return fmt.Errorf("failed to insert article: %v", err)
	}

	log.Printf("Successfully saved article: %s", article.Title)
	return nil
}

// setupCollector 设置爬虫规则
func (s *Scraper) setupCollector() {
	// 处理文章详情页
	s.collector.OnHTML("body > main > article", func(e *colly.HTMLElement) {
		content, _ := e.DOM.Html()
		article, ok := e.Request.Ctx.GetAny("article").(Article)
		if !ok {
			log.Println("Failed to get article from context")
			return
		}

		article.Content = content
		if err := s.saveArticle(article); err != nil {
			log.Printf("Error saving article: %v", err)
		}
	})

	// 处理文章列表页
	s.collector.OnHTML("body > main > section.preview-list", func(e *colly.HTMLElement) {
		if !strings.Contains(e.Request.URL.Path, "/p/") {
			return
		}

		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			article := Article{
				Title:       el.ChildText("h3>a"),
				Summary:     el.ChildText("p"),
				Category:    strings.Split(el.ChildText("article > small"), " · ")[0],
				PublishTime: el.ChildAttr("article > small > time", "datetime"),
				ImageURL:    "https://www.30secondsofcode.org" + el.ChildAttr("img", "src"),
			}

			link := "https://www.30secondsofcode.org" + el.ChildAttr("a", "href")
			if link == "" {
				return
			}

			ctx := colly.NewContext()
			ctx.Put("article", article)

			s.limiter.Wait(context.Background())
			if err := s.collector.Request("GET", link, nil, ctx, nil); err != nil {
				log.Printf("Error visiting article: %v", err)
			}
		})
	})

	// 错误处理
	s.collector.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %v failed with response: %v\nError: %v", r.Request.URL, r, err)
	})
}

// Start 开始爬取
func (s *Scraper) Start() {
	var wg sync.WaitGroup

	for _, category := range s.config.Scraper.Categories {
		for i := 1; i <= s.config.Scraper.CatchPages; i++ {
			wg.Add(1)
			go func(cat string, page int) {
				defer wg.Done()
				url := fmt.Sprintf(s.config.Scraper.BaseURL, cat, page)
				s.limiter.Wait(context.Background())
				if err := s.collector.Visit(url); err != nil {
					log.Printf("Error visiting category %s: %v", cat, err)
				}
			}(category, i)
		}
	}

	wg.Wait()
	s.collector.Wait()
}

func main() {
	// 初始化配置
	cfg := newConfig()

	// 初始化数据库连接
	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建爬虫实例
	scraper := NewScraper(cfg, db)
	scraper.setupCollector()

	// 开始爬取
	scraper.Start()
}
