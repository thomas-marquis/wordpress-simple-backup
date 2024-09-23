package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/common"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a backup",
	Long: `Restore a backup.

    Restore the last version:
    wsb restore

    Restore a specific version:
    wsb restore --version 1`,
	Run: func(cmd *cobra.Command, args []string) {
		argsVal, _ := common.ParseCommonArgs(cmd)
		versionID, err := cmd.Flags().GetInt("version")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		hasError := false
		for _, arg := range common.CommonArgs {
			if arg.Required {
				if cmd.Flag(arg.Name).Value.String() == "" {
					fmt.Printf("Error: %s is required\n", arg.Name)
					hasError = true
				}
			}
		}
		if hasError {
			os.Exit(1)
		}

		app := common.GetBackupApp(
			argsVal.DbUsername,
			argsVal.DbPassword,
			argsVal.DbContainer,
			argsVal.WpContentPath,
			argsVal.S3AccessKeyId,
			argsVal.S3SecretAccessKey,
			argsVal.S3Region,
			argsVal.S3Bucket,
		)

		err = app.Restore(argsVal.Name, versionID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().IntP("version", "v", 0, "Version ID to restore")
	for _, arg := range common.CommonArgs {
		saveCmd.Flags().StringP(arg.Name, arg.Short, "", arg.Description)
	}
}
