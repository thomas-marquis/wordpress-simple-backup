package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/common"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/infrastructure"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("push called")
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

		err = s3.UploadFile(cfg.BackupTmpPath+"/wp-content-backup.tar.gz", "wp-content-backup.tar.gz")
		if err != nil {
			fmt.Printf("Error downloading file: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	s3Cmd.AddCommand(pushCmd)
	common.SetupCommonArgs(pushCmd)
}
