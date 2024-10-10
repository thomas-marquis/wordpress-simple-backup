package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/common"
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

		app, err := common.GetBackupApp(cfg)
		if err != nil {
			fmt.Printf("Error getting backup app: %s\n", err)
			os.Exit(1)
		}

		versions, err := app.List()
		if err != nil {
			fmt.Printf("Error listing versions: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("Existing versions:")
		for _, v := range versions {
			fmt.Println(v.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	common.SetupCommonArgs(listCmd)
}
