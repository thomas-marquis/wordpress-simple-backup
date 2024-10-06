package infrastructure

import (
	"encoding/json"
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
	key := infrautils.FormatVersionDirName(b.siteName, v)
	dbFileName, err := v.DbBackupFile().FileName()
	if err != nil {
		return err
	}
	dbFileKey := key + "/" + dbFileName

	return b.s3.UploadFile(v.DbDumpFile.Path, dbFileKey)
}

func (b *BackupRepositoryImpl) RemoveVersion(vid core.VersionID) error {
	return nil
}

func (b *BackupRepositoryImpl) ClearDump(d core.DumpFile) error {
	return nil
}

func (b *BackupRepositoryImpl) getContentArchiveFileName() string {
	return fmt.Sprintf("%s.tar.gz", b.siteName)
}

func (b *BackupRepositoryImpl) getDbDumpFileName() string {
	return fmt.Sprintf("%s.sql", b.siteName)
}

// func (b *BackupRepositoryImpl) CreateNewBackupVersion(name string) (*core.Version, error) {
// 	return nil, nil
// }
//
// func (b *BackupRepositoryImpl) SaveVersion(backupName string, v *core.Version) error {
// 	return nil
// }
//
// func (b *BackupRepositoryImpl) IsBackupExists(name string) (bool, error) {
// 	exists, err := b.s3.IsFolderExists(name)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	return exists, nil
// }
//
// func (b *BackupRepositoryImpl) CreateNewBackup(name string) (*core.Backup, error) {
// 	alreadyExists, err := b.IsBackupExists(name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if alreadyExists {
// 		return nil, errors.New("backup already exists")
// 	}
// 	bu := core.NewBackup(name, 100)
//
// 	m := backupMeta{Name: name}
//
// 	mFileKey := fmt.Sprintf("%s/%s", name, metadataFile)
// 	if err := b.saveJson(m, mFileKey); err != nil {
// 		return nil, err
// 	}
//
// 	return bu, nil
// }
//
// func (b *BackupRepositoryImpl) CreateContentBackupVersion(name string) (*core.Version, error) {
// 	v := core.NewVersion(b.getLastVersionId(name) + 1)
// 	return v, nil
// }
//
// func (b *BackupRepositoryImpl) GetExistingBackup(name string) (*core.Backup, error) {
// 	bu := &core.Backup{Name: name, VersionsToKeep: 5, Versions: []core.Version{}}
// 	return bu, nil
// }
//
// // func (b *BackupRepositoryImpl) CreateNewBackupVersion(name string) (*core.Version, error) {
// // 	dbFilePath, err := b.db.BackupDatabaseFromDocker(name)
// // 	if err != nil {
// // 		b.logger.Println(err)
// // 		return nil, err
// // 	}
// //
// // 	var contentFilePath string
// // 	contentFilePath, err = b.wp.CreateBackup(name)
// // 	if err != nil {
// // 		b.logger.Println(err)
// // 		if err := b.wp.CleanupArchive(name); err != nil {
// // 			b.logger.Println(err)
// // 			return nil, err
// // 		}
// // 		return nil, err
// // 	}
// //
// // 	v := core.NewVersion(b.getLastVersionId(name) + 1)
// //
// // 	v.WithLocalDbBackupFile(dbFilePath)
// // 	v.WithLocalContentBackupFile(contentFilePath)
// //
// // 	return v, nil
// // }
//
// func (b *BackupRepositoryImpl) SaveBackup(bu *core.Backup) error {
// 	// dirsVersions, err := b.s3.ListDir(bu.Name)
// 	// if err != nil && err != ErrBackupNotFound {
// 	// 	return err
// 	// }
// 	//
// 	// if err == ErrBackupNotFound {
// 	// 	// Create backup folder
// 	// }
//
// 	return nil
// }
//
// func (b *BackupRepositoryImpl) RestoreToVersion(bu *core.Backup, versionID int) error {
// 	return nil
// }
//
// func (b *BackupRepositoryImpl) getLastVersionId(_ string) int {
// 	return 0
// }
//
// func (b *BackupRepositoryImpl) getMetadata(backupName string) (backupMeta, error) {
// 	// versionFolders, er := b.s3.ListFolders(backupName)
// 	// if er != nil {
// 	// 	return backupMeta{}, er
// 	// }
// 	//
// 	// m := backupMeta{Name: backupName, Versions: []backupVersionMeta{}}
// 	//
// 	// for _, v := range versionFolders {
// 	// 	vFiles, err := b.s3.ListDir(v)
// 	// 	if err != nil {
// 	// 		return backupMeta{}, err
// 	// 	}
// 	//
// 	// 	vSplited := strings.Split(v, "$")
// 	// 	if len(vSplited) != 2 {
// 	// 		return backupMeta{}, ErrCorruptedBackup
// 	// 	}
// 	// 	vId := vSplited[1]
// 	//
// 	// 	vIdInt, err := strconv.Atoi(vId)
// 	// 	if err != nil {
// 	// 		return backupMeta{}, ErrCorruptedBackup
// 	// 	}
// 	//
// 	// 	var dbBackupFile, contentBackupFile string
// 	// 	for _, f := range vFiles {
// 	// 		if strings.HasSuffix(f, dbBackupFile) {
// 	// 			dbBackupFile = f
// 	// 		}
// 	// 		if strings.HasSuffix(f, contentBackupFile) {
// 	// 			contentBackupFile = f
// 	// 		}
// 	// 	}
// 	// 	if dbBackupFile != "" || contentBackupFile != "" {
// 	// 		return backupMeta{}, ErrCorruptedBackup
// 	// 	}
// 	// 	m.Versions = append(m.Versions, backupVersionMeta{
// 	// 		Id: vIdInt,
// 	// 	})
// 	// }
//
// 	return backupMeta{}, nil
//
// 	// return m, nil
// }
//
// func (b *BackupRepositoryImpl) saveJson(d interface{}, destFile string) error {
// 	buff, err := json.Marshal(d)
// 	if err != nil {
// 		return err
// 	}
//
// 	return b.s3.UploadFileFromMemory(buff, destFile)
// }
