package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yourusername/blog-cli/utils"
	"gopkg.in/yaml.v3"
)

var (
	configOutputPath string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate configuration file",
	Long:  `Generate a default configuration file (config.yaml) with default settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 创建默认配置
		cfg := &utils.Config{}
		cfg.Database.Host = "localhost"
		cfg.Database.Port = 3306
		cfg.Database.User = "go_blog"
		cfg.Database.Password = "go_blog"
		cfg.Database.DBName = "go_blog"
		cfg.Scan.Dir = ""
		cfg.Scan.Workers = 5

		// 将配置转换为YAML
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %v", err)
		}

		// 确保输出目录存在
		outputDir := filepath.Dir(configOutputPath)
		if outputDir != "." {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %v", err)
			}
		}

		// 写入配置文件
		if err := os.WriteFile(configOutputPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write config file: %v", err)
		}

		fmt.Printf("Configuration file generated successfully at: %s\n", configOutputPath)
		fmt.Println("\nDefault configuration:")
		fmt.Println(string(data))
		fmt.Println("\nYou can modify this file according to your needs.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// 配置文件输出路径
	configCmd.Flags().StringVarP(&configOutputPath, "output", "o", "config.yaml", "Output path for the configuration file")
}
