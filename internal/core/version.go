package core

import (
	"os"
	"time"
)

type Version struct {
	ID        int
	CreatedAt time.Time

	localBackupDbPath      string
	localBackupContentPath string
}

func NewVersion(id int) *Version {
	return &Version{
		ID:        id,
		CreatedAt: time.Now(),
	}
}

func (v *Version) Age() time.Duration {
	return time.Since(v.CreatedAt)
}

func (v *Version) WithLocalDbBackupFile(path string) {
	v.localBackupDbPath = path
}

func (v *Version) WithLocalContentBackupFile(path string) {
	v.localBackupContentPath = path
}

func (v *Version) LocalDbBackupFilePath() string {
	return v.localBackupDbPath
}

func (v *Version) LocalContentBackupFilePath() string {
	return v.localBackupContentPath
}

func (v *Version) AreLocalBackupsAllReady() bool {
	var dbExists, contentExists bool
	_, err := os.Stat(v.localBackupDbPath)
	if err == nil {
		dbExists = true
	}

	_, err = os.Stat(v.localBackupContentPath)
	if err == nil {
		contentExists = true
	}

	return contentExists && dbExists
}
