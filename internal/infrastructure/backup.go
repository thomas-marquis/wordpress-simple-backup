package infrastructure

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
	infrautils "github.com/thomas-marquis/wordpress-simple-backup/internal/infrastructure/utils"
)

const (
	metadataFile            = "metadata.json"
	dbBackupFileSuffix      = ".sql"
	contentBackupFileSuffix = ".tar.gz"
)

var (
	ErrBackupNotFound  = errors.New("backup not found")
	ErrCorruptedBackup = errors.New("corrupted backup")
)

type backupVersionMeta struct {
	Id                int       `json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	DbBackupFile      string    `json:"db_backup_file"`
	ContentBackupFile string    `json:"content_backup_file"`
}

type backupMeta struct {
	Name string `json:"name"`
}

type BackupRepositoryImpl struct {
	siteName        string
	s3              *S3Impl
	db              *MariaDbImpl
	logger          *log.Logger
	tmp             string
	wpContentFolder string
}

var _ core.Repository = &BackupRepositoryImpl{}

func NewBackupRepositoryImpl(siteName, tmpPath, wpContentFolder string, s3 *S3Impl, db *MariaDbImpl) *BackupRepositoryImpl {
	l := log.New(os.Stdout, "BackupRepositoryImpl", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	return &BackupRepositoryImpl{
		siteName:        siteName,
		s3:              s3,
		db:              db,
		logger:          l,
		tmp:             tmpPath,
		wpContentFolder: wpContentFolder,
	}
}

func (b *BackupRepositoryImpl) CreateContentDump() (core.DumpFile, error) {
	destFile := b.tmp + "/" + b.getContentArchiveFileName()
	err := CompressFolder(b.wpContentFolder, destFile)
	if err != nil {
		return core.NewDumpFile(""), err
	}
	return core.NewDumpFile(destFile), nil
}

func (b *BackupRepositoryImpl) CreateDbDump() (core.DumpFile, error) {
	destFile := b.tmp + "/" + b.getDbDumpFileName()
	err := b.db.BackupDatabaseFromDocker(destFile)
	if err != nil {
		return core.NewDumpFile(""), err
	}
	return core.NewDumpFile(destFile), nil
}

func (b *BackupRepositoryImpl) ListVersions() ([]*core.Version, error) {
	content, err := b.s3.ListFolders(b.siteName)
	if err != nil {
		return nil, err
	}
	var versions []*core.Version
	for _, c := range content {
		v, err := infrautils.ParseVersionDirName(c)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}

func (b *BackupRepositoryImpl) SaveVersion(v *core.Version) error {
	key := infrautils.FormatVersionDirName(v)
	db := v.DbBackupFile()
	dbFileName, err := db.FileName()
	if err != nil {
		return err
	}
	dbFileKey := b.siteName + "/" + key + "/" + dbFileName

	cont := v.ContentBackupFile()
	contFileName, err := cont.FileName()
	if err != nil {
		return err
	}
	contFileKey := b.siteName + "/" + key + "/" + contFileName

	err = b.s3.UploadFile(db.Path, dbFileKey)
	if err != nil {
		return err
	}
	err = b.s3.UploadFile(cont.Path, contFileKey)
	if err != nil {
		if err := b.s3.DeleteFile(dbFileKey); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (b *BackupRepositoryImpl) RemoveVersion(vid core.VersionID) error {
	v, err := b.getVersion(vid)
	if err != nil {
		return err
	}
	dirKey := b.siteName + "/" + infrautils.FormatVersionDirName(v)

	if err := b.s3.DeleteFolder(dirKey); err != nil {
		return err
	}
	return nil
}

func (b *BackupRepositoryImpl) ClearDump(d core.DumpFile) error {
	return d.Remove()
}

func (b *BackupRepositoryImpl) getVersion(vid core.VersionID) (*core.Version, error) {
	versions, err := b.ListVersions()
	if err != nil {
		return nil, err
	}
	for _, ver := range versions {
		if ver.ID == vid {
			return ver, nil
		}
	}
	return nil, ErrBackupNotFound
}

func (b *BackupRepositoryImpl) getContentArchiveFileName() string {
	return fmt.Sprintf("%s.tar.gz", b.siteName)
}

func (b *BackupRepositoryImpl) getDbDumpFileName() string {
	return fmt.Sprintf("%s.sql", b.siteName)
}
