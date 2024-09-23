package infrastructure

import (
	"log"
	"os"

	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
)

type BackupRepositoryImpl struct {
	wp     *WordpressImpl
	s3     *S3Impl
	db     *MariaDbImpl
	logger *log.Logger
}

var _ core.Repository = &BackupRepositoryImpl{}

func NewBackupRepositoryImpl(wp *WordpressImpl, s3 *S3Impl, db *MariaDbImpl) *BackupRepositoryImpl {
	l := log.New(os.Stdout, "BackupRepositoryImpl", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	return &BackupRepositoryImpl{
		wp:     wp,
		s3:     s3,
		db:     db,
		logger: l,
	}
}

func (b *BackupRepositoryImpl) CreateContentBackupVersion(name string) (*core.Version, error) {
	v := core.NewVersion(b.getLastVersionId(name) + 1)
	return v, nil
}

func (b *BackupRepositoryImpl) GetExistingBackup(name string) (*core.Backup, error) {
	bu := &core.Backup{Name: name, VersionsToKeep: 5, Versions: []core.Version{}}
	return bu, nil
}

func (b *BackupRepositoryImpl) CreateNewBackupVersion(name string) (*core.Version, error) {
	dbFilePath, err := b.db.BackupDatabaseFromDocker(name)
	if err != nil {
		b.logger.Println(err)
		return nil, err
	}

	var contentFilePath string
	contentFilePath, err = b.wp.CreateBackup(name)
	if err != nil {
		b.logger.Println(err)
		if err := b.wp.CleanupArchive(name); err != nil {
			b.logger.Println(err)
			return nil, err
		}
		return nil, err
	}

	v := core.NewVersion(b.getLastVersionId(name) + 1)

	v.WithLocalDbBackupFile(dbFilePath)
	v.WithLocalContentBackupFile(contentFilePath)

	return v, nil
}

func (b *BackupRepositoryImpl) SaveBackup(bu *core.Backup) error {
	return nil
}

func (b *BackupRepositoryImpl) RestoreToVersion(bu *core.Backup, versionID int) error {
	return nil
}

func (b *BackupRepositoryImpl) getLastVersionId(_ string) int {
	return 0
}
