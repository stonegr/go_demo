package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "blog_cli",
	Short: "A CLI tool for syncing blog posts to database",
	Long:  `A CLI tool that syncs markdown blog posts to a database, tracking changes and updating only modified files.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
