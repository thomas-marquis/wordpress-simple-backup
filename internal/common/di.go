package common

import (
	"github.com/thomas-marquis/wordpress-simple-backup/internal/application"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/infrastructure"
)

func GetBackupApp(
	dbUsername string,
	dbPassword string,
	dbContainer string,
	wpContentPath string,
	s3AccessKey string,
	s3SecretKey string,
	s3Region string,
	s3Bucket string,
) *application.BackupApplication {
	repo := infrastructure.NewBackupRepositoryImpl(
		infrastructure.NewWordPressImpl(wpContentPath),
		infrastructure.NewS3Impl(s3AccessKey, s3SecretKey, s3Region, s3Bucket),
		infrastructure.NewMariaDbImpl(dbUsername, dbPassword, dbContainer),
	)
	app := application.NewBackupApplication(repo)
	return app
}
