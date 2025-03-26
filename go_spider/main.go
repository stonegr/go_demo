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

// Article 结构体用于存储文章信息
type Article struct {
	Title       string
	Summary     string
	Category    string
	Content     string
	ImageURL    string
	PublishTime string
}

// 数据库配置
const (
	dbUser     = "go_blog"
	dbPassword = "go_blog"
	dbName     = "go_blog"
)

func main() {
	// 初始化数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:33306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建限速器，每秒允许5个请求
	limiter := rate.NewLimiter(rate.Every(200*time.Millisecond), 1)

	// 创建爬虫实例，允许异步
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		colly.Async(true), // 启用异步模式
	)

	// 设置并发数
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5, // 同时运行的爬虫数量
	})

	// 设置爬虫规则
	c.OnHTML("body > main > article", func(e *colly.HTMLElement) {
		// 获取文章内容
		content, _ := e.DOM.Html()

		// 从上下文中获取之前保存的文章信息
		article, ok := e.Request.Ctx.GetAny("article").(Article)
		if !ok {
			log.Println("Failed to get article from context")
			return
		}

		// 更新文章内容
		article.Content = content

		// 将文章保存到数据库
		insertSQL := `
		INSERT INTO posts (title, summary, category, publish_time, content, image_url)
		VALUES (?, ?, ?, ?, ?, ?)`

		_, err := db.Exec(insertSQL,
			article.Title,
			article.Summary,
			article.Category,
			article.PublishTime,
			article.Content,
			article.ImageURL)

		if err != nil {
			log.Printf("Error inserting article: %v", err)
			return
		}

		log.Printf("Successfully saved article: %s", article.Title)
	})

	// 处理文章列表页面
	c.OnHTML("body > main > section.preview-list", func(e *colly.HTMLElement) {
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

			// 使用限速器
			limiter.Wait(context.Background())
			err := c.Request("GET", link, nil, ctx, nil)
			if err != nil {
				log.Printf("Error visiting article: %v", err)
			}
		})
	})

	// 设置错误处理
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %v failed with response: %v\nError: %v", r.Request.URL, r, err)
	})

	// 使用 WaitGroup 来等待所有爬虫完成
	var wg sync.WaitGroup
	baseURL := "https://www.30secondsofcode.org/%s/p/%d/"
	categories := []string{"js", "css", "html", "react", "node", "git", "python"}

	for _, category := range categories {
		for i := 1; i <= 1; i++ {
			wg.Add(1)
			go func(cat string, page int) {
				defer wg.Done()
				url := fmt.Sprintf(baseURL, cat, page)
				limiter.Wait(context.Background())
				err := c.Visit(url)
				if err != nil {
					log.Printf("Error visiting category %s: %v", cat, err)
				}
			}(category, i)
		}
	}

	// 等待所有爬虫完成
	wg.Wait()
	c.Wait()
}
