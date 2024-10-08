package common

import (
	"github.com/thomas-marquis/wordpress-simple-backup/internal/application"
	"github.com/thomas-marquis/wordpress-simple-backup/internal/infrastructure"
)

func GetBackupApp(cfg Config) (*application.BackupApplication, error) {
	s3, err := infrastructure.NewS3Impl(
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Region,
		cfg.S3BucketName,
		cfg.S3Endpoint,
	)
	if err != nil {
		return nil, err
	}
	db := infrastructure.NewMariaDbImpl(cfg.DbUser, cfg.DbPassword, cfg.DbContainerName)
	repo := infrastructure.NewBackupRepositoryImpl(
		cfg.Host,
		cfg.BackupTmpPath,
		cfg.WpContentPath,
		s3,
		db,
	)
	return application.NewBackupApplication(repo, cfg.VersionToKeep), nil
}
