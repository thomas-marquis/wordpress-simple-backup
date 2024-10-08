package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/common"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/infrastructure"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all existing backups versions.",
	Long: `List all existing backups versions.

    Usage:
    wsb list
    `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
		cfgPath := cmd.Flag("config").Value.String()
		cfg, err := common.LoadConfig(cfgPath)
		if err != nil {
			fmt.Printf("Error loading config: %s\n", err)
			os.Exit(1)
		}

		s3, err := infrastructure.NewS3Impl(cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Region, cfg.S3BucketName, cfg.S3Endpoint)
		if err != nil {
			fmt.Printf("Error creating S3 client: %s\n", err)
			os.Exit(1)
		}

		content, err := s3.ListFolders("")
		if err != nil {
			fmt.Printf("Error downloading file: %s\n", err)
			os.Exit(1)
		}

		for _, c := range content {
			fmt.Println(c)
		}
	},
}

func init() {
	s3Cmd.AddCommand(listCmd)
	common.SetupCommonArgs(listCmd)
}
