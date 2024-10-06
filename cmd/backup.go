package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/common"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/infrastructure"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("backup called")
		cfgPath := cmd.Flag("config").Value.String()
		cfg, err := common.LoadConfig(cfgPath)
		if err != nil {
			fmt.Printf("Error loading config: %s\n", err)
			os.Exit(1)
		}
		mbd := infrastructure.NewMariaDbImpl(cfg.DbUser, cfg.DbPassword, cfg.DbContainerName)
		_, err = mbd.BackupDatabaseFromDocker(cfg.BackupTmpPath + "/db")
		if err != nil {
			fmt.Printf("Error backing up database: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	dbCmd.AddCommand(backupCmd)
	common.SetupCommonArgs(backupCmd)
}
