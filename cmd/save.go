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
	Long: `Save the actual site as a new backup version.

    Usage:
    wsb save --name mybackup --db-container mariadb-container --db-password azert123! --db-username toto --wp-content-path ./wp-content --s3-bucket mybucket --s3-region eu-west-1 --s3-access-key-id AKIA1234 --s3-secret-access-key 1234
    `,
	Run: func(cmd *cobra.Command, args []string) {
		argsVal, _ := common.ParseCommonArgs(cmd)

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
			argsVal.Name,
			10,
			argsVal.DbUsername,
			argsVal.DbPassword,
			argsVal.DbContainer,
			argsVal.WpContentPath,
			argsVal.S3AccessKeyId,
			argsVal.S3SecretAccessKey,
			argsVal.S3Region,
			argsVal.S3Bucket,
		)

		err := app.Save()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
}
