package models

import (
	"html/template"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID          uint          `gorm:"primarykey;comment:文章ID"`
	Title       string        `gorm:"size:200;not null;comment:文章标题"`
	Summary     string        `gorm:"size:500;comment:文章摘要"`
	Content     string        `gorm:"type:longtext;comment:文章内容"`
	HTMLContent template.HTML `gorm:"-"`
	Category    string        `gorm:"size:20;index;comment:文章分类"`
	PublishTime time.Time     `gorm:"type:date;not null;comment:发布时间"`
	ImageUrl    string        `gorm:"size:255;comment:文章配图URL"`
}

// 获取文章列表
func GetPosts(page int, pageSize int, category string) ([]Post, int64, error) {
	var posts []Post
	var total int64

	query := DB
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 获取总数
	query.Model(&Post{}).Count(&total)

	// 获取分页数据
	err := query.Select("id, title, summary, category, publish_time, image_url, created_at, updated_at, deleted_at").
		Order("publish_time desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts).Error

	return posts, total, err
}

// GetCategories 获取所有分类
func GetCategories() ([]string, error) {
	var categories []string
	err := DB.Model(&Post{}).
		Distinct().
		Pluck("category", &categories).
		Error
	return categories, err
}

func GetPostByID(id int) (*Post, error) {
	var post Post
	result := DB.First(&post, id)
	if result.Error != nil {
		return nil, result.Error
	}
	// 替换HTML内容中的src="/为src="https://www.30secondsofcode.org/
	content := strings.ReplaceAll(post.Content, `/assets/cover/`, `https://www.30secondsofcode.org/assets/cover/`)
	content = strings.ReplaceAll(content, `href="/`, `href="https://www.30secondsofcode.org/`)
	post.HTMLContent = template.HTML(content)
	return &post, nil
}
