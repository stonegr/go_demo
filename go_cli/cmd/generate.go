package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	count     int
	outputDir string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate markdown files with random content",
	Long: `Generate multiple markdown files with random content in the specified format.
Each file will contain front matter with meta information and random markdown content.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if count <= 0 {
			return fmt.Errorf("count must be greater than 0")
		}

		if outputDir == "" {
			outputDir = "articles"
		}

		// Create output directory if it doesn't exist
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}

		// Generate files
		for i := 0; i < count; i++ {
			if err := generateFile(i + 1); err != nil {
				return err
			}
		}

		fmt.Printf("Successfully generated %d markdown files in %s\n", count, outputDir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().IntVarP(&count, "count", "c", 5, "Number of files to generate")
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "articles", "Output directory for generated files")
}

func generateFile(id int) error {
	// Generate random title
	title := fmt.Sprintf("Sample Article %d: %s", id, getRandomTitle())

	// Generate random tags
	tags := getRandomTags()

	// Generate random date
	date := time.Now().AddDate(0, -rand.Intn(12), -rand.Intn(30)).Format("2006-01-02")

	// Generate content
	content := generateContent()

	// Create front matter
	frontMatter := fmt.Sprintf(`---
id: %d
title: %s
tags: [%s]
cover: https://www.30secondsofcode.org/assets/cover/compass-400.webp
excerpt: %s
listed: true
dateModified: %s
---

%s`, id, title, strings.Join(tags, ", "), getRandomExcerpt(), date, content)

	// Create filename
	filename := fmt.Sprintf("article-%d.md", id)
	filepath := filepath.Join(outputDir, filename)

	// Write file
	return os.WriteFile(filepath, []byte(frontMatter), 0644)
}

func getRandomTitle() string {
	titles := []string{
		"Understanding Modern Web Development",
		"Best Practices for Code Organization",
		"Introduction to Cloud Computing",
		"Data Structures and Algorithms",
		"Software Testing Fundamentals",
		"DevOps and CI/CD Pipeline",
		"Microservices Architecture",
		"Database Design Principles",
		"Security in Modern Applications",
		"Performance Optimization Techniques",
	}
	return titles[rand.Intn(len(titles))]
}

func getRandomTags() []string {
	allTags := []string{
		"programming", "technology", "web", "development", "coding",
		"software", "engineering", "cloud", "devops", "database",
		"security", "performance", "testing", "architecture", "design",
	}

	// Select 2-4 random tags
	numTags := 2 + rand.Intn(3)
	selectedTags := make([]string, numTags)
	used := make(map[string]bool)

	for i := 0; i < numTags; i++ {
		for {
			tag := allTags[rand.Intn(len(allTags))]
			if !used[tag] {
				selectedTags[i] = tag
				used[tag] = true
				break
			}
		}
	}

	return selectedTags
}

func getRandomExcerpt() string {
	excerpts := []string{
		"Learn the fundamentals of modern web development practices.",
		"Discover best practices for organizing your code effectively.",
		"An introduction to cloud computing concepts and services.",
		"Essential knowledge about data structures and algorithms.",
		"Understanding the importance of software testing.",
		"Building efficient CI/CD pipelines for your projects.",
		"Designing scalable microservices architectures.",
		"Best practices for database design and optimization.",
		"Security considerations in modern applications.",
		"Techniques for optimizing application performance.",
	}
	return excerpts[rand.Intn(len(excerpts))]
}

func generateContent() string {
	contents := []string{
		`# Introduction

Welcome to this comprehensive guide on modern development practices.

## Getting Started

Here's a simple example:

` + "```" + `python
def hello_world():
    print("Hello, World!")
` + "```" + `

## Key Concepts

1. First point
2. Second point
3. Third point

### Important Note

> This is an important note about the topic.

## Conclusion

Thank you for reading!`,

		`# Main Topic

Let's explore some interesting concepts.

## Code Example

` + "```" + `javascript
const greeting = (name) => {
    return 'Hello, ' + name + '!';
};
` + "```" + `

## Benefits

- Improved efficiency
- Better organization
- Enhanced readability

### Tips

1. Keep it simple
2. Follow best practices
3. Test thoroughly

## Summary

That's all for now!`,

		`# Overview

In this article, we'll discuss important concepts.

## Implementation

Here's a practical example:

` + "```" + `go
func main() {
    fmt.Println("Welcome to Go!")
}
` + "```" + `

## Best Practices

1. Write clean code
2. Document properly
3. Review regularly

### Additional Resources

- Link 1
- Link 2
- Link 3

## Final Thoughts

Thank you for your attention!`,
	}
	return contents[rand.Intn(len(contents))]
}
