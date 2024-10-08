package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/common"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save the actual site as a new backup version.",
	Long:  `Save the actual site as a new backup version.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		if err := app.Save(); err != nil {
			fmt.Printf("Error saving backup: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
	common.SetupCommonArgs(saveCmd)
}
