package common

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Host            string `mapstructure:"SITE_HOST"`
	DbName          string `mapstructure:"SITE_DB_NAME"`
	DbUser          string `mapstructure:"SITE_DB_USER"`
	DbPassword      string `mapstructure:"SITE_DB_PASSWORD"`
	DbContainerName string `mapstructure:"SITE_DB_CONTAINER_NAME"`
	WpContentPath   string `mapstructure:"SITE_WP_CONTENT_PATH"`
	BackupTmpPath   string `mapstructure:"SITE_BACKUP_TMP_PATH"`
	S3AccessKey     string `mapstructure:"SITE_BACKUP_S3_ACCESS_KEY"`
	S3SecretKey     string `mapstructure:"SITE_BACKUP_S3_SECRET_KEY"`
	S3BucketName    string `mapstructure:"SITE_BACKUP_S3_BUCKET_NAME"`
	S3Region        string `mapstructure:"SITE_BACKUP_S3_REGION"`
	S3Endpoint      string `mapstructure:"SITE_BACKUP_S3_ENDPOINT"`
	VersionToKeep   int    `mapstructure:"SITE_BACKUP_VERSION_TO_KEEP"`
}

func LoadConfig(filename string) (Config, error) {
	var c Config
	viper.SetConfigFile(filename)
	viper.SetConfigType("env")

	viper.SetDefault("SITE_BACKUP_TMP_PATH", "/tmp")

	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}

	err := viper.Unmarshal(&c)
	if err != nil {
		return c, err
	}

	c.WpContentPath = strings.TrimSuffix(c.WpContentPath, "/")
	c.BackupTmpPath = strings.TrimSuffix(c.BackupTmpPath, "/")

	return c, nil
}
