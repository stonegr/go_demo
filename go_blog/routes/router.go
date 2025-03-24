package routes

import (
	"go_blog/controllers"
	"go_blog/utils"
	"html/template"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	if os.Getenv("BLOG_ENV") == "DEV" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// 初始化日志
	utils.InitLogger()

	// 添加日志中间件
	r.Use(utils.GinLogger())

	// 添加自定义模板函数
	r.SetFuncMap(template.FuncMap{
		"subtract": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"max": func(a, b int) int {
			if a > b {
				return a
			}
			return b
		},
		"min": func(a, b int) int {
			if a < b {
				return a
			}
			return b
		},
		"iterate": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
	})

	// 加载模板
	r.LoadHTMLGlob("templates/*")

	// 设置路由
	r.GET("/", controllers.PostList)
	r.GET("/category/:category", controllers.PostList)
	r.GET("/post/:id", controllers.PostDetail)
	r.POST("/post/:id/summary", controllers.GeneratePostSummary)

	return r
}
