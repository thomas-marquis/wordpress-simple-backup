package common

import "github.com/spf13/cobra"

type Arg struct {
	Name        string
	Short       string
	Description string
	Required    bool
}

var (
	CommonArgs = []Arg{
		{"name", "n", "Name of the backup", true},
		{"db-container", "d", "Name of the mariadb container", true},
		{"db-password", "p", "Password of the mariadb user", true},
		{"db-username", "u", "Username of the mariadb user", true},
		{"wp-content-path", "w", "Path to the wp-content directory", true},
		{"s3-bucket", "b", "Name of the S3 bucket", true},
		{"s3-region", "r", "Region of the S3 bucket", true},
		{"s3-access-key-id", "a", "Access key ID of the S3 bucket", true},
		{"s3-secret-access-key", "s", "Secret access-key of the S3 bucket", true},
	}
)

type CommonArgsValues struct {
	Name              string
	DbContainer       string
	DbPassword        string
	DbUsername        string
	WpContentPath     string
	S3Bucket          string
	S3Region          string
	S3AccessKeyId     string
	S3SecretAccessKey string
}

func ParseCommonArgs(cmd *cobra.Command) (CommonArgsValues, error) {
	var name, dbContainer, dbPassword, dbUsername, wpContentPath, s3Bucket, s3Region, s3AccessKeyId, s3SecretAccessKey string
	name, _ = cmd.Flags().GetString("name")
	dbContainer, _ = cmd.Flags().GetString("db-container")
	dbPassword, _ = cmd.Flags().GetString("db-password")
	dbUsername, _ = cmd.Flags().GetString("db-username")
	wpContentPath, _ = cmd.Flags().GetString("wp-content-path")
	s3Bucket, _ = cmd.Flags().GetString("s3-bucket")
	s3Region, _ = cmd.Flags().GetString("s3-region")
	s3AccessKeyId, _ = cmd.Flags().GetString("s3-access-key-id")
	s3SecretAccessKey, _ = cmd.Flags().GetString("s3-secret-access-key")

	return CommonArgsValues{
		Name:              name,
		DbContainer:       dbContainer,
		DbPassword:        dbPassword,
		DbUsername:        dbUsername,
		WpContentPath:     wpContentPath,
		S3Bucket:          s3Bucket,
		S3Region:          s3Region,
		S3AccessKeyId:     s3AccessKeyId,
		S3SecretAccessKey: s3SecretAccessKey,
	}, nil
}
