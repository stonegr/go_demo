package models

import (
	"time"
)

type Article struct {
	ID           int64     `gorm:"primaryKey;comment:文章ID"`
	Title        string    `gorm:"size:255;comment:文章标题"`
	Tags         string    `gorm:"type:text;comment:文章标签,以逗号分隔"`
	Cover        string    `gorm:"size:255;comment:封面图片URL"`
	Excerpt      string    `gorm:"type:text;comment:文章摘要"`
	Listed       bool      `gorm:"comment:是否在列表中显示"`
	DateModified time.Time `gorm:"comment:最后修改时间"`
	Content      string    `gorm:"type:longtext;comment:文章内容"`
	MD5Check     string    `gorm:"size:32;comment:文件内容的MD5哈希值"`
	CreatedAt    time.Time `gorm:"comment:创建时间"`
	UpdatedAt    time.Time `gorm:"comment:更新时间"`
}

func (Article) TableName() string {
	return "articles"
}
