package controllers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go_blog/models"
	"go_blog/utils"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

func PostList(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 10

	// 获取分类
	category := c.Param("category")

	// 获取文章列表
	posts, total, err := models.GetPosts(page, pageSize, category)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	// 获取所有分类
	categories, err := models.GetCategories()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	// 计算总页数
	totalPages := (int(total) + pageSize - 1) / pageSize

	c.HTML(http.StatusOK, "index.html", gin.H{
		"posts":      posts,
		"page":       page,
		"totalPages": totalPages,
		"category":   category,
		"categories": categories,
		"totalPosts": total,
	})
}

func PostDetail(c *gin.Context) {
	// 获取文章ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "无效的文章ID",
		})
		return
	}

	// 获取文章详情
	post, err := models.GetPostByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	// // 将内容转换为 template.HTML 类型
	// post.Content = template.HTML(post.Content)

	c.HTML(http.StatusOK, "post.html", gin.H{
		"post": post,
	})
}

func GeneratePostSummary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "无效的文章ID")
		return
	}

	post, err := models.GetPostByID(id)
	if err != nil {
		c.String(http.StatusInternalServerError, "获取文章失败")
		return
	}

	// 提取纯文本内容
	plainText := extractText(post.Content)

	// 构建 OpenAI API 请求
	prompt := fmt.Sprintf("%s\n\n%s", utils.AppConfig.AI.Prompt, plainText)

	// 设置响应头，启用流式响应
	c.Header("Content-Type", "text/plain")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 调用 OpenAI API 并流式传输响应
	err = streamOpenAIResponse(c.Writer, prompt)
	if err != nil {
		c.String(http.StatusInternalServerError, "生成摘要失败")
		return
	}
}

// 提取HTML中的纯文本
func extractText(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	var buf bytes.Buffer
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data + " ")
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(doc)
	return strings.TrimSpace(buf.String())
}

// 流式调用 OpenAI API
func streamOpenAIResponse(w io.Writer, prompt string) error {
	// OpenAI API 配置
	var reader *bufio.Reader
	// 测试环境下使用模拟数据
	if os.Getenv("OPENAI_ENV") == "DEV" {
		// 创建一个模拟的响应内容
		mockResponse := `data: {"choices":[{"delta":{"content":"这是一个"}}]}
data: {"choices":[{"delta":{"content":"测试摘要"}}]}
data: {"choices":[{"delta":{"content":"。这篇文章"}}]}
data: {"choices":[{"delta":{"content":"主要讨论"}}]}
data: {"choices":[{"delta":{"content":"了某个"}}]}
data: {"choices":[{"delta":{"content":"技术主题"}}]}
data: {"choices":[{"delta":{"content":"。"}}]}
data: [DONE]`

		// 创建一个包含模拟数据的Reader
		reader = bufio.NewReader(strings.NewReader(mockResponse))
	} else {
		// 生产环境下的实际API调用
		apiKey := utils.AppConfig.AI.ApiKey
		url := utils.AppConfig.AI.Url

		requestBody := map[string]interface{}{
			"model": utils.AppConfig.AI.Model,
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": prompt,
				},
			},
			"stream": true,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// 读取并转发流式响应
		reader = bufio.NewReader(resp.Body)
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var response struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &response); err != nil {
				continue
			}

			if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
				w.Write([]byte(response.Choices[0].Delta.Content))
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
		}
	}

	return nil
}
