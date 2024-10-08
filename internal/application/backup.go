package application

import (
	"errors"

	"github.com/thomas-marquis/wordpress-simple-backup/internal/core"
)

type BackupApplication struct {
	repo           core.Repository
	versionsToKeep int
}

func NewBackupApplication(repo core.Repository, keep int) *BackupApplication {
	return &BackupApplication{
		repo:           repo,
		versionsToKeep: keep,
	}
}

func (a *BackupApplication) Restore(backupName string, versionID int) error {
	// Get backup

	// Get version if exists

	// Restore version
	return errors.New("not implemented")
}

func (a *BackupApplication) List() ([]core.Version, error) {
	return []core.Version{}, errors.New("not implemented")
}

func versionListToString(versions []*core.Version) string {
	var str string

	for _, v := range versions {
		str += v.String() + "\n"
	}

	return str
}
