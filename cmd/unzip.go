/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/common"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/infrastructure"
)

// unzipCmd represents the unzip command
var unzipCmd = &cobra.Command{
	Use:   "unzip",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("unzip called")
		cfgPath := cmd.Flag("config").Value.String()
		cfg, err := common.LoadConfig(cfgPath)
		if err != nil {
			fmt.Printf("Error loading config: %s\n", err)
			os.Exit(1)
		}

		if err := infrastructure.UncompressFolder(cfg.BackupTmpPath+"/wp-content-backup.tar.gz", cfg.BackupTmpPath+"/wp-content-restore"); err != nil {
			fmt.Printf("Error compressing folder: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	wpCmd.AddCommand(unzipCmd)
	common.SetupCommonArgs(unzipCmd)
}
