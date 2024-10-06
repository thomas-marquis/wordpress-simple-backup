package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/common"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/infrastructure"
)

// zipCmd represents the zip command
var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("zip called")
		cfgPath := cmd.Flag("config").Value.String()
		cfg, err := common.LoadConfig(cfgPath)
		if err != nil {
			fmt.Printf("Error loading config: %s\n", err)
			os.Exit(1)
		}

		if err := infrastructure.CompressFolder(cfg.WpContentPath, cfg.BackupTmpPath+"/wp-content-backup.tar.gz"); err != nil {
			fmt.Printf("Error compressing folder: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	wpCmd.AddCommand(zipCmd)
	common.SetupCommonArgs(zipCmd)
}
